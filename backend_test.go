package main

import (
	"encoding/base64"
	"testing"
)

func TestTestTransRequiresAPIKey(t *testing.T) {
	_, err := TestTrans(translationRequest{APIKey: " \t\n"})
	if err == nil || err.Error() != "Gemini API key is required" {
		t.Fatalf("TestTrans() error = %v, want Gemini API key is required", err)
	}
}

func TestValidateDocumentSize(t *testing.T) {
	tests := []struct {
		name          string
		document      string
		maxBytes      int64
		wantOverLimit bool
	}{
		{
			name:     "at limit",
			document: base64.StdEncoding.EncodeToString([]byte("1234")),
			maxBytes: 4,
		},
		{
			name:          "over limit",
			document:      base64.StdEncoding.EncodeToString([]byte("12345")),
			maxBytes:      4,
			wantOverLimit: true,
		},
		{
			name:          "invalid base64",
			document:      "!",
			maxBytes:      4,
			wantOverLimit: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := validateDocumentSize(test.document, test.maxBytes); got != test.wantOverLimit {
				t.Fatalf("validateDocumentSize() = %t, want %t", got, test.wantOverLimit)
			}
		})
	}
}
