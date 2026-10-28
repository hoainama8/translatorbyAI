package main

import (
	"strings"
	"testing"

	"google.golang.org/genai"
)

func TestBacthTransRequiresAPIKey(t *testing.T) {
	_, err := BacthTrans(translationRequest{APIKey: " \t\n"})
	if err == nil || err.Error() != "Gemini API key is required" {
		t.Fatalf("BacthTrans() error = %v, want Gemini API key is required", err)
	}
}

func TestIsBatchUnsupportedError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "model does not support batch",
			err:  genai.APIError{Code: 404, Message: "not supported for batchGenerateContent"},
			want: true,
		},
		{
			name: "batch precondition failed",
			err:  genai.APIError{Code: 400, Message: "Precondition check failed."},
			want: true,
		},
		{
			name: "unrelated API error",
			err:  genai.APIError{Code: 403, Message: "permission denied"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := isBatchUnsupportedError(test.err); got != test.want {
				t.Fatalf("isBatchUnsupportedError() = %t, want %t", got, test.want)
			}
		})
	}
}

func TestBatchResultToInteraction(t *testing.T) {
	batch := &genai.BatchJob{
		Dest: &genai.BatchJobDestination{
			InlinedResponses: []*genai.InlinedResponse{
				{
					Response: &genai.GenerateContentResponse{
						Candidates: []*genai.Candidate{
							{Content: &genai.Content{Parts: []*genai.Part{{Text: "Bản dịch"}}}},
						},
						UsageMetadata: &genai.GenerateContentResponseUsageMetadata{
							PromptTokenCount:     10,
							CandidatesTokenCount: 20,
							TotalTokenCount:      30,
						},
					},
				},
			},
		},
	}

	got, err := batchResultToInteraction(batch, "gemini-test")
	if err != nil {
		t.Fatalf("batchResultToInteraction() error = %v", err)
	}
	if text := got.GetOutputText(); text == nil || *text != "Bản dịch" {
		t.Fatalf("GetOutputText() = %v, want Bản dịch", text)
	}
	if got.GetModel() == nil || string(*got.GetModel()) != "gemini-test" {
		t.Fatalf("GetModel() = %v, want gemini-test", got.GetModel())
	}
	usage := got.GetUsage()
	if usage == nil || usage.GetTotalTokens() == nil || *usage.GetTotalTokens() != 30 {
		t.Fatalf("GetUsage() = %v, want total tokens 30", usage)
	}
}

func TestBatchResultToInteractionReturnsRequestError(t *testing.T) {
	batch := &genai.BatchJob{
		Dest: &genai.BatchJobDestination{
			InlinedResponses: []*genai.InlinedResponse{
				{Error: &genai.JobError{Message: "request rejected"}},
			},
		},
	}

	_, err := batchResultToInteraction(batch, "gemini-test")
	if err == nil || !strings.Contains(err.Error(), "request rejected") {
		t.Fatalf("batchResultToInteraction() error = %v, want request error", err)
	}
}
