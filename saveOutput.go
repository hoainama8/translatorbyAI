package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
)

func showSavedTranslation(win fyne.Window, req translationRequest, translatedText string) {
	directories, err := defaultApplicationDirectories()
	if err != nil {
		dialog.ShowError(err, win)
		return
	}
	textFile, err := saveTranslationFile(directories.output, req.FileName, req.TargetLang, translatedText, req.Extension)
	if err != nil {
		dialog.ShowError(fmt.Errorf("Không thể lưu tệp dịch: %v", err), win)
		return
	}

	message := fmt.Sprintf("Đã lưu tệp bản dịch: %s", textFile)
	dialog.ShowInformation("Thành công", message, win)

}

func saveTranslationFile(outputDir, sourceName, targetLang, translatedText string, extension string) (string, error) {
	if !isSupportedOutputExtension(extension) {
		return "", fmt.Errorf("định dạng đầu ra không được hỗ trợ: %q", extension)
	}
	sourceBase := strings.TrimSuffix(safeStoredName(sourceName), filepath.Ext(sourceName))
	if sourceBase == "document" && strings.TrimSpace(sourceName) == "" {
		sourceBase = "translation"
	}
	language := safeStoredName(strings.NewReplacer("/", "-", "\\", "-").Replace(targetLang))
	filename := fmt.Sprintf("%s - %s%s", sourceBase, language, extension)
	if err := os.MkdirAll(outputDir, 0o700); err != nil {
		return "", fmt.Errorf("không thể tạo thư mục đầu ra: %w", err)
	}
	path := filepath.Join(outputDir, filename)
	if isBinaryOutputExtension(extension) {
		if err := writeBinaryTranslation(path, extension, translatedText); err != nil {
			return "", fmt.Errorf("không thể tạo tệp %s: %w", extension, err)
		}
	} else if err := os.WriteFile(path, []byte(translatedText), 0o600); err != nil {
		return "", fmt.Errorf("không thể ghi nội dung dịch: %w", err)
	}
	return path, nil
}

func isSupportedOutputExtension(extension string) bool {
	for _, supportedExtension := range supportedOutputExtensions {
		if extension == supportedExtension {
			return true
		}
	}
	return false
}

func isBinaryOutputExtension(extension string) bool {
	switch extension {
	case ".docx", ".pdf", ".odt":
		return true
	default:
		return false
	}
}
