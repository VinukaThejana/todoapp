package auth

import (
	"context"
	"errors"

	env "github.com/VinukaThejana/todoapp/internal/config"
	"github.com/VinukaThejana/todoapp/internal/database"
	pb "github.com/VinukaThejana/todoapp/pkg/auth"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

// Server is used to implement auth.AuthServiceServer
type Server struct {
	pb.UnimplementedAuthServiceServer
	E  *env.Env
	DB *gorm.DB
}

// NewServer creates a new auth server
func NewServer(e *env.Env, db *gorm.DB) *Server {
	return &Server{
		E:  e,
		DB: db,
	}
}

// Register is a gRPC endpoint to register a new user
func (s *Server) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	user := &database.User{
		Name:     req.Name,
		Email:    req.Email,
		Username: req.Username,
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return &pb.RegisterResponse{
			Success: false,
			Message: "Failed to hash the password",
		}, status.Error(codes.Internal, "failed to hash the password")
	}
	user.Password = string(hashedPassword)

	err = s.DB.Create(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return &pb.RegisterResponse{
				Success: false,
				Message: "User already exists",
			}, status.Error(codes.AlreadyExists, "user already exists")
		}

		return &pb.RegisterResponse{
			Success: false,
			Message: "Internal server error",
		}, status.Error(codes.Internal, "internal server error")
	}

	return &pb.RegisterResponse{
		Success: true,
		Message: "User registered successfully",
	}, nil
}
