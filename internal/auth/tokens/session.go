package tokens

import (
	"context"
	"errors"
	"time"

	env "github.com/VinukaThejana/todoapp/internal/config"
	"github.com/golang-jwt/jwt/v5"
	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

// SessionToken is a struct to represent a session token
type SessionToken struct {
	E  *env.Env
	DB *gorm.DB
}

func NewSessionToken(e *env.Env, db *gorm.DB) *SessionToken {
	return &SessionToken{
		E:  e,
		DB: db,
	}
}

// SessionTokenDetails is a struct to represent the details of a session token
type SessionTokenDetails struct {
	tokendetails
	Email    string
	Username string
	Name     string
}

func (st *SessionToken) Create(
	ctx context.Context,
	userID uint,
	email,
	username,
	name string,
) (std *SessionTokenDetails, err error) {
	std = &SessionTokenDetails{}
	now := time.Now()

	std.Iat = now.Unix()
	std.JTI = ulid.Make().String()
	std.ExpiresIn = now.Add(st.E.RefreshTokenExpiresIn).Unix()
	std.Sub = userID

	std.Email = email
	std.Username = username
	std.Name = name

	claims := make(jwt.MapClaims)
	claims["sub"] = userID
	claims["jti"] = std.JTI
	claims["iat"] = std.Iat
	claims["nbf"] = std.Iat
	claims["exp"] = std.ExpiresIn
	claims["email"] = std.Email
	claims["username"] = std.Username
	claims["name"] = std.Name

	std.Token, err = jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	).SignedString([]byte(st.E.SessionSecret))
	if err != nil {
		return nil, err
	}

	return std, nil
}

// Validate validates the session token
func (st *SessionToken) Validate(ctx context.Context, token string) (std *SessionTokenDetails, err error) {
	std = &SessionTokenDetails{}

	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}

		return []byte(st.E.SessionSecret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok || !parsedToken.Valid {
		return nil, errors.New("invalid token")
	}

	std.Sub = uint(claims["sub"].(float64))
	std.JTI = claims["jti"].(string)
	std.Iat = int64(claims["iat"].(float64))
	std.ExpiresIn = int64(claims["exp"].(float64))
	std.Email = claims["email"].(string)
	std.Username = claims["username"].(string)
	std.Name = claims["name"].(string)
	std.Token = token

	return std, nil
}
