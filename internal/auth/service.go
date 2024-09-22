package auth

import (
	"context"
	"errors"

	"github.com/VinukaThejana/todoapp/internal/auth/tokens"
	env "github.com/VinukaThejana/todoapp/internal/config"
	"github.com/VinukaThejana/todoapp/internal/database"
	rdb "github.com/VinukaThejana/todoapp/internal/redis"
	pb "github.com/VinukaThejana/todoapp/pkg/auth"
	"github.com/redis/go-redis/v9"
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
	R  *redis.Client
}

// NewServer creates a new auth server
func NewServer(e *env.Env, db *gorm.DB, r *redis.Client) *Server {
	return &Server{
		E:  e,
		DB: db,
		R:  r,
	}
}

// Register is a gRPC endpoint to register a new user
// returns Internal, AlreadyExists, nil
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

func (s *Server) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	user := &database.User{}
	var err error

	switch loginType := req.Login.(type) {
	case *pb.LoginRequest_Email:
		err = s.DB.Where("email = ?", loginType.Email).First(&user).Error
	case *pb.LoginRequest_Username:
		err = s.DB.Where("username = ?", loginType.Username).First(&user).Error
	default:
		return &pb.LoginResponse{
			Success: false,
			Message: "invalid login type",
		}, status.Error(codes.InvalidArgument, "must provide the username or email")
	}

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &pb.LoginResponse{
				Success: false,
				Message: "User not found",
			}, status.Error(codes.NotFound, "user not found")
		}

		return &pb.LoginResponse{
			Success: false,
			Message: "Internal server error",
		}, status.Error(codes.Internal, "internal server error")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return &pb.LoginResponse{
			Success: false,
			Message: "Invalid password",
		}, status.Error(codes.Unauthenticated, "invalid password")
	}

	rt := tokens.NewRefreshToken(s.E, s.DB, s.R)
	refreshToken, err := rt.Create(ctx, user.ID)
	if err != nil {
		return &pb.LoginResponse{
			Success: false,
			Message: "Failed to create refresh token",
		}, status.Error(codes.Internal, "failed to create refresh token")
	}

	at := tokens.NewAccessToken(s.E, s.DB, s.R)
	accessToken, err := at.Create(ctx, user.ID, refreshToken.AccessTokenJTI)
	if err != nil {
		return &pb.LoginResponse{
			Success: false,
			Message: "Failed to create access token",
		}, status.Error(codes.Internal, "failed to create access token")
	}

	st := tokens.NewSessionToken(s.E, s.DB)
	sessionToken, err := st.Create(ctx, user.ID, user.Email, user.Username, user.Name)
	if err != nil {
		return &pb.LoginResponse{
			Success: false,
			Message: "Failed to create session token",
		}, status.Error(codes.Internal, "failed to create session token")
	}

	return &pb.LoginResponse{
		Success: true,
		Message: "User logged in successfully",
		TokenSet: &pb.TokenSet{
			AccessToken:  accessToken.Token,
			RefreshToken: refreshToken.Token,
			SessionToken: sessionToken.Token,
		},
	}, nil
}

// Refresh is a gRPC endpoint to refresh the access token
func (s *Server) Refresh(ctx context.Context, req *pb.RefreshRequest) (*pb.RefreshResponse, error) {
	rt := tokens.NewRefreshToken(s.E, s.DB, s.R)
	rtd, err := rt.Validate(ctx, req.RefreshToken)
	if err != nil {
		return &pb.RefreshResponse{
			Success: false,
			Message: "Invalid refresh token",
		}, status.Error(codes.Unauthenticated, "invalid refresh token")
	}

	at := tokens.NewAccessToken(s.E, s.DB, s.R)
	accessToken, err := at.Create(ctx, rtd.Sub, rtd.JTI)
	if err != nil {
		return &pb.RefreshResponse{
			Success: false,
			Message: "Failed to create access token",
		}, status.Error(codes.Internal, "failed to create access token")
	}

	return &pb.RefreshResponse{
		Success:     true,
		Message:     "Access token refreshed successfully",
		AccessToken: accessToken.Token,
	}, nil
}

// Logout is a gRPC endpoint to logout the user
func (s *Server) Logout(ctx context.Context, req *pb.LogoutRequest) (*pb.LogoutResponse, error) {
	rt := tokens.NewRefreshToken(s.E, s.DB, s.R)
	rtd, err := rt.Validate(ctx, req.RefreshToken)
	if err != nil {
		return &pb.LogoutResponse{
			Success: false,
			Message: "Invalid refresh token",
		}, status.Error(codes.Unauthenticated, "invalid refresh token")
	}

	val := s.R.Get(ctx, rdb.RefreshTokenKey(rtd.JTI)).Val()
	if val == "" {
		return &pb.LogoutResponse{
			Success: false,
			Message: "Invalid refresh token",
		}, status.Error(codes.Unauthenticated, "invalid refresh token")
	}

	pipe := s.R.Pipeline()
	pipe.Del(ctx, rdb.RefreshTokenKey(rtd.JTI))
	pipe.Del(ctx, rdb.AccessTokenKey(val))
	_, err = pipe.Exec(ctx)
	if err != nil {
		return &pb.LogoutResponse{
			Success: false,
			Message: "Failed to delete tokens",
		}, status.Error(codes.Internal, "failed to delete tokens")
	}

	return &pb.LogoutResponse{
		Success: true,
		Message: "User logged out successfully",
	}, nil
}
