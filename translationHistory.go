package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type translationHistoryEntry struct {
	Date             time.Time `json:"date"`
	DocumentName     string    `json:"document_name"`
	DocumentLanguage string    `json:"document_language"`
	TargetLanguage   string    `json:"target_language"`
	Result           string    `json:"result"`
}

func ensureTranslationHistoryFile(filePath string) error {
	if err := os.MkdirAll(filepath.Dir(filePath), 0o700); err != nil {
		return fmt.Errorf("không thể tạo thư mục lịch sử dịch: %w", err)
	}
	if _, err := os.Stat(filePath); err == nil {
		return nil
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("không thể kiểm tra tệp lịch sử dịch: %w", err)
	}
	if err := os.WriteFile(filePath, []byte("[]\n"), 0o600); err != nil {
		return fmt.Errorf("không thể tạo tệp lịch sử dịch: %w", err)
	}
	return nil
}

func appendTranslationHistory(filePath string, entry translationHistoryEntry) error {
	if entry.Result != "thành công" && entry.Result != "thất bại" {
		return fmt.Errorf("kết quả dịch không hợp lệ: %q", entry.Result)
	}
	if entry.Date.IsZero() {
		entry.Date = time.Now().UTC()
	}
	if err := ensureTranslationHistoryFile(filePath); err != nil {
		return err
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("không thể đọc lịch sử dịch: %w", err)
	}
	var entries []translationHistoryEntry
	if err := json.Unmarshal(data, &entries); err != nil {
		return fmt.Errorf("không thể giải mã lịch sử dịch: %w", err)
	}
	entries = append(entries, entry)

	encoded, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return fmt.Errorf("không thể mã hóa lịch sử dịch: %w", err)
	}
	if err := os.WriteFile(filePath, append(encoded, '\n'), 0o600); err != nil {
		return fmt.Errorf("không thể lưu lịch sử dịch: %w", err)
	}
	return nil
}

func readTranslationHistory(filePath string) ([]translationHistoryEntry, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("không thể đọc lịch sử dịch: %w", err)
	}

	var entries []translationHistoryEntry
	if err := json.Unmarshal(data, &entries); err != nil {
		return nil, fmt.Errorf("không thể giải mã lịch sử dịch: %w", err)
	}
	return entries, nil
}
