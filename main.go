package main

import (
	"fmt"
	"image/color"
	"mime"
	"os"
	"path/filepath"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type translationRequest struct {
	APIKey     string
	SourceLang string
	TargetLang string
	FileName   string
	FileType   string
	Model      string
	Document   string // Dữ liệu file đính kèm đã mã hóa Base64
	Extension  string // Định dạng file xuất ra
}

func mimeTypeForFile(fileName string) string {
	switch strings.ToLower(filepath.Ext(fileName)) {
	case ".pdf":
		return "application/pdf"
	case ".docx":
		return "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
	}

	detected := mime.TypeByExtension(strings.ToLower(filepath.Ext(fileName)))
	mediaType, _, err := mime.ParseMediaType(detected)
	if err != nil {
		return "application/octet-stream"
	}
	return mediaType
}

// CHƯƠNG TRÌNH CHÍNH
// ============================================================================

func main() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "không thể tìm thư mục home người dùng: %v\n", err)
		return
	}
	directories, err := applicationDirectoriesFor(homeDir)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	for _, directory := range []string{directories.root, directories.input, directories.output} {
		if err := os.MkdirAll(directory, 0o700); err != nil {
			fmt.Fprintf(os.Stderr, "không thể tạo thư mục %s: %v\n", directory, err)
			return
		}
	}
	if _, err := newTokenUsageTracker(directories.tokenUsage); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	myApp := app.NewWithID(applicationID)

	myWindow := myApp.NewWindow("Phần Mềm Dịch Thuật AI Đa Ngôn Ngữ")

	// CỘT 1: INPUT
	//, refreshInputColumn := createInputColumn(models)

	// CỘT 3: OUTPUT
	outputEditor := widget.NewMultiLineEntry()
	row22, refeshInputColumn := createInputColumn(models)
	row21 := Control(outputEditor, myApp, myWindow, Lang, models, refeshInputColumn, directories.input, directories.output)
	row2 := container.NewGridWithColumns(2, row21, row22)
	row1 := row1(outputEditor)
	// CỘT 2: GIỮA
	mainContent := container.NewGridWithRows(2,
		row1,
		row2,
	)
	lgCustom := layout.CustomPaddedLayout{
		TopPadding:   20,
		LeftPadding:  20,
		RightPadding: 20,
	}
	var bgStart, bgEnd color.NRGBA
	switch GetSystemTheme() {
	case Dark:
		bgStart = color.NRGBA{R: 38, G: 31, B: 29, A: 255}
		bgEnd = color.NRGBA{R: 66, G: 51, B: 46, A: 255}
	default:
		bgStart = color.NRGBA{R: 250, G: 245, B: 237, A: 255}
		bgEnd = color.NRGBA{R: 232, G: 220, B: 205, A: 255}
	}
	bgColor := canvas.NewRadialGradient(bgStart, bgEnd)
	myWindow.SetContent(container.NewStack(bgColor, container.New(lgCustom,
		mainContent)))
	myWindow.Resize(fyne.NewSize(1200, 650))
	// myWindow.SetCloseIntercept(func() {
	// 	myWindow.Hide()
	// })
	myWindow.ShowAndRun()
}
