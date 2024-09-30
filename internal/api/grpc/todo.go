package grpc

import (
	"context"
	"crypto/tls"
	"sync"

	"github.com/VinukaThejana/todoapp/pkg/todo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// TodoClientManager is a struct that manages the gRPC client connection to the todo service.
type TodoClientManager struct {
	client todo.TodoServiceClient
	conn   *grpc.ClientConn
	mu     sync.Mutex
}

var (
	todoClientManager *TodoClientManager
	todoClientOnce    sync.Once
)

// NewTodoClientManager creates a new TodoClientManager.
func NewTodoClientManager(cfg ClientConfig) (*TodoClientManager, error) {
	var err error
	todoClientOnce.Do(func() {
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

		todoClientManager = &TodoClientManager{
			client: todo.NewTodoServiceClient(conn),
			conn:   conn,
		}
	})
	if err != nil {
		return nil, err
	}

	return todoClientManager, nil
}

// Client returns the gRPC client connection to the todo service.
func (m *TodoClientManager) Client() todo.TodoServiceClient {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.client
}

// Close closes the gRPC client connection to the todo service.
func (m *TodoClientManager) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.conn != nil {
		return m.conn.Close()
	}

	return nil
}
