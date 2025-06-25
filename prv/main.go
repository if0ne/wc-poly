//go:generate wit-bindgen-wrpc go --out-dir bindings --package github.com/wasmCloud/wasmCloud/examples/go/providers/custom-template/bindings wit

package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ollama/ollama/api"
	"github.com/ollama/ollama/cmd"
	"github.com/shirou/gopsutil/v4/process"
	"github.com/spf13/cobra"
	server "github.com/wasmCloud/wasmCloud/examples/go/providers/custom-template/bindings"
	"go.wasmcloud.dev/provider"
)

func main() {
	if _, err := os.UserHomeDir(); err != nil {
		if err := os.Setenv("HOME", "~"); err != nil {
			return
		}
	}

	command := cmd.NewCLI()
	command.Run = nil
	command.RunE = func(_ *cobra.Command, args []string) error {
		go func() {
			if err := cmd.RunServer(nil, nil); err != nil {

			}
		}()

		return nil
	}
	cobra.CheckErr(command.ExecuteContext(context.Background()))

	StartMemoryMonitor(context.Background(), MemoryMonitorConfig{
		ParentPID:  int32(os.Getpid()),
		LimitBytes: 1300 * 1024 * 1024,
		Interval:   2 * time.Second,
		OnLimitExceed: func(proc *process.Process) {
			_ = proc.Kill()
		},
	})

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

type MemoryMonitorConfig struct {
	ParentPID     int32
	LimitBytes    uint64
	Interval      time.Duration
	OnLimitExceed func(proc *process.Process)
}

func StartMemoryMonitor(ctx context.Context, cfg MemoryMonitorConfig) {
	go func() {
		parent, err := process.NewProcess(cfg.ParentPID)
		if err != nil {
			log.Printf("Failed to get parent process: %v", err)
			return
		}

		ticker := time.NewTicker(cfg.Interval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				log.Println("Memory monitor stopped.")
				return
			case <-ticker.C:
				children, err := parent.Children()
				if err != nil {
					log.Printf("Error getting children: %v", err)
					continue
				}

				var totalMemory uint64
				for _, child := range children {
					mem, err := child.MemoryInfo()
					if err != nil {
						continue
					}
					totalMemory += mem.RSS
				}

				if totalMemory > cfg.LimitBytes {
					log.Printf("Memory limit exceeded: %.2f MB", float64(totalMemory)/1024/1024)
					if cfg.OnLimitExceed != nil {
						for _, child := range children {
							cfg.OnLimitExceed(child)
						}
					}
					return
				}
			}
		}
	}()
}
