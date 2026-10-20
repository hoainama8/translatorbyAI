package main

import (
	"encoding/json"
	"fmt"
	"os"
)

func latestTranslationResult(filePath string, request translationHistoryEntry) (string, bool, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", false, fmt.Errorf("không thể đọc lịch sử dịch: %w", err)
	}

	var entries []translationHistoryEntry
	if err := json.Unmarshal(data, &entries); err != nil {
		return "", false, fmt.Errorf("không thể giải mã lịch sử dịch: %w", err)
	}
	for index := len(entries) - 1; index >= 0; index-- {
		entry := entries[index]
		if entry.DocumentName == request.DocumentName &&
			entry.DocumentLanguage == request.DocumentLanguage &&
			entry.TargetLanguage == request.TargetLanguage {
			return entry.Result, true, nil
		}
	}
	return "", false, nil
}
