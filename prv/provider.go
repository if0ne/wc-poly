package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ollama/ollama/api"
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

func (h *Handler) Chat(ctx context.Context, req string) (string, error) {
	h.provider.Logger.Info("Got raw request", "req", req)

	var input api.ChatRequest

	json.Valid([]byte(req))

	if err := json.Unmarshal([]byte(req), &input); err != nil {
		return "", fmt.Errorf("invalid JSON input: %w", err)
	}

	h.provider.Logger.Info("Got request", "input", input.Messages[0].Content)

	pullReq := &api.PullRequest{Model: input.Model}
	err := h.client.Pull(ctx, pullReq, func(status api.ProgressResponse) error {
		return nil
	})

	if err != nil {
		return "", fmt.Errorf("failed to pull model: %w", err)
	}

	var response api.ChatResponse
	fn := func(oresponse api.ChatResponse) error {
		response = oresponse
		return nil
	}

	err = h.client.Chat(ctx, &input, fn)

	if err != nil {
		return "", fmt.Errorf("chat failed: %w", err)
	}

	respBytes, err := json.Marshal(&response)
	if err != nil {
		return "", fmt.Errorf("failed to marshal response: %w", err)
	}

	return string(respBytes), nil
}
