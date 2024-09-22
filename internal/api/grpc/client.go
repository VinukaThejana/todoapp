// Pacakage grpc is a client for establishing a connection to the gRPC server.
package grpc

import "time"

// ClientConfig is a struct that contains the configuration for the gRPC client.
type ClientConfig struct {
	Address     string
	DialTimeout time.Duration
	UseTLS      bool
}
