package main

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"

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

func Control(outputEditor *widget.Entry, myApp fyne.App, myWindow fyne.Window, Lang []string, models []string, refreshInputColumn func(), inputDirectory string, outDirectory string, translationStatus *translationActivity) fyne.CanvasObject {
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

	attachmentSection := MakeFileAttachmentCard(myWindow, func(fileName string, content []byte) {
		sourceName = fileName
		document = &attachedDocument{
			name: fileName,
			data: append([]byte(nil), content...),
		}

	})

	sourceLangSelect.SetSelected("Tự động phát hiện")
	targetLangSelect.SetSelected("Tiếng Việt")

	// modelSelect := widget.NewSelect(models, func(s string) {})
	modelSelect.PlaceHolder = models[0]
	extension.PlaceHolder = "Chọn định dạng tài liệu dịch"
	progressBar.Hide()

	var translateBtn *widget.Button

	translateBtn = widget.NewButtonWithIcon("DỊCH NGAY VIA AI", theme.NavigateNextIcon(), func() {
		fmt.Println(extension.Selected)
		if document == nil {
			dialog.ShowInformation("Cảnh báo", "Vui lòng nhập văn bản hoặc đính kèm tệp trước!", myWindow)
			return
		}
		if strings.TrimSpace(apiKeyEntry.Text) == "" {
			dialog.ShowInformation("Cảnh báo", "Vui lòng nhập Gemini API key trước!", myWindow)
			return
		}

		if strings.TrimSpace(extension.Selected) == "" {
			dialog.ShowInformation("Cảnh báo", "Vui lòng chọn định dạng file xuất ra", myWindow)
			return
		}
		request := translationRequest{
			APIKey:     apiKeyEntry.Text,
			SourceLang: sourceLangSelect.Selected,
			TargetLang: targetLangSelect.Selected,
			FileName:   sourceName,
			Extension:  extension.Selected,
			Model:      modelSelect.Selected,
		}
		if document != nil {
			request.FileName = document.name
			request.FileType = mimeTypeForFile(document.name)
			request.Document = base64.StdEncoding.EncodeToString(document.data)
			// Backend uses the uploaded bytes as the source of truth and checks
			// its format before forwarding anything to Gemini.

		}
		fmt.Println(request)
		progressBar.Show()
		translateBtn.Disable()
		translationStatus.setInProgress(true)

		go func() {
			response, err := TestTrans(request)
			fyne.Do(func() {
				defer translationStatus.setInProgress(false)
				progressBar.Hide()
				translateBtn.Enable()
				if err != nil {
					fmt.Println(err)
					dialog.ShowError(err, myWindow)
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
				outputEditor.SetText(*response.GetOutputText())
				showSavedTranslation(myWindow, request, *response.GetOutputText())
			})
		}()
	})
	translateBtn.Importance = widget.HighImportance

	optionsForm := widget.NewForm(
		widget.NewFormItem("API key:", apiKeyEntry),
		widget.NewFormItem("Ngôn ngữ nguồn:", sourceLangSelect),
		widget.NewFormItem("Ngôn ngữ đích:", targetLangSelect),
		widget.NewFormItem("Chọn định dạng file dịch", extension),
		widget.NewFormItem("Mô hình AI:", modelSelect),
		//widget.NewFormItem("Lựa chọn:", radio),
		//widget.NewFormItem("Chọn:", combo),
	)

	center1 := container.NewVBox(
		container.NewCenter(attachmentSection),
		optionsForm,
		container.NewCenter(translateBtn),
		progressBar,
		row2_2(inputDirectory, outDirectory, myApp, myWindow),
	)

	//center_green := container.NewStack(canvas.NewHorizontalGradient(color.NRGBA{R: 92, G: 51, B: 23, A: 255} , color.White), centerColumn)
	return center1
}
