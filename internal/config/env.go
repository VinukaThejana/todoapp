package env

import environ "github.com/VinukaThejana/env"

// Env is the struct that holds the environment variables
type Env struct {
	Domain         string `mapstructure:"DOMAIN" validate:"required"`
	AuthgRPCPort   string `mapstructure:"AUTH_GRPC_PORT" validate:"required"`
	TodogRPCPort   string `mapstructure:"TODO_GRPC_PORT" validate:"required"`
	APIGatewayPort string `mapstructure:"API_GATEWAY_PORT" validate:"required"`
	Environ        string `mapstructure:"ENVIRON" validate:"required,oneof=dev test prod"`
}

func (e *Env) Load(path ...string) {
	environ.Load(e, path...)
}
