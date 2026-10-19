package main

// func TransBtn(refes) *widget.Button {
// 	translateBtn := widget.NewButtonWithIcon("DỊCH NGAY VIA AI", theme.NavigateNextIcon(), func() {
// 		if document == nil {
// 			dialog.ShowInformation("Cảnh báo", "Vui lòng nhập văn bản hoặc đính kèm tệp trước!", myWindow)
// 			return
// 		}
// 		if strings.TrimSpace(apiKeyEntry.Text) == "" {
// 			dialog.ShowInformation("Cảnh báo", "Vui lòng nhập Gemini API key trước!", myWindow)
// 			return
// 		}

// 		request := translationRequest{
// 			APIKey:     apiKeyEntry.Text,
// 			SourceLang: sourceLangSelect.Selected,
// 			TargetLang: targetLangSelect.Selected,
// 			FileName: sourceName,
// 			Model:    modelSelect.Selected,
// 		}
// 		if document != nil {
// 			request.FileName = document.name
// 			request.FileType = mimeTypeForFile(document.name)
// 			request.Document = append([]byte(nil), document.data...)
// 			// Backend uses the uploaded bytes as the source of truth and checks
// 			// its format before forwarding anything to Gemini.

// 		}
// 		fmt.Println(request.Model, request.SourceLang, request.TargetLang, request.FileName, len(request.Document))
// 		progressBar.Show()
// 		translateBtn.Disable()

// 		go func() {
// 			response, err := TestTrans(request)
// 			fyne.Do(func() {
// 				progressBar.Hide()
// 				translateBtn.Enable()
// 				if err != nil {
// 					fmt.Println(err)
// 					dialog.ShowError(err, myWindow)
// 					return
// 				}
// 				if response.GetUsage() != nil {
// 					if err := response.SaveUsage(usageTracker); err != nil {
// 						dialog.ShowError(err, myWindow)
// 					} else {
// 						sessionTokensLabel.SetText(strconv.Itoa(usageTracker.sessionTokens))
// 						totalTokensLabel.SetText(strconv.Itoa(usageTracker.totalTokens))
// 						refreshInputColumn()
// 					}
// 				}
// 				outputEditor.SetText(*response.GetOutputText())
// 				showSavedTranslation(myWindow, request, *response.GetOutputText())
// 			})
// 		}()
// 	})
// 	return translateBtn
// }
