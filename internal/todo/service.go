// Pacakge todo : Implement the todo service
package todo

import (
	env "github.com/VinukaThejana/todoapp/internal/config"
	pb "github.com/VinukaThejana/todoapp/pkg/todo"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// Server is used to implement todo.TodoServiceServer
type Server struct {
	pb.UnimplementedTodoServiceServer
	E  *env.Env
	DB *gorm.DB
	R  *redis.Client
}

// NewServer creates a new todo server
func NewServer(e *env.Env, db *gorm.DB, r *redis.Client) *Server {
	return &Server{
		E:  e,
		DB: db,
		R:  r,
	}
}
