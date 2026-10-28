package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"time"

	"google.golang.org/genai"
	"google.golang.org/genai/interactions/models/interactions"
)

const (
	batchPollInterval = 10 * time.Second
	batchTimeout      = 24 * time.Hour
)

func BacthTrans(req translationRequest) (*InteractionsJSON, error) {
	apikey := strings.TrimSpace(req.APIKey)
	if apikey == "" {
		return nil, errors.New("Gemini API key is required")
	}
	if strings.TrimSpace(req.Model) == "" {
		if len(models) == 0 {
			return nil, errors.New("không có mô hình Gemini mặc định")
		}
		req.Model = models[0]
	}
	if len(req.Document) == 0 {
		return nil, errors.New("tài liệu cần dịch không được để trống")
	}
	if strings.TrimSpace(req.FileType) == "" {
		return nil, errors.New("không xác định được định dạng tài liệu")
	}

	document, err := base64.StdEncoding.Strict().DecodeString(req.Document)
	if err != nil {
		return nil, fmt.Errorf("dữ liệu tài liệu không phải Base64 hợp lệ: %w", err)
	}
	if len(document) == 0 {
		return nil, errors.New("tài liệu cần dịch không được để trống")
	}

	ctx, cancel := context.WithTimeout(context.Background(), batchTimeout)
	defer cancel()

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  apikey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return nil, fmt.Errorf("không thể khởi tạo Gemini client: %w", err)
	}

	docFile, err := client.Files.Upload(ctx, bytes.NewReader(document), &genai.UploadFileConfig{
		MIMEType:    req.FileType,
		DisplayName: safeStoredName(req.FileName),
	})
	if err != nil {
		return nil, fmt.Errorf("không thể tải tài liệu lên Gemini: %w", err)
	}
	if docFile == nil || docFile.URI == "" {
		return nil, errors.New("Gemini không trả về URI của tài liệu đã tải lên")
	}

	if strings.TrimSpace(req.SourceLang) == "" {
		req.SourceLang = "tự động phát hiện"
	}
	prompt := BuildTextPrompt(req)
	source := &genai.BatchJobSource{
		InlinedRequests: []*genai.InlinedRequest{
			{
				Contents: []*genai.Content{
					{
						Role: "user",
						Parts: []*genai.Part{
							genai.NewPartFromURI(docFile.URI, docFile.MIMEType),
						},
					},
				},
				Config: &genai.GenerateContentConfig{
					SystemInstruction: &genai.Content{
						Parts: []*genai.Part{genai.NewPartFromText(prompt)},
					},
				},
			},
		},
	}
	batch, err := client.Batches.Create(ctx, req.Model, source, &genai.CreateBatchJobConfig{
		DisplayName: "translation-" + safeStoredName(req.FileName),
	})
	if err != nil {
		if isBatchUnsupportedError(err) {
			result, directErr := TestTrans(req)
			if directErr != nil {
				return nil, fmt.Errorf(
					"Gemini Batch API không hỗ trợ model %q hoặc yêu cầu này (%v); dịch trực tiếp cũng thất bại: %w",
					req.Model,
					err,
					directErr,
				)
			}
			return result, nil
		}
		return nil, fmt.Errorf("không thể tạo Gemini batch job: %w", err)
	}
	if batch == nil || batch.Name == "" {
		return nil, errors.New("Gemini không trả về mã batch job")
	}

	ticker := time.NewTicker(batchPollInterval)
	defer ticker.Stop()

	for {
		switch batch.State {
		case genai.JobStateSucceeded, genai.JobStatePartiallySucceeded:
			return batchResultToInteraction(batch, req.Model)
		case genai.JobStateFailed, genai.JobStateCancelled, genai.JobStateExpired:
			if batch.Error != nil {
				if batch.Error.Message != "" {
					return nil, fmt.Errorf("Gemini batch job kết thúc với trạng thái %s: %s", batch.State, batch.Error.Message)
				}
				if len(batch.Error.Details) > 0 {
					return nil, fmt.Errorf("Gemini batch job kết thúc với trạng thái %s: %s", batch.State, strings.Join(batch.Error.Details, "; "))
				}
			}
			return nil, fmt.Errorf("Gemini batch job kết thúc với trạng thái %s", batch.State)
		}

		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("đã hết thời gian chờ Gemini batch job %s: %w", batch.Name, ctx.Err())
		case <-ticker.C:
		}

		batch, err = client.Batches.Get(ctx, batch.Name, nil)
		if err != nil {
			return nil, fmt.Errorf("không thể kiểm tra trạng thái Gemini batch job %s: %w", batch.Name, err)
		}
		if batch == nil {
			return nil, fmt.Errorf("Gemini không trả về trạng thái của batch job %s", batch.Name)
		}
	}
}

func isBatchUnsupportedError(err error) bool {
	var apiErr genai.APIError
	if !errors.As(err, &apiErr) {
		return false
	}

	switch apiErr.Code {
	case 404:
		return strings.Contains(apiErr.Message, "batchGenerateContent") ||
			strings.Contains(apiErr.Message, "not found")
	case 400:
		return strings.Contains(apiErr.Message, "Precondition check failed")
	default:
		return false
	}
}

func batchResultToInteraction(batch *genai.BatchJob, model string) (*InteractionsJSON, error) {
	if batch.Dest == nil || len(batch.Dest.InlinedResponses) == 0 {
		return nil, errors.New("Gemini batch job không trả về kết quả")
	}

	result := batch.Dest.InlinedResponses[0]
	if result == nil {
		return nil, errors.New("Gemini batch job trả về kết quả trống")
	}
	if result.Error != nil {
		switch {
		case result.Error.Message != "":
			return nil, fmt.Errorf("Gemini không thể dịch tài liệu: %s", result.Error.Message)
		case len(result.Error.Details) > 0:
			return nil, fmt.Errorf("Gemini không thể dịch tài liệu: %s", strings.Join(result.Error.Details, "; "))
		default:
			return nil, errors.New("Gemini không thể dịch tài liệu")
		}
	}
	if result.Response == nil {
		return nil, errors.New("Gemini batch job không trả về phản hồi")
	}

	outputText := result.Response.Text()
	if strings.TrimSpace(outputText) == "" {
		return nil, errors.New("Gemini không trả về nội dung văn bản trong batch job")
	}
	interaction := &InteractionsJSON{
		Model:      genai.Ptr(interactions.Model(model)),
		OutputText: &outputText,
	}
	if usage := result.Response.UsageMetadata; usage != nil && usage.TotalTokenCount > 0 {
		interaction.Usage = &interactions.Usage{
			TotalInputTokens:  genai.Ptr(int(usage.PromptTokenCount)),
			TotalOutputTokens: genai.Ptr(int(usage.CandidatesTokenCount)),
			TotalTokens:       genai.Ptr(int(usage.TotalTokenCount)),
		}
	}

	return interaction, nil
}
