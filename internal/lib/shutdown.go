package lib

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

// GracefulShutdowngRPC gracefully shuts down the gRPC server
func GracefulShutdowngRPC(s *grpc.Server, timeout ...time.Duration) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)

	<-ch

	log.Info().Msg("Shutting down gRPC server...")

	t := 5 * time.Second
	if len(timeout) > 0 {
		t = timeout[0]
	}

	ctx, cancel := context.WithTimeout(context.Background(), t)
	defer cancel()

	stopped := make(chan struct{})
	go func() {
		s.GracefulStop()
		close(stopped)
	}()

	select {
	case <-ctx.Done():
		log.Info().Msg("gRPC server shutdown timed out, force shutting down")
		s.Stop()
	case <-stopped:
		log.Info().Msg("gRPC server stopped")
	}
}
