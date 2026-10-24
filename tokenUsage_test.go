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

func costTableCell(table *widget.Table, row, col int) (*canvas.Rectangle, *widget.Label) {
	cell := table.CreateCell().(*fyne.Container)
	table.UpdateCell(widget.TableCellID{Row: row, Col: col}, cell)
	return cell.Objects[0].(*canvas.Rectangle), cell.Objects[1].(*widget.Label)
}

func TestCostTableShowsModelUsageAndCost(t *testing.T) {
	app := test.NewApp()
	t.Cleanup(app.Quit)
	entries := tokenUsageData{
		"gemini-3.8-flash": {{
			StartedAt:     time.Date(2026, time.October, 21, 3, 30, 0, 0, time.UTC),
			InputTokens:   1000,
			CachedTokens:  200,
			OutputTokens:  100,
			ThoughtTokens: 50,
		}},
	}
	table := CostTable(entries, "gemini-3.8-flash")
	if rows, columns := table.Length(); rows != 3 || columns != 6 {
		t.Fatalf("table dimensions = %dx%d, want 3x6", rows, columns)
	}

	want := []string{"2026-10-21", "03:30:00", "1000", "150", "200", "$0.001178"}
	for col, expected := range want {
		_, label := costTableCell(table, 1, col)
		if got := label.Text; got != expected {
			t.Errorf("cell (1,%d) = %q, want %q", col, got, expected)
		}
	}
	headerBackground, _ := costTableCell(table, 0, 0)
	if got, want := headerBackground.FillColor, theme.Color(theme.ColorNameBackground); got != want {
		t.Errorf("header background = %v, want default background %v", got, want)
	}
	firstRowBackground, _ := costTableCell(table, 1, 0)
	if got, want := firstRowBackground.FillColor, theme.Color(theme.ColorNameInputBackground); got != want {
		t.Errorf("first row background = %v, want alternate background %v", got, want)
	}
	totalBackground, totalLabel := costTableCell(table, 2, 5)
	if got, want := totalBackground.FillColor, theme.Color(theme.ColorNameBackground); got != want {
		t.Errorf("total row background = %v, want default background %v", got, want)
	}
	if got := totalLabel.Text; got != "$0.001178" {
		t.Errorf("total cost = %q, want %q", got, "$0.001178")
	}
}

func TestCostTableShowsEmptyStateForModel(t *testing.T) {
	app := test.NewApp()
	t.Cleanup(app.Quit)
	table := CostTable(tokenUsageData{"another-model": {{InputTokens: 1}}}, "gemini-3.8-flash")
	if rows, columns := table.Length(); rows != 3 || columns != 6 {
		t.Fatalf("table dimensions = %dx%d, want 3x6", rows, columns)
	}
	_, label := costTableCell(table, 1, 0)
	if got := label.Text; got != "Chưa có dữ liệu" {
		t.Errorf("empty-state text = %q, want %q", got, "Chưa có dữ liệu")
	}
	_, totalLabel := costTableCell(table, 2, 5)
	if got := totalLabel.Text; got != "$0.000000" {
		t.Errorf("empty total cost = %q, want %q", got, "$0.000000")
	}
}

