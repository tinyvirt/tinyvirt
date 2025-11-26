package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/phsym/console-slog"
	"github.com/urfave/cli/v3"
)

var rootCmd = &cli.Command{
	Name: "tinyvirt",
	Commands: []*cli.Command{
		serverCmd,
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
