package main

import (
	"fmt"
	"net"
	"os"

	"github.com/VinukaThejana/go-utils/logger"
	env "github.com/VinukaThejana/todoapp/internal/config"
	"github.com/VinukaThejana/todoapp/internal/database"
	"github.com/VinukaThejana/todoapp/internal/enums"
	"github.com/VinukaThejana/todoapp/internal/lib"
	rdbc "github.com/VinukaThejana/todoapp/internal/redis"
	"github.com/VinukaThejana/todoapp/internal/todo"
	pb "github.com/VinukaThejana/todoapp/pkg/todo"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

var e = &env.Env{}
var db *gorm.DB
var rdb *redis.Client

func init() {
	e.Load()
	db = database.Init(e)
	rdb = rdbc.Init(e)

	if e.Environ == string(enums.Dev) {
		log.Logger = log.Output(zerolog.ConsoleWriter{
			Out: os.Stderr,
		})
	}
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", e.TodogRPCPort))
	if err != nil {
		logger.Errorf(fmt.Errorf("failed to listen: %v", err))
	}

	s := grpc.NewServer()
	pb.RegisterTodoServiceServer(s, &todo.Server{
		E:  e,
		DB: db,
		R:  rdb,
	})

	go func() {
		log.Info().Msg(fmt.Sprintf("starting the todo gRPC server on port %s", e.TodogRPCPort))
		if err := s.Serve(lis); err != nil {
			log.Error().Msg(fmt.Sprintf("failed to serve: %v", err))
		}
	}()

	lib.GracefulShutdowngRPC(s)
}
