package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/VinukaThejana/go-utils/logger"
	"github.com/VinukaThejana/todoapp/internal/api/grpc"
	"github.com/VinukaThejana/todoapp/internal/api/router"
	env "github.com/VinukaThejana/todoapp/internal/config"
	"github.com/VinukaThejana/todoapp/internal/database"
	"github.com/VinukaThejana/todoapp/internal/enums"
	rdbc "github.com/VinukaThejana/todoapp/internal/redis"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

var (
	e   = &env.Env{}
	db  *gorm.DB
	rdb *redis.Client
	acm *grpc.AuthClientManager
	err error
)

func init() {
	e.Load()
	db = database.Init(e)
	rdb = rdbc.Init(e)
	isProd := e.Environ == string(enums.Prd)

	if isProd {
		log.Logger = log.Output(zerolog.ConsoleWriter{
			Out: os.Stderr,
		})
	}

	authCfg := grpc.ClientConfig{
		Address:     fmt.Sprintf("%s:%s", e.AuthgRPCDomain, e.AuthgRPCPort),
		DialTimeout: 5 * time.Second,
		UseTLS:      isProd,
	}

	acm, err = grpc.NewAuthClientManager(authCfg)
	if err != nil {
		logger.Errorf(fmt.Errorf("failed to create auth client manager: %w", err))
	}
}

func main() {
	defer func() {
		acm.Close()
	}()

	r := router.Init(acm, e, db, rdb)

	server := &http.Server{
		Addr:    ":" + e.APIGatewayPort,
		Handler: r,
	}

	go func() {
		log.Info().Msgf("starting server on port %s", e.APIGatewayPort)
		if err := server.ListenAndServe(); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				log.Fatal().Err(err).Msg("server failed to start")
			}
		}
	}()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch

	log.Info().Msg("shutting down server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	stopped := make(chan struct{})
	go func() {
		if err := server.Shutdown(ctx); err != nil {
			log.Error().Err(err).Msg("server shutdown failed")
		}

		close(stopped)
	}()

	select {
	case <-ctx.Done():
		log.Error().Msg("server shutdown timed out, forcing shutdown")
	case <-stopped:
		log.Info().Msg("server shutdown complete")
	}
}
