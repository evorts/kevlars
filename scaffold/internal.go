/**
 * @Author: steven
 * @Description:
 * @File: internal
 * @Date: 16/01/24 09.13
 */

package scaffold

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"
)

type GracefulEngine interface {
	Start(port string) error
	Shutdown(ctx context.Context) error
}

type GracefulLogger interface {
	Info(i ...interface{})
	Fatal(i ...interface{})
}

func GracefulStop(port int, timeout time.Duration, engine GracefulEngine, logger GracefulLogger) {
	GracefulStopWithContext(context.Background(), port, timeout, engine, logger)
}

func GracefulStopWithContext(ctx context.Context, port int, timeout time.Duration, engine GracefulEngine, logger GracefulLogger) {
	// Start server
	go func() {
		if err := engine.Start(fmt.Sprintf(":%d", port)); err != nil {
			logger.Fatal("shutting down the server", err)
		}
	}()
	// Wait for interrupt signal to gracefully shut down the server with a timeout of 10 seconds.
	// Use a buffered channel to avoid missing signals as recommended for signal.Notify
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	if err := engine.Shutdown(ctx); err != nil {
		logger.Fatal(err)
	}
}
