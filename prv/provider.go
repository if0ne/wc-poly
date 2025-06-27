package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ollama/ollama/api"
	"github.com/ollama/ollama/openai"
	sdk "go.wasmcloud.dev/provider"
)

type Handler struct {
	client     *Client
	provider   *sdk.WasmcloudProvider
	linkedFrom map[string]map[string]string
	linkedTo   map[string]map[string]string
}

func (h *Handler) Chat(ctx context.Context, req string) (string, error) {
	h.provider.Logger.Info("Got raw request", "req", req)

	var input openai.ChatCompletionRequest

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

	var response openai.ChatCompletion
	fn := func(oresponse openai.ChatCompletion) error {
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
