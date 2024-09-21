package tokens

import (
	"context"
	"encoding/base64"
	"errors"
	"time"

	env "github.com/VinukaThejana/todoapp/internal/config"
	rdb "github.com/VinukaThejana/todoapp/internal/redis"
	"github.com/golang-jwt/jwt/v5"
	"github.com/oklog/ulid/v2"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// AccessToken is a struct to represent an access token
type AccessToken struct {
	E  *env.Env
	DB *gorm.DB
	R  *redis.Client
}

// NewAccessToken creates a new access token
func NewAccessToken(e *env.Env, db *gorm.DB, r *redis.Client) *AccessToken {
	return &AccessToken{
		E:  e,
		DB: db,
		R:  r,
	}
}

// AccessTokenDetails is a struct to represent the details of an access token
type AccessTokenDetails struct {
	tokendetails
}

// Create creates a new access token
func (at *AccessToken) Create(
	ctx context.Context,
	userID uint,
	refreshTokenJTI string,
	jti ...string,
) (atd *AccessTokenDetails, err error) {
	atd = &AccessTokenDetails{}
	now := time.Now()

	isRefreshTokenCreated := false

	atd.Iat = now.Unix()
	if len(jti) > 0 {
		isRefreshTokenCreated = true
		atd.JTI = jti[0]
	} else {
		atd.JTI = ulid.Make().String()
	}
	atd.ExpiresIn = now.Add(at.E.AccessTokenExpiresIn).Unix()
	atd.Sub = userID

	if !isRefreshTokenCreated {
		val := at.R.Get(ctx, rdb.RefreshTokenKey(refreshTokenJTI)).Val()
		if val == "" {
			return nil, errors.New("refresh token not found")
		}
	}

	privateKey, err := base64.StdEncoding.DecodeString(at.E.AccessTokenPrivateKey)
	if err != nil {
		return nil, err
	}
	key, err := jwt.ParseRSAPrivateKeyFromPEM(privateKey)
	if err != nil {
		return nil, err
	}

	claims := make(jwt.MapClaims)
	claims["sub"] = atd.Sub
	claims["jti"] = atd.JTI
	claims["iat"] = atd.Iat
	claims["nbf"] = atd.Iat
	claims["exp"] = atd.ExpiresIn

	atd.Token, err = jwt.NewWithClaims(
		jwt.SigningMethodRS256,
		claims,
	).SignedString(key)
	if err != nil {
		return nil, err
	}

	if isRefreshTokenCreated {
		return atd, nil
	}

	pipe := at.R.Pipeline()
	pipe.Set(
		ctx,
		rdb.RefreshTokenKey(refreshTokenJTI),
		atd.JTI,
		redis.KeepTTL,
	)
	pipe.Set(
		ctx,
		rdb.AccessTokenKey(atd.JTI),
		userID,
		at.E.AccessTokenExpiresIn,
	)
	_, err = pipe.Exec(ctx)
	if err != nil {
		return nil, err
	}

	return atd, nil
}
