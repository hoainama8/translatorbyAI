package main

func latestTranslationResult(filePath string, request translationHistoryEntry) (string, bool, error) {
	entries, err := readTranslationHistory(filePath)
	if err != nil {
		return "", false, err
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
