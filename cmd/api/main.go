package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"test-question/cmd"
	"test-question/internal/infra"

	"github.com/pkg/errors"
)

const (
	initResourcesTimeout    = 10 * time.Second
	gracefulShutdownTimeout = 30 * time.Second
)

func main() {
	infraCtx, cancel := context.WithTimeout(context.Background(), initResourcesTimeout)
	defer cancel()

	resources, err := infra.Init(infraCtx)
	if err != nil {
		panic(err)
	}

	srv := &http.Server{
		Addr:              resources.Env.ListenPort,
		Handler:           cmd.SetupRouter(resources),
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		resources.Logger.Info("server started", "addr", srv.Addr)
		if errListen := srv.ListenAndServe(); err != nil && !errors.Is(errListen, http.ErrServerClosed) {
			resources.Logger.Error("server error", "err", errListen)
			os.Exit(1)
		}
	}()

	// graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	resources.Logger.Info("shutdown initiated")

	ctx, cancel := context.WithTimeout(context.Background(), gracefulShutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		resources.Logger.Error("graceful shutdown failed", "err", err)
	}

	resources.Logger.Info("server exited")
}
