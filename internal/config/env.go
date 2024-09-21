package env

import environ "github.com/VinukaThejana/env"

// Env is the struct that holds the environment variables
type Env struct {
	Domain         string `mapstructure:"DOMAIN" validate:"required"`
	AuthgRPCPort   string `mapstructure:"AUTH_GRPC_PORT" validate:"required"`
	TodogRPCPort   string `mapstructure:"TODO_GRPC_PORT" validate:"required"`
	APIGatewayPort string `mapstructure:"API_GATEWAY_PORT" validate:"required"`
	DatabaseURL    string `mapstructure:"DATABASE_URL" validate:"required"`
	RedisURL       string `mapstructure:"REDIS_URL" validate:"required"`
	RedisPassword  string `mapstructure:"REDIS_PASSWORD" validate:"required"`
	Environ        string `mapstructure:"ENVIRON" validate:"required,oneof=dev stg prod"`
}

func (e *Env) Load(path ...string) {
	environ.Load(e, path...)
}