func TestTokenUsageTrackerPersistsPerModelAndSession(t *testing.T) {
	filePath := filepath.Join(t.TempDir(), "use_token.json")
	tracker, err := newTokenUsageTracker(filePath)
	if err != nil {
		t.Fatalf("newTokenUsageTracker() error = %v", err)
	}
	if err := tracker.Add("gemini-flash", 120); err != nil {
		t.Fatalf("Add(gemini-flash, 120) error = %v", err)
	}
	if tracker.sessionTokens != 120 || tracker.totalTokens != 120 {
		t.Fatalf("after first Add: session=%d total=%d, want 120/120", tracker.sessionTokens, tracker.totalTokens)
	}
	if err := tracker.Add("gemma", 30); err != nil {
		t.Fatalf("Add(gemma, 30) error = %v", err)
	}

	reloaded, err := newTokenUsageTracker(filePath)
	if err != nil {
		t.Fatalf("reloading tracker error = %v", err)
	}
	if reloaded.sessionTokens != 0 || reloaded.totalTokens != 150 {
		t.Fatalf("reloaded tracker: session=%d total=%d, want 0/150", reloaded.sessionTokens, reloaded.totalTokens)
	}
	if err := reloaded.Add("gemini-flash", 20); err != nil {
		t.Fatalf("Add(gemini-flash, 20) error = %v", err)
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}
	var saved tokenUsageData
	if err := json.Unmarshal(data, &saved); err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}
	totalTokens := 0
	for _, sessions := range saved {
		for _, session := range sessions {
			totalTokens += session.TotalTokens
		}
	}
	if totalTokens != 170 {
		t.Errorf("saved total tokens = %d, want 170", totalTokens)
	}
	if got := findTokenUsageSession(saved["gemini-flash"], tracker.sessionID).TotalTokens; got != 120 {
		t.Errorf("first session gemini-flash tokens = %d, want 120", got)
	}
	if got := findTokenUsageSession(saved["gemma"], tracker.sessionID).TotalTokens; got != 30 {
		t.Errorf("first session gemma tokens = %d, want 30", got)
	}
	if got := findTokenUsageSession(saved["gemini-flash"], reloaded.sessionID).TotalTokens; got != 20 {
		t.Errorf("second session gemini-flash tokens = %d, want 20", got)
	}
}

func findTokenUsageSession(sessions []tokenUsageSession, sessionID string) tokenUsageSession {
	for _, session := range sessions {
		if session.SessionID == sessionID {
			return session
		}
	}
	return tokenUsageSession{}
}

func TestDecodeTokenUsageDataMigratesLegacySchema(t *testing.T) {
	legacy := `{"total_tokens":12,"sessions":{"old-session":{"started_at":"2026-10-12T02:24:27Z","total_tokens":12,"model_tokens":{"gemini-test":12},"model_usage":{"gemini-test":{"total_input_tokens":8,"total_output_tokens":4,"total_tokens":12}}}}}`

	data, err := decodeTokenUsageData([]byte(legacy))
	if err != nil {
		t.Fatalf("decodeTokenUsageData() error = %v", err)
	}
	session := findTokenUsageSession(data["gemini-test"], "old-session")
	if session.SessionID != "old-session" || session.InputTokens != 8 || session.OutputTokens != 4 || session.TotalTokens != 12 {
		t.Errorf("migrated session = %+v", session)
	}
}

func TestEnsureApplicationDirectoriesCreatesExpectedFiles(t *testing.T) {
	directories, err := ensureApplicationDirectories(t.TempDir())
	if err != nil {
		t.Fatalf("ensureApplicationDirectories() error = %v", err)
	}
	for _, directory := range []string{directories.root, directories.input, directories.output} {
		info, err := os.Stat(directory)
		if err != nil {
			t.Errorf("expected directory %s: %v", directory, err)
		} else if !info.IsDir() {
			t.Errorf("%s is not a directory", directory)
		}
	}
	if _, err := os.Stat(directories.tokenUsage); err != nil {
		t.Errorf("expected token usage file: %v", err)
	}
}

func TestSaveUploadedInput(t *testing.T) {
	inputDir := filepath.Join(t.TempDir(), "input")
	filePath, err := saveUploadedInput(inputDir, "../document.txt", []byte("uploaded"))
	if err != nil {
		t.Fatalf("saveUploadedInput() error = %v", err)
	}
	if filepath.Dir(filePath) != inputDir {
		t.Fatalf("saved file outside input directory: %s", filePath)
	}
	data, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}
	if string(data) != "uploaded" {
		t.Errorf("saved content = %q, want uploaded", data)
	}
}
