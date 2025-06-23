package main

import (
	"context"

	"github.com/ollama/ollama/api"
	sdk "go.wasmcloud.dev/provider"
)

type OllamaAdaptor interface {
	Chat(context.Context, *api.ChatRequest, api.ChatResponseFunc) error
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

func (h *Handler) Prompt(ctx context.Context, prompt string) (string, error) {
	var stream bool
	var request = api.ChatRequest{
		Model:  "gemma3:1b",
		Stream: &stream,
	}

	var response api.ChatResponse
	fn := func(oresponse api.ChatResponse) error {
		response = oresponse
		return nil
	}

	err := h.client.Chat(ctx, &request, fn)

	if err != nil {
		return "Failed", nil
	}

	return response.Message.Content, nil
}
