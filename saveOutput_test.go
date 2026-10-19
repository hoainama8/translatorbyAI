package main

import (
	"archive/zip"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestSaveTranslationTextUsesStandardExtension(t *testing.T) {
	outputDir := t.TempDir()
	path, err := saveTranslationFile(outputDir, "source.pdf", "Tiếng Việt", "Xin chào", ".txt")
	if err != nil {
		t.Fatalf("saveTranslationText returned an error: %v", err)
	}
	if got, want := filepath.Base(path), "source - Tiếng Việt.txt"; got != want {
		t.Fatalf("unexpected output filename: got %q, want %q", got, want)
	}
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("reading output: %v", err)
	}
	if string(content) != "Xin chào" {
		t.Fatalf("unexpected output content: %q", content)
	}
}

func TestSaveTranslationTextRejectsUnsupportedExtension(t *testing.T) {
	if _, err := saveTranslationFile(t.TempDir(), "source.txt", "Tiếng Việt", "text", "pdf"); err == nil {
		t.Fatal("expected unsupported extension to be rejected")
	}
}

func TestWriteDOCXCreatesValidDocumentArchive(t *testing.T) {
	path := filepath.Join(t.TempDir(), "translation.docx")
	if err := writeDOCX(path, "Xin chào\nDịch thuật"); err != nil {
		t.Fatalf("writeDOCX returned an error: %v", err)
	}
	archive, err := zip.OpenReader(path)
	if err != nil {
		t.Fatalf("opening DOCX archive: %v", err)
	}
	defer archive.Close()
	document := readZipEntry(t, archive.File, "word/document.xml")
	for _, want := range []string{"Xin chào", "Dịch thuật", "<w:p>"} {
		if !strings.Contains(document, want) {
			t.Errorf("DOCX document is missing %q", want)
		}
	}
}

func TestWriteODTCreatesValidDocumentArchive(t *testing.T) {
	path := filepath.Join(t.TempDir(), "translation.odt")
	if err := writeODT(path, "Xin chào\nDịch thuật"); err != nil {
		t.Fatalf("writeODT returned an error: %v", err)
	}
	archive, err := zip.OpenReader(path)
	if err != nil {
		t.Fatalf("opening ODT archive: %v", err)
	}
	defer archive.Close()
	if len(archive.File) == 0 || archive.File[0].Name != "mimetype" || archive.File[0].Method != zip.Store {
		t.Fatal("ODT mimetype must be the first, uncompressed archive entry")
	}
	mimetype := readZipEntry(t, archive.File, "mimetype")
	if mimetype != "application/vnd.oasis.opendocument.text" {
		t.Fatalf("unexpected ODT mimetype: %q", mimetype)
	}
	content := readZipEntry(t, archive.File, "content.xml")
	for _, want := range []string{"Xin chào", "Dịch thuật", "<text:p"} {
		if !strings.Contains(content, want) {
			t.Errorf("ODT content is missing %q", want)
		}
	}
}

func TestWritePDFCreatesPDFWhenLibreOfficeIsAvailable(t *testing.T) {
	libreOfficePath(t)
	path := filepath.Join(t.TempDir(), "translation.pdf")
	if err := writePDF(path, "Xin chào, thế giới"); err != nil {
		t.Fatalf("writePDF returned an error: %v", err)
	}
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("reading PDF: %v", err)
	}
	if len(content) < 5 || string(content[:5]) != "%PDF-" {
		t.Fatal("output does not have a PDF signature")
	}
}

func TestDOCXAndODTOpenInLibreOffice(t *testing.T) {
	converter := libreOfficePath(t)
	tempDir := t.TempDir()
	for _, extension := range []string{".docx", ".odt"} {
		path := filepath.Join(tempDir, "translation"+extension)
		if err := writeBinaryTranslation(path, extension, "Xin chào\nDịch thuật"); err != nil {
			t.Fatalf("creating %s: %v", extension, err)
		}
		outputDir := filepath.Join(tempDir, strings.TrimPrefix(extension, "."))
		if err := os.Mkdir(outputDir, 0o700); err != nil {
			t.Fatalf("creating conversion directory: %v", err)
		}
		output, err := exec.Command(converter, "--headless", "--convert-to", "pdf", "--outdir", outputDir, path).CombinedOutput()
		if err != nil {
			t.Fatalf("LibreOffice could not open %s: %v: %s", extension, err, strings.TrimSpace(string(output)))
		}
		if _, err := os.Stat(filepath.Join(outputDir, "translation.pdf")); err != nil {
			t.Fatalf("LibreOffice did not create a PDF from %s: %v", extension, err)
		}
	}
}

func libreOfficePath(t *testing.T) string {
	t.Helper()
	converter, err := exec.LookPath("libreoffice")
	if err != nil {
		converter, err = exec.LookPath("soffice")
	}
	if err != nil {
		t.Skip("LibreOffice is not installed")
	}
	return converter
}

func readZipEntry(t *testing.T, files []*zip.File, name string) string {
	t.Helper()
	for _, file := range files {
		if file.Name != name {
			continue
		}
		reader, err := file.Open()
		if err != nil {
			t.Fatalf("opening %s: %v", name, err)
		}
		defer reader.Close()
		content, err := io.ReadAll(reader)
		if err != nil {
			t.Fatalf("reading %s: %v", name, err)
		}
		return string(content)
	}
	t.Fatalf("archive entry %s not found", name)
	return ""
}
