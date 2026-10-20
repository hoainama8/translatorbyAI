package main

import (
	"archive/zip"
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"unicode/utf8"

	"github.com/ledongthuc/pdf"
	"github.com/pemistahl/lingua-go"
	"golang.org/x/net/html"
)

const maxLanguageDetectionTextSize = 1 << 20

var (
	languageDetector = lingua.NewLanguageDetectorBuilder().
				FromLanguages(
			lingua.Vietnamese,
			lingua.English,
			lingua.Japanese,
			lingua.Chinese,
			lingua.French,
			lingua.Korean,
			lingua.German,
			lingua.Latin,
			lingua.Russian,
		).
		Build()

	targetLanguageCodes = map[string]lingua.Language{
		"Tiếng Việt":  lingua.Vietnamese,
		"Tiếng Anh":   lingua.English,
		"Tiếng Nhật":  lingua.Japanese,
		"Tiếng Trung": lingua.Chinese,
		"Tiếng Pháp":  lingua.French,
		"Tiếng Hàn":   lingua.Korean,
		"Tiếng Đức":   lingua.German,
		"Tiếng Latin": lingua.Latin,
		"Tiếng Nga":   lingua.Russian,
	}
)

func documentLanguageMatchesTarget(content []byte, fileName, targetLanguage string) (string, bool, error) {
	target, ok := targetLanguageCodes[targetLanguage]
	if !ok {
		return "", false, fmt.Errorf("ngôn ngữ đích không hợp lệ: %q", targetLanguage)
	}

	text, err := extractTextForLanguageDetection(content, fileName)
	if err != nil {
		return "", false, err
	}
	detected, reliable := languageDetector.DetectLanguageOf(text)
	if !reliable {
		return "", false, nil
	}

	return detected.String(), detected == target, nil
}

func extractTextForLanguageDetection(content []byte, fileName string) (string, error) {
	switch strings.ToLower(filepath.Ext(fileName)) {
	case ".pdf":
		return extractPDFText(content)
	case ".docx":
		return extractZipXMLText(content, func(name string) bool {
			return name == "word/document.xml"
		})
	case ".odt":
		return extractZipXMLText(content, func(name string) bool {
			return name == "content.xml"
		})
	case ".epub":
		return extractEPUBText(content)
	case ".html":
		return extractHTMLText(content), nil
	case ".xml":
		return extractXMLText(content)
	case ".txt", ".md", ".json", ".py", ".js", ".sh", ".csv", ".tex", ".rtf":
		if len(content) > maxLanguageDetectionTextSize {
			content = content[:maxLanguageDetectionTextSize]
		}
		return strings.ToValidUTF8(string(content), " "), nil
	default:
		return "", fmt.Errorf("không hỗ trợ trích xuất nội dung %q để xác định ngôn ngữ", filepath.Ext(fileName))
	}
}

func extractPDFText(content []byte) (string, error) {
	reader, err := pdf.NewReader(bytes.NewReader(content), int64(len(content)))
	if err != nil {
		return "", fmt.Errorf("không thể đọc nội dung PDF để xác định ngôn ngữ: %w", err)
	}

	var text strings.Builder
	for pageNumber := 1; pageNumber <= reader.NumPage() && text.Len() < maxLanguageDetectionTextSize; pageNumber++ {
		pageText, err := reader.Page(pageNumber).GetPlainText(nil)
		if err != nil {
			return "", fmt.Errorf("không thể trích xuất nội dung trang %d của PDF: %w", pageNumber, err)
		}
		appendDetectionText(&text, pageText)
	}
	return text.String(), nil
}

