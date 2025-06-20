package main

import (
	"context"

	sdk "go.wasmcloud.dev/provider"
)

type Handler struct {
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
