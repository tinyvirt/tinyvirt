package main

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/fanyang89/easyvirt/proto"
	ev "github.com/fanyang89/easyvirt/v1"
	"github.com/phsym/console-slog"
	"github.com/urfave/cli/v3"
	"google.golang.org/grpc"
)

var rootCmd = &cli.Command{
	Name: "easyvirt",
	Commands: []*cli.Command{
		{
			Name: "server",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "addr",
					Aliases: []string{"a"},
					Value:   ":50051",
				},
			},
			Action: func(ctx context.Context, cmd *cli.Command) error {
				addr := cmd.String("addr")

				service, err := ev.NewEasyVirtServer()
				if err != nil {
					return err
				}

				listener, err := net.Listen("tcp", addr)
				if err != nil {
					return fmt.Errorf("listen failed, %w", err)
				}

				server := grpc.NewServer()
				proto.RegisterEasyVirtServer(server, service)

				go func(ctx context.Context) {
					err := server.Serve(listener)
					if err != nil {
						slog.Error("")
						os.Exit(1)
					}
				}(ctx)

				<-ctx.Done()
				server.GracefulStop()
				return nil
			},
		},
	},
}

func main() {
	slog.SetDefault(slog.New(console.NewHandler(os.Stderr, &console.HandlerOptions{Level: slog.LevelInfo})))

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	err := rootCmd.Run(ctx, os.Args)
	if err != nil {
		slog.Error("unexpected error", "err", err)
	}
}
