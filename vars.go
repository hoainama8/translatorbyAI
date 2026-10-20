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
var models = []string{"gemini-3.8-flash", "gemini-3.6-flash", "gemini-3.5-flash-lite", "gemini-3.1-flash-lite", "gemma-4-31b-it", "gemma-4-26b-a4b-it"}
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

type Price struct {
	Input  float64
	Output float64
	Cached float64
}

var limitToken float64 = 1000000
var modelPrices = map[string]Price{
	"gemini-3.8-flash":      {Input: 0.75 / limitToken, Output: 3.75 / limitToken, Cached: 0.075 / limitToken},
	"gemini-3.6-flash":      {Input: 0.75 / limitToken, Output: 3.75 / limitToken, Cached: 0.075 / limitToken},
	"gemini-3.5-flash-lite": {Input: 0.3 / limitToken, Output: 2.5 / limitToken, Cached: 0.03 / limitToken},
	"gemini-3.1-flash-lite": {Input: 0.25 / limitToken, Output: 1.5 / limitToken, Cached: 0.025 / limitToken},
	// "gemini-3.6-flash": {Input: ..., Output: ...},
	"gemma-4-31b-it":     {Input: 0.15 / limitToken, Output: 0.4 / limitToken, Cached: 0.06 / limitToken},
	"gemma-4-26b-a4b-it": {Input: 0.09 / limitToken, Output: 0.3 / limitToken, Cached: 0.05 / limitToken},
}
