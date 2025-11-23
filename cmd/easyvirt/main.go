package main

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/fanyang89/easyvirt/ent"
	"github.com/fanyang89/easyvirt/v1/disk"
	"github.com/fanyang89/easyvirt/v1/domain"
	_ "github.com/mattn/go-sqlite3"
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
				&cli.StringFlag{
					Name:  "libvirt-addr",
					Value: "qemu:///system",
				},
				&cli.StringFlag{
					Name:    "dsn",
					Aliases: []string{"d"},
					Value:   "file:ent?mode=memory&cache=shared&_fk=1",
				},
			},
			Action: func(ctx context.Context, cmd *cli.Command) error {
				addr := cmd.String("addr")
				libvirtAddr := cmd.String("libvirt-addr")
				dsn := cmd.String("dsn")

				db, err := ent.Open("sqlite3", dsn)
				if err != nil {
					return fmt.Errorf("ent.Open: %v", err)
				}
				defer func() { _ = db.Close() }()

				// RPC server listen
				listener, err := net.Listen("tcp", addr)
				if err != nil {
					return fmt.Errorf("listen failed, %w", err)
				}
				server := grpc.NewServer()

				// virt service
				virtService, err := domain.NewVirtService(libvirtAddr, db)
				if err != nil {
					return err
				}
				domain.RegisterEasyVirtServer(server, virtService)

				// disk service
				diskService := disk.NewService(db)
				disk.RegisterDiskManagerServer(server, diskService)

				// run server
				go func(ctx context.Context) {
					err := server.Serve(listener)
					if err != nil {
						slog.Error("")
						os.Exit(1)
					}
				}(ctx)

				// wait for stop signal
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
