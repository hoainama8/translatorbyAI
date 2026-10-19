package main

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func writeBinaryTranslation(path, extension, translatedText string) error {
	switch extension {
	case ".docx":
		return writeDOCX(path, translatedText)
	case ".odt":
		return writeODT(path, translatedText)
	case ".pdf":
		return writePDF(path, translatedText)
	default:
		return fmt.Errorf("định dạng nhị phân không được hỗ trợ: %q", extension)
	}
}

func writeDOCX(path, text string) error {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o600)
	if err != nil {
		return err
	}
	archive := zip.NewWriter(file)
	files := []struct {
		name string
		body string
	}{
		{
			name: "[Content_Types].xml",
			body: `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types">
<Default Extension="rels" ContentType="application/vnd.openxmlformats-package.relationships+xml"/>
<Default Extension="xml" ContentType="application/xml"/>
<Override PartName="/word/document.xml" ContentType="application/vnd.openxmlformats-officedocument.wordprocessingml.document.main+xml"/>
</Types>`,
		},
		{
			name: "_rels/.rels",
			body: `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
<Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/officeDocument" Target="word/document.xml"/>
</Relationships>`,
		},
		{
			name: "word/document.xml",
			body: `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<w:document xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main">
<w:body>` + docxParagraphs(text) + `<w:sectPr/></w:body>
</w:document>`,
		},
	}
	for _, entry := range files {
		writer, err := archive.Create(entry.name)
		if err != nil {
			_ = archive.Close()
			_ = file.Close()
			return err
		}
		if _, err := writer.Write([]byte(entry.body)); err != nil {
			_ = archive.Close()
			_ = file.Close()
			return err
		}
	}
	if err := archive.Close(); err != nil {
		_ = file.Close()
		return err
	}
	return file.Close()
}

func docxParagraphs(text string) string {
	var paragraphs strings.Builder
	for _, line := range strings.Split(strings.ReplaceAll(text, "\r\n", "\n"), "\n") {
		var escaped bytes.Buffer
		_ = xml.EscapeText(&escaped, []byte(line))
		paragraphs.WriteString(`<w:p><w:r><w:t xml:space="preserve">`)
		paragraphs.Write(escaped.Bytes())
		paragraphs.WriteString(`</w:t></w:r></w:p>`)
	}
	return paragraphs.String()
}

func writeODT(path, text string) error {
	return convertWithLibreOffice(path, text, "odt")
}

func writePDF(path, text string) error {
	if err := convertWithLibreOffice(path, text, "pdf"); err != nil {
		return err
	}
	pdf, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("không thể đọc tệp PDF vừa tạo: %w", err)
	}
	if len(pdf) < 5 || string(pdf[:5]) != "%PDF-" {
		return fmt.Errorf("LibreOffice trả về nội dung không phải tệp PDF")
	}
	return nil
}

func convertWithLibreOffice(path, text, outputFormat string) error {
	converter, err := exec.LookPath("libreoffice")
	if err != nil {
		converter, err = exec.LookPath("soffice")
	}
	if err != nil {
		return fmt.Errorf("không tìm thấy LibreOffice (libreoffice/soffice) để tạo tệp %s: %w", outputFormat, err)
	}
	tempDir, err := os.MkdirTemp("", "haitrans-output-*")
	if err != nil {
		return fmt.Errorf("không thể tạo thư mục tạm: %w", err)
	}
	defer os.RemoveAll(tempDir)

	inputPath := filepath.Join(tempDir, "translation.html")
	htmlDocument := `<!doctype html><html><head><meta charset="utf-8"><style>
body { font-family: sans-serif; font-size: 11pt; }
pre { white-space: pre-wrap; overflow-wrap: anywhere; }
</style></head><body><pre>` + html.EscapeString(text) + `</pre></body></html>`
	if err := os.WriteFile(inputPath, []byte(htmlDocument), 0o600); err != nil {
		return fmt.Errorf("không thể tạo nội dung nguồn để chuyển đổi: %w", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	command := exec.CommandContext(ctx, converter, "--headless", "--convert-to", outputFormat, "--outdir", tempDir, inputPath)
	output, err := command.CombinedOutput()
	if err != nil {
		return fmt.Errorf("LibreOffice không thể tạo tệp %s: %w: %s", outputFormat, err, strings.TrimSpace(string(output)))
	}
	outputFile, err := os.ReadFile(filepath.Join(tempDir, "translation."+outputFormat))
	if err != nil {
		return fmt.Errorf("LibreOffice không tạo được tệp %s: %w", outputFormat, err)
	}
	if err := os.WriteFile(path, outputFile, 0o600); err != nil {
		return fmt.Errorf("không thể ghi tệp %s: %w", outputFormat, err)
	}
	return nil
}
