// Package tokens creates and validates JWT tokens
package tokens

type tokendetails struct {
	Token     string
	JTI       string
	ExpiresIn int64
	Iat       int64
	Sub       uint
}
