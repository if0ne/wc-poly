package main

import (
	"context"
	"fmt"

	"github.com/ollama/ollama/api"
	"github.com/wasmCloud/wasmCloud/examples/go/providers/custom-template/bindings/exports/yasp/llm/ollama"
	sdk "go.wasmcloud.dev/provider"
)

type OllamaAdaptor interface {
	Chat(context.Context, *api.ChatRequest, api.ChatResponseFunc) error
	Pull(context.Context, *api.PullRequest, api.PullProgressFunc) error
}

type Handler struct {
	client     OllamaAdaptor
	provider   *sdk.WasmcloudProvider
	linkedFrom map[string]map[string]string
	linkedTo   map[string]map[string]string
}

func (h *Handler) Chat(ctx context.Context, req *ollama.ChatRequest) (*ollama.ChatResponse, error) {
	pullReq := &api.PullRequest{Model: req.Model}
	err := h.client.Pull(ctx, pullReq, func(status api.ProgressResponse) error {
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to pull model: %w", err)
	}

	var messages []api.Message
	for _, m := range req.Messages {
		messages = append(messages, api.Message{
			Role:    m.Role,
			Content: m.Content,
		})
	}

	var stream bool
	var request = api.ChatRequest{
		Model:    req.Model,
		Messages: messages,
		Stream:   &stream,
	}

	var response api.ChatResponse
	fn := func(oresponse api.ChatResponse) error {
		response = oresponse
		return nil
	}

	err = h.client.Chat(ctx, &request, fn)

	if err != nil {
		return nil, fmt.Errorf("chat failed: %w", err)
	}

	var images [][]uint8
	for _, img := range response.Message.Images {
		images = append(images, []byte(img))
	}

	msg := &ollama.ChatMessageResponse{
		Role:     response.Message.Role,
		Content:  response.Message.Content,
		Thinking: response.Message.Thinking,
		Images:   images,
	}

	metrics := &ollama.Metrics{
		Total:              uint64(response.Metrics.TotalDuration),
		Load:               uint64(response.Metrics.LoadDuration),
		PromptEvalCount:    uint32(response.Metrics.PromptEvalCount),
		PromptEvalDuration: uint64(response.Metrics.PromptEvalDuration),
		EvalCount:          uint32(response.Metrics.EvalCount),
		EvalDuration:       uint64(response.Metrics.EvalDuration),
	}

	resp := &ollama.ChatResponse{
		Model:      req.Model,
		CreateAt:   response.CreatedAt.Format("2006-01-02 15:04:05"),
		Message:    msg,
		DoneReason: response.DoneReason,
		Metrics:    metrics,
	}

	return resp, nil
}
