package env

import (
	"time"

	environ "github.com/VinukaThejana/env"
)

// Env is the struct that holds the environment variables
type Env struct {
	Domain                 string        `mapstructure:"DOMAIN" validate:"required"`
	AuthgRPCPort           string        `mapstructure:"AUTH_GRPC_PORT" validate:"required"`
	TodogRPCPort           string        `mapstructure:"TODO_GRPC_PORT" validate:"required"`
	APIGatewayPort         string        `mapstructure:"API_GATEWAY_PORT" validate:"required"`
	DatabaseURL            string        `mapstructure:"DATABASE_URL" validate:"required"`
	RedisURL               string        `mapstructure:"REDIS_URL" validate:"required"`
	RedisPassword          string        `mapstructure:"REDIS_PASSWORD" validate:"required"`
	RefreshTokenMaxAge     string        `mapstructure:"REFRESH_TOKEN_MAX_AGE" validate:"required"`
	RefreshTokenPrivateKey string        `mapstructure:"REFRESH_TOKEN_PRIVATE_KEY" validate:"required"`
	RefreshTokenPublicKey  string        `mapstructure:"REFRESH_TOKEN_PUBLIC_KEY" validate:"required"`
	AccessTokenMaxAge      string        `mapstructure:"ACCESS_TOKEN_MAX_AGE" validate:"required"`
	AccessTokenPrivateKey  string        `mapstructure:"ACCESS_TOKEN_PRIVATE_KEY" validate:"required"`
	AccessTokenPublicKey   string        `mapstructure:"ACCESS_TOKEN_PUBLIC_KEY" validate:"required"`
	Environ                string        `mapstructure:"ENVIRON" validate:"required,oneof=dev stg prod"`
	SessionSecret          string        `mapstructure:"SESSION_SECRET" validate:"required"`
	AccessTokenExpiresIn   time.Duration `mapstructure:"ACCESS_TOKEN_EXPIRES_IN" validate:"required"`
	RefreshTokenExpiresIn  time.Duration `mapstructure:"REFRESH_TOKEN_EXPIRES_IN" validate:"required"`
}

func (e *Env) Load(path ...string) {
	environ.Load(e, path...)
}
