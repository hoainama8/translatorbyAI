package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func TestTranslationHistoryTableClipsTextAndAlternatesRows(t *testing.T) {
	app := test.NewApp()
	t.Cleanup(app.Quit)

	table := translationHistoryTable([]translationHistoryEntry{
		{DocumentName: "document-with-a-very-long-name.pdf"},
		{DocumentName: "second-document.pdf"},
	})
	newCell := func(row int) (*canvas.Rectangle, *widget.Label) {
		t.Helper()
		cell := table.CreateCell().(*fyne.Container)
		table.UpdateCell(widget.TableCellID{Row: row, Col: 1}, cell)
		return cell.Objects[0].(*canvas.Rectangle), cell.Objects[1].(*widget.Label)
	}

	headerBackground, header := newCell(0)
	if header.Truncation != fyne.TextTruncateEllipsis {
		t.Errorf("header truncation = %v, want ellipsis", header.Truncation)
	}
	if got, want := headerBackground.FillColor, theme.Color(theme.ColorNameBackground); got != want {
		t.Errorf("header background = %v, want default background %v", got, want)
	}

	firstBackground, first := newCell(1)
	if first.Truncation != fyne.TextTruncateEllipsis {
		t.Errorf("data truncation = %v, want ellipsis", first.Truncation)
	}
	if first.Text != "document-with-a-very-long-name.pdf" {
		t.Errorf("document name = %q, want full source text for renderer truncation", first.Text)
	}
	if got, want := firstBackground.FillColor, theme.Color(theme.ColorNameInputBackground); got != want {
		t.Errorf("first row background = %v, want alternate background %v", got, want)
	}

	secondBackground, _ := newCell(2)
	if got, want := secondBackground.FillColor, theme.Color(theme.ColorNameBackground); got != want {
		t.Errorf("second row background = %v, want default background %v", got, want)
	}
}

func TestTranslationHistoryTableViewFillsContainer(t *testing.T) {
	app := test.NewApp()
	t.Cleanup(app.Quit)

	table := translationHistoryTable(nil)
	view := translationHistoryTableView(table)
	size := fyne.NewSize(1500, 700)
	view.Resize(size)

	if got := table.Size(); got != size {
		t.Errorf("table size = %v, want container size %v", got, size)
	}
}

func TestEnsureTranslationHistoryFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "haitranslator", "history.json")
	if err := ensureTranslationHistoryFile(path); err != nil {
		t.Fatalf("ensureTranslationHistoryFile() error = %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}
	if string(data) != "[]\n" {
		t.Errorf("new history content = %q, want []", data)
	}
}

func TestAppendTranslationHistory(t *testing.T) {
	path := filepath.Join(t.TempDir(), "history.json")
	firstDate := time.Date(2026, time.October, 20, 1, 2, 3, 0, time.UTC)
	entries := []translationHistoryEntry{
		{
			Date:             firstDate,
			DocumentName:     "source.docx",
			DocumentLanguage: "Vietnamese",
			TargetLanguage:   "Tiếng Anh",
			Result:           "thành công",
		},
		{
			DocumentName:     "failed.pdf",
			DocumentLanguage: "English",
			TargetLanguage:   "Tiếng Việt",
			Result:           "thất bại",
		},
	}

	for _, entry := range entries {
		if err := appendTranslationHistory(path, entry); err != nil {
			t.Fatalf("appendTranslationHistory() error = %v", err)
		}
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}
	var saved []translationHistoryEntry
	if err := json.Unmarshal(data, &saved); err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}
	if len(saved) != len(entries) {
		t.Fatalf("history entries = %d, want %d", len(saved), len(entries))
	}
	if saved[0] != entries[0] {
		t.Errorf("first history entry = %+v, want %+v", saved[0], entries[0])
	}
	if saved[1].DocumentName != entries[1].DocumentName ||
		saved[1].DocumentLanguage != entries[1].DocumentLanguage ||
		saved[1].TargetLanguage != entries[1].TargetLanguage ||
		saved[1].Result != entries[1].Result {
		t.Errorf("second history entry = %+v, want values from %+v", saved[1], entries[1])
	}
	if saved[1].Date.IsZero() {
		t.Error("second history entry date is zero")
	}
}

func TestAppendTranslationHistoryRejectsInvalidResult(t *testing.T) {
	err := appendTranslationHistory(filepath.Join(t.TempDir(), "history.json"), translationHistoryEntry{
		Result: "unknown",
	})
	if err == nil {
		t.Fatal("appendTranslationHistory() error = nil, want invalid result error")
	}
}

func TestAppendTranslationHistoryDoesNotOverwriteInvalidJSON(t *testing.T) {
	path := filepath.Join(t.TempDir(), "history.json")
	if err := os.WriteFile(path, []byte("{invalid"), 0o600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	err := appendTranslationHistory(path, translationHistoryEntry{Result: "thành công"})
	if err == nil {
		t.Fatal("appendTranslationHistory() error = nil, want malformed JSON error")
	}
	data, readErr := os.ReadFile(path)
	if readErr != nil {
		t.Fatalf("ReadFile() error = %v", readErr)
	}
	if string(data) != "{invalid" {
		t.Errorf("malformed history changed to %q", data)
	}
}
