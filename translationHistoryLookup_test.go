package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLatestTranslationResultMatchesDocumentAndLanguages(t *testing.T) {
	path := filepath.Join(t.TempDir(), "history.json")
	if err := appendTranslationHistory(path, translationHistoryEntry{
		DocumentName:     "report.docx",
		DocumentLanguage: "English",
		TargetLanguage:   "Tiếng Việt",
		Result:           "thành công",
	}); err != nil {
		t.Fatalf("appendTranslationHistory() error = %v", err)
	}
	if err := appendTranslationHistory(path, translationHistoryEntry{
		DocumentName:     "report.docx",
		DocumentLanguage: "English",
		TargetLanguage:   "Tiếng Nhật",
		Result:           "thất bại",
	}); err != nil {
		t.Fatalf("appendTranslationHistory() error = %v", err)
	}

	result, found, err := latestTranslationResult(path, translationHistoryEntry{
		DocumentName:     "report.docx",
		DocumentLanguage: "English",
		TargetLanguage:   "Tiếng Việt",
	})
	if err != nil {
		t.Fatalf("latestTranslationResult() error = %v", err)
	}
	if !found || result != "thành công" {
		t.Errorf("latestTranslationResult() = (%q, %t), want (%q, true)", result, found, "thành công")
	}
}

func TestLatestTranslationResultReturnsLatestMatchingStatus(t *testing.T) {
	path := filepath.Join(t.TempDir(), "history.json")
	for _, result := range []string{"thành công", "thất bại"} {
		if err := appendTranslationHistory(path, translationHistoryEntry{
			DocumentName:     "report.docx",
			DocumentLanguage: "English",
			TargetLanguage:   "Tiếng Việt",
			Result:           result,
		}); err != nil {
			t.Fatalf("appendTranslationHistory() error = %v", err)
		}
	}

	result, found, err := latestTranslationResult(path, translationHistoryEntry{
		DocumentName:     "report.docx",
		DocumentLanguage: "English",
		TargetLanguage:   "Tiếng Việt",
	})
	if err != nil {
		t.Fatalf("latestTranslationResult() error = %v", err)
	}
	if !found || result != "thất bại" {
		t.Errorf("latestTranslationResult() = (%q, %t), want (%q, true)", result, found, "thất bại")
	}
}

func TestLatestTranslationResultRequiresAllFieldsToMatch(t *testing.T) {
	path := filepath.Join(t.TempDir(), "history.json")
	if err := appendTranslationHistory(path, translationHistoryEntry{
		DocumentName:     "report.docx",
		DocumentLanguage: "English",
		TargetLanguage:   "Tiếng Việt",
		Result:           "thành công",
	}); err != nil {
		t.Fatalf("appendTranslationHistory() error = %v", err)
	}

	tests := []struct {
		name    string
		request translationHistoryEntry
	}{
		{
			name: "different document",
			request: translationHistoryEntry{
				DocumentName:     "other.docx",
				DocumentLanguage: "English",
				TargetLanguage:   "Tiếng Việt",
			},
		},
		{
			name: "different source language",
			request: translationHistoryEntry{
				DocumentName:     "report.docx",
				DocumentLanguage: "French",
				TargetLanguage:   "Tiếng Việt",
			},
		},
		{
			name: "different target language",
			request: translationHistoryEntry{
				DocumentName:     "report.docx",
				DocumentLanguage: "English",
				TargetLanguage:   "Tiếng Nhật",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, found, err := latestTranslationResult(path, test.request)
			if err != nil {
				t.Fatalf("latestTranslationResult() error = %v", err)
			}
			if found || result != "" {
				t.Errorf("latestTranslationResult() = (%q, %t), want empty result and no match", result, found)
			}
		})
	}
}

func TestLatestTranslationResultReportsInvalidHistory(t *testing.T) {
	path := filepath.Join(t.TempDir(), "history.json")
	if err := os.WriteFile(path, []byte("{invalid"), 0o600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	_, _, err := latestTranslationResult(path, translationHistoryEntry{})
	if err == nil {
		t.Fatal("latestTranslationResult() error = nil, want malformed history error")
	}
}
