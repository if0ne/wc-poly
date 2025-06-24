package main

import (
	"context"
	"fmt"

	"github.com/ollama/ollama/api"
	"github.com/wasmCloud/wasmCloud/examples/go/providers/custom-template/bindings/exports/yasp/giga/wc"
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

func (h *Handler) Whoami(ctx context.Context, name string) (string, error) {
	switch name {
	case "Паша":
		return "Привет, растоманы", nil
	case "Саид":
		return "Привет вайбкодерам, остальным соболезную", nil
	case "Дима":
		return "Че за х*йня", nil
	default:
		return "Привет, ноунеймы", nil
	}
}

func (h *Handler) Prompt(ctx context.Context, req *wc.PromptRequest) (*wc.PromptResponse, error) {
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

	msg := &wc.MessageResponse{
		Role:     response.Message.Role,
		Content:  response.Message.Content,
		Thinking: response.Message.Thinking,
		Images:   images,
	}

	metrics := &wc.Metrics{
		Total:              uint64(response.Metrics.TotalDuration),
		Load:               uint64(response.Metrics.LoadDuration),
		PromptEvalCount:    uint32(response.Metrics.PromptEvalCount),
		PromptEvalDuration: uint64(response.Metrics.PromptEvalDuration),
		EvalCount:          uint32(response.Metrics.EvalCount),
		EvalDuration:       uint64(response.Metrics.EvalDuration),
	}

	resp := &wc.PromptResponse{
		Model:      req.Model,
		CreateAt:   response.CreatedAt.Format("2006-01-02 15:04:05"),
		Message:    msg,
		DoneReason: response.DoneReason,
		Metrics:    metrics,
	}

	return resp, nil
}
