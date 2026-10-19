package main

import "fyne.io/fyne/v2/widget"

// Contasnt
const applicationID = "vn.hoainam.fynetranslator"

type lang string

const (
	vietnam  lang = "Tiếng Việt"
	English  lang = "Tiếng Anh"
	Japanese lang = "Tiếng Nhật"
	Chinese  lang = "Tiếng Trung"
	French   lang = "Tiếng Pháp"
	Korean   lang = "Tiếng Hàn"
	German   lang = "Tiếng Đức"
	Latin    lang = "Tiếng Latin"
	Russian  lang = "Tiếng Nga"
)

// Variable
var Lang = []string{"Tiếng Việt", "Tiếng Anh", "Tiếng Nhật", "Tiếng Trung", "Tiếng Pháp", "Tiếng Hàn", "Tiếng Đức", "Tiếng Latin", "Tiếng Nga"}
var models = []string{"gemini-3.8-flash", "gemini-3.7-flash", "gemini-3.6-flash", "gemma-4-31b-it"}
var sourceLang = append([]string{"Tự động phát hiện"}, Lang...)
var sourceLangSelect = widget.NewSelect(sourceLang, func(s string) {})

//sourceLangSelect.SetSelected("Tự động phát hiện")

var targetLangSelect = widget.NewSelect(Lang, func(s string) {})

//targetLangSelect.SetSelected("Tiếng Việt")

var modelSelect = widget.NewSelect(models, func(s string) {})

//modelSelect.PlaceHolder = "Chọn mô hình AI"

var apiKeyEntry = widget.NewPasswordEntry()

// apiKeyEntry.SetPlaceHolder("Nhập API key")
var progressBar = widget.NewProgressBarInfinite()

// progressBar.Hide()
var supportedOutputExtensions = []string{".txt", ".docx", ".pdf", ".md", ".tex", ".odt"}
var extension = widget.NewSelect(supportedOutputExtensions, func(s string) {})
