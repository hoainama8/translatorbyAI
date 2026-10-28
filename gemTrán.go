package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"google.golang.org/genai"
	"google.golang.org/genai/interactions/models/interactions"
	"google.golang.org/genai/interactions/models/operations"
)

func TestTrans(req translationRequest) (*InteractionsJSON, error) {
	apikey := strings.TrimSpace(req.APIKey)
	if apikey == "" {
		return nil, errors.New("Gemini API key is required")
	}
	ctx := context.Background()
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  apikey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return nil, fmt.Errorf("không thể khởi tạo Gemini client: %w", err)
	}
	document, err := base64.StdEncoding.Strict().DecodeString(req.Document)
	if err != nil {
		return nil, fmt.Errorf("dữ liệu tài liệu không phải Base64 hợp lệ: %w", err)
	}
	if len(document) == 0 {
		return nil, errors.New("tài liệu cần dịch không được để trống")
	}

	if strings.TrimSpace(req.SourceLang) == "" {
		req.SourceLang = "tự động phát hiện"
	}
	prompt := BuildTextPrompt(req)
	docFile, err := client.Files.Upload(ctx, bytes.NewReader(document), &genai.UploadFileConfig{
		MIMEType: req.FileType,
	})
	if err != nil {
		return nil, err
	}
	input := interactions.NewInteractionsInput([]interactions.Content{
		interactions.NewContent(interactions.DocumentContent{
			URI:      genai.Ptr(docFile.URI),
			MimeType: interactions.DocumentContentMimeType(docFile.MIMEType).ToPointer(),
		})})
	res, err := client.Interactions.Create(ctx, operations.CreateInteractionRequest{
		Body: operations.NewCreateInteractionRequestBody(interactions.CreateModelInteraction{
			Model:             interactions.Model(req.Model),
			Input:             &input,
			SystemInstruction: &prompt,
		}),
	})
	if err != nil {
		return nil, err
	}
	interactionJSON := NewInteractionsJSON(res.Interaction)
	if outputText := interactionJSON.GetOutputText(); outputText != nil {
		fmt.Println(*outputText)
		return interactionJSON, nil
	}
	return nil, errors.New("Gemini không trả về nội dung văn bản trong interaction")
}
