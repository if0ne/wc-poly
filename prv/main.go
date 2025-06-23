//go:generate wit-bindgen-wrpc go --out-dir bindings --package github.com/wasmCloud/wasmCloud/examples/go/providers/custom-template/bindings wit

package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ollama/ollama/api"
	"github.com/ollama/ollama/cmd"
	server "github.com/wasmCloud/wasmCloud/examples/go/providers/custom-template/bindings"
	"go.wasmcloud.dev/provider"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	client, err := api.ClientFromEnvironment()
	if err != nil {
		return err
	}
	providerHandler := Handler{
		client:     client,
		linkedFrom: make(map[string]map[string]string),
		linkedTo:   make(map[string]map[string]string),
	}
	p, err := provider.New(
		provider.SourceLinkPut(func(link provider.InterfaceLinkDefinition) error {
			return handleNewSourceLink(&providerHandler, link)
		}),
		provider.TargetLinkPut(func(link provider.InterfaceLinkDefinition) error {
			return handleNewTargetLink(&providerHandler, link)
		}),
		provider.SourceLinkDel(func(link provider.InterfaceLinkDefinition) error {
			return handleDelSourceLink(&providerHandler, link)
		}),
		provider.TargetLinkDel(func(link provider.InterfaceLinkDefinition) error {
			return handleDelTargetLink(&providerHandler, link)
		}),
		provider.HealthCheck(func() string {
			return handleHealthCheck(&providerHandler)
		}),
		provider.Shutdown(func() error {
			return handleShutdown(&providerHandler)
		}),
	)
	if err != nil {
		return err
	}

	providerHandler.provider = p

	providerCh := make(chan error, 1)
	signalCh := make(chan os.Signal, 1)

	go func() {
		if err := cmd.RunServer(nil, nil); err != nil {
			providerHandler.provider.Logger.Info("Failed to start ollama")
		} else {
			providerHandler.provider.Logger.Info("Started ollama")
		}
	}()

	stopFunc, err := server.Serve(p.RPCClient, &providerHandler)
	if err != nil {
		p.Shutdown()
		return err
	}

	go func() {
		err := p.Start()
		providerCh <- err
	}()

	signal.Notify(signalCh, syscall.SIGINT)

	select {
	case err = <-providerCh:
		stopFunc()
		return err
	case <-signalCh:
		p.Shutdown()
		stopFunc()
	}

	return nil
}

func handleNewSourceLink(handler *Handler, link provider.InterfaceLinkDefinition) error {
	handler.provider.Logger.Info("Handling new source link", "link", link)
	handler.linkedTo[link.Target] = link.SourceConfig
	return nil
}

func handleNewTargetLink(handler *Handler, link provider.InterfaceLinkDefinition) error {
	handler.provider.Logger.Info("Handling new target link", "link", link)
	handler.linkedFrom[link.SourceID] = link.TargetConfig
	return nil
}

func handleDelSourceLink(handler *Handler, link provider.InterfaceLinkDefinition) error {
	handler.provider.Logger.Info("Handling del source link", "link", link)
	delete(handler.linkedTo, link.SourceID)
	return nil
}

func handleDelTargetLink(handler *Handler, link provider.InterfaceLinkDefinition) error {
	handler.provider.Logger.Info("Handling del target link", "link", link)
	delete(handler.linkedFrom, link.Target)
	return nil
}

func handleHealthCheck(handler *Handler) string {
	handler.provider.Logger.Info("Handling health check")
	return "provider healthy"
}

func handleShutdown(handler *Handler) error {
	handler.provider.Logger.Info("Handling shutdown")
	clear(handler.linkedFrom)
	clear(handler.linkedTo)
	return nil
}
