package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"google.golang.org/genai"
	"google.golang.org/genai/interactions/models/interactions"
	"google.golang.org/genai/interactions/models/operations"
)

func TestTrans(req translationRequest) (*InteractionsJSON, error) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	fileover50mb := validateDocumentSize(req.Document, maxGeminiBatchSize)

	if strings.TrimSpace(req.SourceLang) == "" {
		req.SourceLang = "tự động phát hiện"
	}
	prompt := BuildTextPrompt(req)
	var input interactions.InteractionsInput
	if fileover50mb {
		// res, err := TranbyDoc(ctx, client, req, prompt)
		// if err != nil {
		// 	return nil, err
		// }
		readerDocument := bytes.NewReader([]byte(req.Document))
		docFile, err := client.Files.Upload(ctx, readerDocument, &genai.UploadFileConfig{
			MIMEType: req.FileType,
		})
		if err != nil {
			return nil, err
		}

		input = interactions.NewInteractionsInput([]interactions.Content{
			interactions.NewContent(interactions.DocumentContent{
				URI:      genai.Ptr(docFile.URI),
				MimeType: interactions.DocumentContentMimeType(docFile.MIMEType).ToPointer(),
			})})
	} else {
		// res, err := TransbyFile(ctx, client, req, prompt)
		// if err != nil {
		// 	return nil, err
		// }
		readerDocument := bytes.NewReader([]byte(req.Document))
		docFile, err := client.Files.Upload(ctx, readerDocument, &genai.UploadFileConfig{
			MIMEType: req.FileType,
		})
		if err != nil {
			return nil, err
		}

		input = interactions.NewInteractionsInput([]interactions.Content{
			interactions.NewContent(interactions.DocumentContent{
				URI:      genai.Ptr(docFile.URI),
				MimeType: interactions.DocumentContentMimeType(docFile.MIMEType).ToPointer(),
			})})
	}
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
