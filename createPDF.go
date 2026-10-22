package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"github.com/signintech/gopdf"
)

func WritePdf(text string, req translationRequest, outputDir string, myWindow fyne.Window) error {
	if strings.TrimSpace(text) == "" {
		return fmt.Errorf("không có nội dung để xuất ra PDF")
	}
	fontBytes, err := readSystemPDFFont()
	if err != nil {
		return err
	}

	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})
	pdf.AddPage()

	if err := pdf.AddTTFFontData("VietnameseFont", fontBytes); err != nil {
		return fmt.Errorf("không thể nạp font PDF: %w", err)
	}
	if err := pdf.SetFont("VietnameseFont", "", 14); err != nil {
		return fmt.Errorf("không thể thiết lập font PDF: %w", err)
	}

	pdf.SetX(30)
	pdf.SetY(40)

	// if err := pdf.CellWithOption(&gopdf.Rect{W: 500, H: 200}, text, gopdf.CellOption{
	// 	BreakOption: &gopdf.BreakOption{
	// 		Mode:           gopdf.BreakModeIndicatorSensitive,
	// 		BreakIndicator: ' ',
	// 	},
	//});

	files, err := pdf.SplitTextWithWordWrap(text, 535)
	if err != nil {
		dialog.ShowError(err, myWindow)
		return fmt.Errorf("không thể ngắt dòng nội dung PDF: %w", err)
	}
	for _, line := range files {
		if err := pdf.Cell(nil, line); err != nil {
			return fmt.Errorf("không thể ghi nội dung PDF: %w", err)
		}
		pdf.Br(20)
	}

	var outputBuf bytes.Buffer
	if _, err := pdf.WriteTo(&outputBuf); err != nil {
		return fmt.Errorf("không thể tạo dữ liệu PDF: %w", err)
	}

	sourceBase := strings.TrimSuffix(safeStoredName(req.FileName), filepath.Ext(req.FileName))
	if strings.TrimSpace(req.FileName) == "" {
		sourceBase = "translation"
	}
	language := safeStoredName(strings.NewReplacer("/", "-", "\\", "-").Replace(req.TargetLang))
	if strings.TrimSpace(language) == "" {
		return fmt.Errorf("chưa chọn ngôn ngữ đích để đặt tên tệp PDF")
	}
	if err := os.MkdirAll(outputDir, 0o700); err != nil {
		return fmt.Errorf("không thể tạo thư mục lưu PDF: %w", err)
	}
	outputPath := filepath.Join(outputDir, fmt.Sprintf("%s - %s.pdf", sourceBase, language))
	if err := os.WriteFile(outputPath, outputBuf.Bytes(), 0o600); err != nil {
		return fmt.Errorf("không thể lưu tệp PDF: %w", err)
	}
	return nil
}

func readSystemPDFFont() ([]byte, error) {
	candidates := systemPDFFontPaths()
	for _, path := range candidates {
		fontBytes, err := os.ReadFile(path)
		if err == nil {
			return fontBytes, nil
		}
		if !os.IsNotExist(err) {
			return nil, fmt.Errorf("không thể đọc font hệ thống %q: %w", path, err)
		}
	}
	return nil, fmt.Errorf("không tìm thấy font TrueType phù hợp; hãy cài Noto Sans, DejaVu Sans hoặc Arial")
}

func systemPDFFontPaths() []string {
	switch runtime.GOOS {
	case "windows":
		windowsDir := os.Getenv("WINDIR")
		if windowsDir == "" {
			windowsDir = `C:\Windows`
		}
		return []string{
			filepath.Join(windowsDir, "Fonts", "arial.ttf"),
			filepath.Join(windowsDir, "Fonts", "segoeui.ttf"),
			filepath.Join(windowsDir, "Fonts", "tahoma.ttf"),
		}
	case "darwin":
		return []string{
			"/System/Library/Fonts/Supplemental/Arial.ttf",
			"/System/Library/Fonts/Supplemental/Arial Unicode.ttf",
			"/System/Library/Fonts/Supplemental/Verdana.ttf",
		}
	default:
		paths := []string{
			"/usr/share/fonts/truetype/noto/NotoSans-Regular.ttf",
			"/usr/share/fonts/noto/NotoSans-Regular.ttf",
			"/usr/share/fonts/truetype/dejavu/DejaVuSans.ttf",
			"/usr/share/fonts/truetype/liberation2/LiberationSans-Regular.ttf",
		}
		if homeDir, err := os.UserHomeDir(); err == nil {
			paths = append(paths,
				filepath.Join(homeDir, ".local", "share", "fonts", "NotoSans-Regular.ttf"),
				filepath.Join(homeDir, ".fonts", "NotoSans-Regular.ttf"),
			)
		}
		return paths
	}
}
