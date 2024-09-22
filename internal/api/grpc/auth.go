package grpc

import (
	"context"
	"crypto/tls"
	"sync"

	"github.com/VinukaThejana/todoapp/pkg/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// AuthClientManager is a struct that manages the gRPC client connection to the auth service.
type AuthClientManager struct {
	client auth.AuthServiceClient
	conn   *grpc.ClientConn
	mu     sync.Mutex
}

var (
	authClientManager *AuthClientManager
	authClientOnce    sync.Once
)

// NewAuthClientManager creates a new AuthClientManager.
func NewAuthClientManager(cfg ClientConfig) (*AuthClientManager, error) {
	var err error
	authClientOnce.Do(func() {
		var conn *grpc.ClientConn
		ctx, cancel := context.WithTimeout(context.Background(), cfg.DialTimeout)
		defer cancel()

		opts := []grpc.DialOption{
			grpc.WithBlock(),
		}

		if cfg.UseTLS {
			opts = append(
				opts,
				grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{})),
			)
		} else {
			opts = append(opts, grpc.WithInsecure())
		}

		conn, err = grpc.DialContext(ctx, cfg.Address, opts...)
		if err != nil {
			return
		}

		authClientManager = &AuthClientManager{
			client: auth.NewAuthServiceClient(conn),
			conn:   conn,
		}
	})
	if err != nil {
		return nil, err
	}

	return authClientManager, nil
}

// Client returns the gRPC client connection to the auth service.
func (m *AuthClientManager) Client() auth.AuthServiceClient {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.client
}

// Close closes the gRPC client connection to the auth service.
func (m *AuthClientManager) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.conn != nil {
		return m.conn.Close()
	}

	return nil
}
