package main

import (
	"archive/zip"
	"bytes"
	"testing"
)

func TestDocumentLanguageMatchesTarget(t *testing.T) {
	tests := []struct {
		name           string
		content        string
		targetLanguage string
		wantDetected   string
		wantMatch      bool
	}{
		{
			name: "same language",
			content: "Đây là một đoạn văn dài bằng tiếng Việt để kiểm tra khả năng nhận diện ngôn ngữ " +
				"của tài liệu. Nội dung này mô tả quy trình làm việc và những thông tin cần thiết.",
			targetLanguage: "Tiếng Việt",
			wantDetected:   "Vietnamese",
			wantMatch:      true,
		},
		{
			name: "different language",
			content: "This is a longer English paragraph used to verify the local language detector. " +
				"The document describes a process and includes several complete sentences.",
			targetLanguage: "Tiếng Việt",
			wantDetected:   "English",
			wantMatch:      false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			detected, matches, err := documentLanguageMatchesTarget([]byte(test.content), "sample.txt", test.targetLanguage)
			if err != nil {
				t.Fatalf("documentLanguageMatchesTarget() error = %v", err)
			}
			if detected != test.wantDetected {
				t.Errorf("detected language = %q, want %q", detected, test.wantDetected)
			}
			if matches != test.wantMatch {
				t.Errorf("matches target = %t, want %t", matches, test.wantMatch)
			}
		})
	}
}

func TestDocumentLanguageMatchesTargetExtractsDocxText(t *testing.T) {
	content := createTestDocx(t, `<w:document xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main"><w:body><w:p><w:r><w:t>Đây là tài liệu tiếng Việt với nội dung văn bản đủ dài để hệ thống có thể nhận diện đúng ngôn ngữ và so sánh với lựa chọn của người dùng.</w:t></w:r></w:p></w:body></w:document>`)

	detected, matches, err := documentLanguageMatchesTarget(content, "sample.docx", "Tiếng Việt")
	if err != nil {
		t.Fatalf("documentLanguageMatchesTarget() error = %v", err)
	}
	if detected != "Vietnamese" {
		t.Errorf("detected language = %q, want Vietnamese", detected)
	}
	if !matches {
		t.Error("matches target = false, want true")
	}
}

func TestDocumentLanguageMatchesTargetRejectsInvalidTarget(t *testing.T) {
	_, _, err := documentLanguageMatchesTarget([]byte("sample"), "sample.txt", "")
	if err == nil {
		t.Fatal("documentLanguageMatchesTarget() error = nil, want invalid target error")
	}
}

func createTestDocx(t *testing.T, documentXML string) []byte {
	t.Helper()

	var output bytes.Buffer
	archive := zip.NewWriter(&output)
	file, err := archive.Create("word/document.xml")
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if _, err := file.Write([]byte(documentXML)); err != nil {
		t.Fatalf("Write() error = %v", err)
	}
	if err := archive.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}
	return output.Bytes()
}
