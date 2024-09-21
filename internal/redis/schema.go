package redis

import (
	"fmt"

	env "github.com/VinukaThejana/todoapp/internal/config"
)

// RefreshTokenKey returns the key for a refresh token
func RefreshTokenKey(jti string) string {
	return fmt.Sprintf("refresh_token:%s", jti)
}

// RefreshTokenTTL returns the TTL for a refresh token
func RefreshTokenTTL(e *env.Env) int {
	return int(e.RefreshTokenExpiresIn.Seconds())
}

// AccessTokenKey returns the key for an access token
func AccessTokenKey(jti string) string {
	return fmt.Sprintf("access_token:%s", jti)
}

// AccessTokenTTL returns the TTL for an access token
func AccessTokenTTL(e *env.Env) int {
	return int(e.AccessTokenExpiresIn.Seconds())
}
