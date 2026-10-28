package main

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type attachedDocument struct {
	name string
	data []byte
}

var document *attachedDocument
var sourceName string

func Control(outputText *string, outputEditor *widget.Entry, myApp fyne.App, myWindow fyne.Window, Lang []string, models []string, refreshInputColumn func(), inputDirectory string, outDirectory string, translationStatus *translationActivity, historyPath string, refreshSummaryTabs func()) fyne.CanvasObject {
	tokenUsagePath, pathErr := defaultTokenUsagePath()
	if pathErr != nil {
		tokenUsagePath = "token_usage.json"
	}
	usageTracker, trackerErr := newTokenUsageTracker(tokenUsagePath)
	if trackerErr != nil {
		usageTracker = &tokenUsageTracker{filePath: tokenUsagePath}
	}
	sessionTokensLabel := widget.NewLabelWithStyle(strconv.Itoa(usageTracker.sessionTokens), fyne.TextAlignLeading, fyne.TextStyle{})
	totalTokensLabel := widget.NewLabelWithStyle(strconv.Itoa(usageTracker.totalTokens), fyne.TextAlignLeading, fyne.TextStyle{})
	if pathErr != nil {
		dialog.ShowError(pathErr, myWindow)
	}
	if trackerErr != nil {
		dialog.ShowError(trackerErr, myWindow)
	}

	sourceLangSelect := widget.NewSelect(sourceLang, func(string) {})
	targetLangSelect := widget.NewSelect(Lang, func(string) {})
	modelSelect := widget.NewSelect(models, func(string) {})
	accountOptions := []string{"Miễn phí", "Trả phí"}
	accountSelect := widget.NewSelect(accountOptions, func(string) {})
	apiKeyEntry := widget.NewPasswordEntry()
	progressBar := widget.NewProgressBarInfinite()
	extension := widget.NewSelect(supportedOutputExtensions, func(string) {})

	attachmentSection := MakeFileAttachmentCard(myWindow, func(fileName string, content []byte) {
		if len(content) == 0 {
			sourceName = ""
			document = nil
			return
		}
		sourceName = fileName
		document = &attachedDocument{
			name: fileName,
			data: append([]byte(nil), content...),
		}

	})

	sourceLangSelect.SetSelected("Tự động phát hiện")
	targetLangSelect.SetSelected("Tiếng Việt")
	accountSelect.SetSelected(accountOptions[0])

	// modelSelect := widget.NewSelect(models, func(s string) {})
	modelSelect.PlaceHolder = models[0]
	modelSelect.SetSelected(models[0])

	progressBar.Hide()
	request := translationRequest{
		APIKey:     apiKeyEntry.Text,
		SourceLang: sourceLangSelect.Selected,
		TargetLang: targetLangSelect.Selected,
		FileName:   sourceName,
		Extension:  extension.Selected,
		Model:      modelSelect.Selected,
		Account:    accountSelect.Selected == accountOptions[0],
	}

	var translateBtn *widget.Button

	translateBtn = widget.NewButtonWithIcon("DỊCH == AI", theme.NavigateNextIcon(), func() {
		// fmt.Println(extension.Selected)
		if document == nil {
			dialog.ShowInformation("Cảnh báo", "Vui lòng nhập văn bản hoặc đính kèm tệp trước!", myWindow)
			return
		}
		if strings.TrimSpace(apiKeyEntry.Text) == "" {
			dialog.ShowInformation("Cảnh báo", "Vui lòng nhập Gemini API key trước!", myWindow)
			return
		}

		// if strings.TrimSpace(extension.Selected) == "" {
		// 	dialog.ShowInformation("Cảnh báo", "Vui lòng chọn định dạng file xuất ra", myWindow)
		// 	return
		// }
		request.APIKey = apiKeyEntry.Text
		request.SourceLang = sourceLangSelect.Selected
		request.TargetLang = targetLangSelect.Selected
		request.FileName = sourceName
		// request.Extension = extension.Selected
		request.Model = modelSelect.Selected
		request.Account = accountSelect.Selected == accountOptions[0]
		detectedLanguage, matchesTarget, err := documentLanguageMatchesTarget(
			document.data,
			document.name,
			targetLangSelect.Selected,
		)
		if err != nil {
			dialog.ShowError(err, myWindow)
			return
		}
		if detectedLanguage == "" {
			detectedLanguage = "Không xác định"
			if sourceLangSelect.Selected != "" && sourceLangSelect.Selected != "Tự động phát hiện" {
				detectedLanguage = sourceLangSelect.Selected
			}
		}
		historyRequest := translationHistoryEntry{
			DocumentName:     document.name,
			DocumentLanguage: detectedLanguage,
			TargetLanguage:   targetLangSelect.Selected,
		}
		previousResult, foundPreviousResult, err := latestTranslationResult(historyPath, historyRequest)
		if err != nil {
			dialog.ShowError(err, myWindow)
			return
		}
		if foundPreviousResult && previousResult == "thành công" {
			dialog.ShowInformation(
				"Cảnh báo",
				fmt.Sprintf("Tài liệu %q đã được dịch thành công từ %s sang %s.", document.name, detectedLanguage, targetLangSelect.Selected),
				myWindow,
			)
			return
		}
		if matchesTarget && !(foundPreviousResult && previousResult == "thất bại") {
			dialog.ShowInformation(
				"Cảnh báo",
				fmt.Sprintf("Nội dung tài liệu đã là %s. Vui lòng chọn ngôn ngữ đích khác.", targetLangSelect.Selected),
				myWindow,
			)
			return
		}

		if document != nil {
			request.FileName = document.name
			request.FileType = mimeTypeForFile(document.name)
			request.Document = base64.StdEncoding.EncodeToString(document.data)
			// Backend uses the uploaded bytes as the source of truth and checks
			// its format before forwarding anything to Gemini.

		}
		progressBar.Show()
		translateBtn.Disable()
		translationStatus.setInProgress(true)

		translationRequestForCall := request
		go func() {
			var response *InteractionsJSON
			var err error
			if translationRequestForCall.Account {
				response, err = TestTrans(translationRequestForCall)
			} else {
				response, err = BacthTrans(translationRequestForCall)
			}
			fyne.Do(func() {
				defer translationStatus.setInProgress(false)
				progressBar.Hide()
				translateBtn.Enable()
				//modelOutputText := (*string)(nil)
				// if response != nil {
				// 	modelOutputText = response.GetOutputText()
				// }
				result := "thành công"
				if err != nil || response == nil || response.GetOutputText() == nil {
					result = "thất bại"
				}
				if historyErr := appendTranslationHistory(historyPath, translationHistoryEntry{
					Date:             time.Now().UTC(),
					DocumentName:     translationRequestForCall.FileName,
					DocumentLanguage: detectedLanguage,
					TargetLanguage:   translationRequestForCall.TargetLang,
					Result:           result,
				}); historyErr != nil {
					dialog.ShowError(historyErr, myWindow)
				}
				if refreshSummaryTabs != nil {
					defer refreshSummaryTabs()
				}
				if err != nil {
					fmt.Println(err)
					dialog.ShowError(err, myWindow)
					return
				}
				if response == nil || response.GetOutputText() == nil {
					dialog.ShowError(fmt.Errorf("bản dịch không có nội dung văn bản"), myWindow)
					return
				}
				if response.GetUsage() != nil {
					if err := response.SaveUsage(usageTracker); err != nil {
						dialog.ShowError(err, myWindow)
					} else {
						sessionTokensLabel.SetText(strconv.Itoa(usageTracker.sessionTokens))
						totalTokensLabel.SetText(strconv.Itoa(usageTracker.totalTokens))
						refreshInputColumn()
					}
				}
				*outputText = *response.GetOutputText()
				outputEditor.SetText(*outputText)
			})
		}()
	})
	translateBtn.Importance = widget.HighImportance

	optionsForm := widget.NewForm(
		widget.NewFormItem("API key:", apiKeyEntry),
		widget.NewFormItem("Ngôn ngữ nguồn:", sourceLangSelect),
		widget.NewFormItem("Ngôn ngữ đích:", targetLangSelect),
		widget.NewFormItem("Tài khoản:", accountSelect),
		//widget.NewFormItem("Chọn định dạng file dịch", extension),
		widget.NewFormItem("Mô hình AI:", modelSelect),
		//widget.NewFormItem("Lựa chọn:", radio),
		//widget.NewFormItem("Chọn:", combo),
	)
	center1 := container.NewVBox(
		container.NewCenter(attachmentSection),
		optionsForm,
		container.NewCenter(translateBtn),
		progressBar,
		row2_2(inputDirectory, outDirectory, historyPath, myApp, myWindow, &request, outputText),
	)

	//center_green := container.NewStack(canvas.NewHorizontalGradient(color.NRGBA{R: 92, G: 51, B: 23, A: 255} , color.White), centerColumn)
	return center1
}
