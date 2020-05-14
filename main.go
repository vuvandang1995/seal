package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/vuvandang1995/seal/pkg/server"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
	fmt.Fprintf(os.Stdout, "Server exited properly!\n")
}

func run() error {
	done := make(chan os.Signal, 1)

	errSrv := make(chan error)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	srv := server.New()

	srv.Setup()

	httpServer := http.Server{
		Addr:    ":8000",
		Handler: srv,
	}

	go func() {
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errSrv <- err
		}
		fmt.Fprintf(os.Stdout, "Server is shutting down\n")
		errSrv <- nil
	}()

	<-done

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer func() {
		cancel()
	}()

	if err := httpServer.Shutdown(ctx); err != nil {
		return err
	}

	return <-errSrv
}
