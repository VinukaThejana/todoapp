package tokens

import (
	"context"
	"encoding/base64"
	"errors"
	"time"

	env "github.com/VinukaThejana/todoapp/internal/config"
	"github.com/VinukaThejana/todoapp/internal/database"
	rdb "github.com/VinukaThejana/todoapp/internal/redis"
	"github.com/golang-jwt/jwt/v5"
	"github.com/oklog/ulid/v2"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// RefreshToken is a struct to represent a refresh token
type RefreshToken struct {
	E  *env.Env
	DB *gorm.DB
	R  *redis.Client
}

// RefreshTokenDetails is a struct to represent the details of a refresh token
type RefreshTokenDetails struct {
	tokendetails
	AccessTokenJTI string
}

// NewRefreshToken creates a new refresh token
func NewRefreshToken(e *env.Env, db *gorm.DB, rdb *redis.Client) *RefreshToken {
	return &RefreshToken{
		E:  e,
		DB: db,
		R:  rdb,
	}
}

// Create creates a new refresh token
func (rt *RefreshToken) Create(ctx context.Context, userID uint) (rtd *RefreshTokenDetails, err error) {
	rtd = &RefreshTokenDetails{}
	now := time.Now()

	rtd.Iat = now.Unix()
	rtd.JTI = ulid.Make().String()
	rtd.AccessTokenJTI = ulid.Make().String()
	rtd.ExpiresIn = now.Add(rt.E.RefreshTokenExpiresIn).Unix()
	rtd.Sub = userID

	privateKey, err := base64.StdEncoding.DecodeString(rt.E.RefreshTokenPrivateKey)
	if err != nil {
		return nil, err
	}
	key, err := jwt.ParseRSAPrivateKeyFromPEM(privateKey)
	if err != nil {
		return nil, err
	}

	claims := make(jwt.MapClaims)
	claims["sub"] = rtd.Sub
	claims["jti"] = rtd.JTI
	claims["iat"] = rtd.Iat
	claims["nbf"] = rtd.Iat
	claims["exp"] = rtd.ExpiresIn

	rtd.Token, err = jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(key)
	if err != nil {
		return nil, err
	}

	session := database.Session{
		ID:        rtd.JTI,
		LoginAt:   now,
		ExpiresAt: rtd.ExpiresIn,
		UserID:    userID,
	}
	err = rt.DB.Create(&session).Error
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			err = rt.DB.Delete(session).Error
			if err == nil {
				go func() {
					val := rt.R.Get(ctx, rdb.RefreshTokenKey(rtd.JTI)).Val()
					pipe := rt.R.Pipeline()
					pipe.Del(ctx, rdb.RefreshTokenKey(rtd.JTI))
					pipe.Del(ctx, rdb.AccessTokenKey(val))
					pipe.Exec(ctx)
				}()
				err = rt.DB.Create(&session).Error
			}
		}
		if err != nil {
			return nil, err
		}
	}

	pipe := rt.R.Pipeline()
	pipe.Set(
		ctx,
		rdb.RefreshTokenKey(rtd.JTI),
		rtd.AccessTokenJTI,
		rt.E.RefreshTokenExpiresIn,
	)
	pipe.Set(
		ctx,
		rdb.AccessTokenKey(rtd.AccessTokenJTI),
		userID,
		rt.E.AccessTokenExpiresIn,
	)
	_, err = pipe.Exec(ctx)
	if err != nil {
		return nil, err
	}

	return rtd, nil
}

// Validate validates a refresh token
func (rt *RefreshToken) Validate(ctx context.Context, token string) (rtd *RefreshTokenDetails, err error) {
	rtd = &RefreshTokenDetails{}

	publicKey, err := base64.StdEncoding.DecodeString(rt.E.RefreshTokenPublicKey)
	if err != nil {
		return nil, err
	}
	key, err := jwt.ParseRSAPublicKeyFromPEM(publicKey)
	if err != nil {
		return nil, err
	}

	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errors.New("unexpected signing method")
		}

		return key, nil
	})

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok || !parsedToken.Valid {
		return nil, errors.New("invalid token")
	}

	rtd.Sub = uint(claims["sub"].(float64))
	rtd.JTI = claims["jti"].(string)
	rtd.Iat = int64(claims["iat"].(float64))
	rtd.ExpiresIn = int64(claims["exp"].(float64))

	val := rt.R.Get(ctx, rdb.RefreshTokenKey(rtd.JTI)).Val()
	if val == "" {
		return nil, errors.New("invalid token")
	}

	return rtd, nil
}