func extractZipXMLText(content []byte, include func(string) bool) (string, error) {
	reader, err := zip.NewReader(bytes.NewReader(content), int64(len(content)))
	if err != nil {
		return "", fmt.Errorf("không thể đọc nội dung tài liệu để xác định ngôn ngữ: %w", err)
	}

	var text strings.Builder
	for _, file := range reader.File {
		if text.Len() >= maxLanguageDetectionTextSize {
			break
		}
		if !include(file.Name) {
			continue
		}
		source, err := file.Open()
		if err != nil {
			return "", fmt.Errorf("không thể mở nội dung %s: %w", file.Name, err)
		}
		err = appendXMLText(&text, io.LimitReader(source, int64(maxLanguageDetectionTextSize-text.Len())))
		closeErr := source.Close()
		if err != nil {
			return "", fmt.Errorf("không thể trích xuất nội dung %s: %w", file.Name, err)
		}
		if closeErr != nil {
			return "", fmt.Errorf("không thể đóng nội dung %s: %w", file.Name, closeErr)
		}
	}
	if text.Len() == 0 {
		return "", fmt.Errorf("không tìm thấy nội dung văn bản trong tài liệu")
	}
	return text.String(), nil
}

func extractEPUBText(content []byte) (string, error) {
	reader, err := zip.NewReader(bytes.NewReader(content), int64(len(content)))
	if err != nil {
		return "", fmt.Errorf("không thể đọc nội dung EPUB để xác định ngôn ngữ: %w", err)
	}

	var text strings.Builder
	for _, file := range reader.File {
		if text.Len() >= maxLanguageDetectionTextSize {
			break
		}
		extension := strings.ToLower(filepath.Ext(file.Name))
		if extension != ".xhtml" && extension != ".html" && extension != ".htm" {
			continue
		}
		source, err := file.Open()
		if err != nil {
			return "", fmt.Errorf("không thể mở nội dung %s: %w", file.Name, err)
		}
		htmlBytes, readErr := io.ReadAll(io.LimitReader(source, int64(maxLanguageDetectionTextSize-text.Len())))
		closeErr := source.Close()
		if readErr != nil {
			return "", fmt.Errorf("không thể đọc nội dung %s: %w", file.Name, readErr)
		}
		if closeErr != nil {
			return "", fmt.Errorf("không thể đóng nội dung %s: %w", file.Name, closeErr)
		}
		appendDetectionText(&text, extractHTMLText(htmlBytes))
	}
	if text.Len() == 0 {
		return "", fmt.Errorf("không tìm thấy nội dung văn bản trong tệp EPUB")
	}
	return text.String(), nil
}

func extractXMLText(content []byte) (string, error) {
	var text strings.Builder
	if err := appendXMLText(&text, bytes.NewReader(content)); err != nil {
		return "", fmt.Errorf("không thể trích xuất nội dung XML: %w", err)
	}
	return text.String(), nil
}

func appendXMLText(destination *strings.Builder, source io.Reader) error {
	decoder := xml.NewDecoder(source)
	for destination.Len() < maxLanguageDetectionTextSize {
		token, err := decoder.Token()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		if text, ok := token.(xml.CharData); ok {
			appendDetectionText(destination, string(text))
		}
	}
	return nil
}

func extractHTMLText(content []byte) string {
	root, err := html.Parse(bytes.NewReader(content))
	if err != nil {
		return ""
	}
	var text strings.Builder
	var appendNodeText func(*html.Node)
	appendNodeText = func(node *html.Node) {
		if text.Len() >= maxLanguageDetectionTextSize {
			return
		}
		if node.Type == html.ElementNode && (node.Data == "script" || node.Data == "style" || node.Data == "noscript") {
			return
		}
		if node.Type == html.TextNode {
			appendDetectionText(&text, node.Data)
			return
		}
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			appendNodeText(child)
		}
	}
	appendNodeText(root)
	return text.String()
}

func appendDetectionText(destination *strings.Builder, value string) {
	remaining := maxLanguageDetectionTextSize - destination.Len()
	if remaining <= 0 {
		return
	}
	if len(value) > remaining {
		value = value[:remaining]
		for !utf8.ValidString(value) {
			value = value[:len(value)-1]
		}
	}
	destination.WriteString(value)
	if destination.Len() < maxLanguageDetectionTextSize {
		destination.WriteByte(' ')
	}
}
