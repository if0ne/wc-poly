package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ollama/ollama/api"
	"github.com/ollama/ollama/openai"
	sdk "go.wasmcloud.dev/provider"
	wrpc "wrpc.io/go"
)

type Handler struct {
	client     *Client
	provider   *sdk.WasmcloudProvider
	linkedFrom map[string]map[string]string
	linkedTo   map[string]map[string]string
}

func (h *Handler) Chat(ctx context.Context, req []uint8) (*wrpc.Result[[]uint8, string], error) {
	var input openai.ChatCompletionRequest

	if err := json.Unmarshal(req, &input); err != nil {
		return nil, fmt.Errorf("invalid JSON input: %w", err)
	}

	pullReq := &api.PullRequest{Model: input.Model}
	err := h.client.Pull(ctx, pullReq, func(status api.ProgressResponse) error {
		return nil
	})

	if err != nil {
		return wrpc.Err[[]uint8]("failed to pull model: " + err.Error()), nil
	}

	var response openai.ChatCompletion
	fn := func(oresponse openai.ChatCompletion) error {
		response = oresponse
		return nil
	}

	err = h.client.Chat(ctx, &input, fn)

	if err != nil {
		return wrpc.Err[[]uint8]("chat failed: " + err.Error()), nil
	}

	respBytes, err := json.Marshal(&response)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal response: %w", err)
	}

	return wrpc.Ok[string](respBytes), nil
}

func (h *Handler) Pull(ctx context.Context, req []uint8) (*wrpc.Result[[]uint8, string], error) {
	var input api.PullRequest

	if err := json.Unmarshal(req, &input); err != nil {
		return nil, fmt.Errorf("invalid JSON input: %w", err)
	}

	go func() {
		err := h.client.Pull(ctx, &input, func(status api.ProgressResponse) error {
			return nil
		})

		if err != nil {
			h.provider.Logger.Error("failed to pull model", "error", err)
		}

		h.provider.Logger.Info("model successfully downloaded", "Model name", input.Model)
	}()

	return wrpc.Ok[string]([]uint8("{}")), nil
}

func (h *Handler) ModelList(ctx context.Context) (*wrpc.Result[[]uint8, string], error) {
	response, err := h.client.List(ctx)

	if err != nil {
		return wrpc.Err[[]uint8]("failed to get list: " + err.Error()), nil
	}

	respBytes, err := json.Marshal(response)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal response: %w", err)
	}

	return wrpc.Ok[string]([]uint8(respBytes)), nil
}

func (h *Handler) Delete(ctx context.Context, req []uint8) (*wrpc.Result[[]uint8, string], error) {
	var input api.DeleteRequest

	if err := json.Unmarshal(req, &input); err != nil {
		return nil, fmt.Errorf("invalid JSON input: %w", err)
	}

	err := h.client.Delete(ctx, &input)

	if err != nil {
		return wrpc.Err[[]uint8]("delete failed:" + err.Error()), nil
	}

	return wrpc.Ok[string]([]uint8("{}")), nil
}
