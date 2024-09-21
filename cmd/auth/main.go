package main

import (
	"fmt"
	"net"
	"os"

	"github.com/VinukaThejana/go-utils/logger"
	"github.com/VinukaThejana/todoapp/internal/auth"
	env "github.com/VinukaThejana/todoapp/internal/config"
	"github.com/VinukaThejana/todoapp/internal/database"
	"github.com/VinukaThejana/todoapp/internal/enums"
	"github.com/VinukaThejana/todoapp/internal/lib"
	pb "github.com/VinukaThejana/todoapp/pkg/auth"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

var e = &env.Env{}
var db *gorm.DB

func init() {
	e.Load()
	db = database.Init(e)

	if e.Environ == string(enums.Dev) {
		log.Logger = log.Output(zerolog.ConsoleWriter{
			Out: os.Stderr,
		})
	}
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", e.AuthgRPCPort))
	if err != nil {
		logger.Errorf(fmt.Errorf("failed to listen: %v", err))
	}

	s := grpc.NewServer()
	pb.RegisterAuthServiceServer(s, &auth.Server{E: e, DB: db})

	go func() {
		log.Info().Msg(fmt.Sprintf("starting the auth gRPC server on port %s", e.AuthgRPCPort))
		if err := s.Serve(lis); err != nil {
			log.Error().Msg(fmt.Sprintf("failed to serve: %v", err))
		}
	}()

	lib.GracefulShutdown(s)
}
