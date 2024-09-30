// Pacakge todo : Implement the todo service
package todo

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	env "github.com/VinukaThejana/todoapp/internal/config"
	"github.com/VinukaThejana/todoapp/internal/database"
	pb "github.com/VinukaThejana/todoapp/pkg/todo"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

// Create is a gRPC endpoint to create a new todo
// returns Internal, nil
func (s *Server) Create(ctx context.Context, req *pb.CreateRequest) (*pb.CreateResponse, error) {
	userID, err := strconv.ParseUint(req.UserId, 10, 64)
	if err != nil {
		return &pb.CreateResponse{
			Success: false,
		}, status.Error(codes.Internal, "failed to parse user id")
	}

	todo := &database.Todo{
		Title:       req.Title,
		UserID:      uint(userID),
		Completed:   false,
		Content:     req.Content,
		Description: req.Description,
	}

	err = s.DB.Create(&todo).Error
	if err != nil {
		return &pb.CreateResponse{
			Success: false,
		}, status.Error(codes.Internal, "failed to create the todo")
	}

	return &pb.CreateResponse{
		Success: true,
		Message: "Todo created successfully",
	}, nil
}

// Get is a gRPC endpoint to get a todo
// returns Internal, NotFound, nil
func (s *Server) Get(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error) {
	todo := &database.Todo{}

	err := s.DB.Where("id = ?", req.Id).First(&todo).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &pb.GetResponse{
				Success: false,
			}, status.Error(codes.NotFound, "todo not found")
		}

		return &pb.GetResponse{
			Success: false,
		}, status.Error(codes.Internal, "failed to get the todo")
	}

	return &pb.GetResponse{
		Success: true,
		Todo: &pb.Todo{
			Id:          fmt.Sprint(todo.ID),
			Title:       todo.Title,
			Description: todo.Description,
			Content:     todo.Content,
			Completed:   todo.Completed,
			UserId:      fmt.Sprint(todo.UserID),
		},
	}, nil
}

// List is a gRPC endpoint to list all todos
// returns Internal, nil
func (s *Server) List(ctx context.Context, req *pb.ListRequest) (*pb.ListResponse, error) {
	userID, err := strconv.ParseUint(req.UserId, 10, 64)
	if err != nil {
		return &pb.ListResponse{
			Todos: []*pb.Todo{},
		}, status.Error(codes.Internal, "failed to parse user id")
	}

	todos := []*database.Todo{}

	err = s.DB.Where("user_id = ?", userID).Find(&todos).Error
	if err != nil {
		return &pb.ListResponse{
			Todos: []*pb.Todo{},
		}, status.Error(codes.Internal, "failed to get the todos")
	}

	pbTodos := []*pb.Todo{}
	for _, todo := range todos {
		pbTodos = append(pbTodos, &pb.Todo{
			Id:          fmt.Sprint(todo.ID),
			Title:       todo.Title,
			Description: todo.Description,
			Content:     todo.Content,
			Completed:   todo.Completed,
			UserId:      fmt.Sprint(todo.UserID),
		})
	}

	return &pb.ListResponse{
		Todos: pbTodos,
	}, nil
}

// Update is a gRPC endpoint to update a todo
// returns DataLoss, Internal, NotFound, nil
func (s *Server) Update(ctx context.Context, req *pb.UpdateRequest) (*pb.UpdateResponse, error) {
	todoID, err := strconv.ParseUint(req.Id, 10, 64)
	if err != nil {
		return &pb.UpdateResponse{
			Success: false,
		}, status.Error(codes.Internal, "failed to parse todo id")
	}
	userID, err := strconv.ParseUint(req.UserId, 10, 64)
	if err != nil {
		return &pb.UpdateResponse{
			Success: false,
		}, status.Error(codes.Internal, "failed to parse user id")
	}

	todo := &database.Todo{}

	if req.Title != "" {
		todo.Title = req.Title
	}
	if req.Description != "" {
		todo.Description = req.Description
	}
	if req.Content != "" {
		todo.Content = req.Content
	}

	todo.ID = uint(todoID)
	todo.UserID = uint(userID)
	todo.Completed = req.Completed

	err = s.DB.Save(&todo).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &pb.UpdateResponse{
				Success: false,
			}, status.Error(codes.NotFound, "todo not found")
		}

		return &pb.UpdateResponse{
			Success: false,
		}, status.Error(codes.Internal, "failed to update the todo")
	}

	return &pb.UpdateResponse{
		Success: true,
	}, nil
}

// Delete is a gRPC endpoint to delete a todo
// returns Unauthenticated, Internal, nil
func (s *Server) Delete(ctx context.Context, req *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	todoID, err := strconv.ParseUint(req.Id, 10, 64)
	if err != nil {
		return &pb.DeleteResponse{
			Success: false,
		}, status.Error(codes.Internal, "failed to parse todo id")
	}
	userID, err := strconv.ParseUint(req.UserId, 10, 64)
	if err != nil {
		return &pb.DeleteResponse{
			Success: false,
		}, status.Error(codes.Internal, "failed to parse user id")
	}

	todo := &database.Todo{}
	todo.ID = uint(todoID)
	todo.UserID = uint(userID)

	err = s.DB.Delete(&todo).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &pb.DeleteResponse{
				Success: false,
			}, status.Error(codes.Unauthenticated, "you are not authorized to delete this todo")
		}

		return &pb.DeleteResponse{
			Success: false,
		}, status.Error(codes.Internal, "failed to delete the todo")
	}

	return &pb.DeleteResponse{
		Success: true,
	}, nil
}
