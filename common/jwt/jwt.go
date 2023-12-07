package jwt

import "github.com/google/uuid"
import "github.com/golang-jwt/jwt/v5"

// JWT interface for JWTImpl
type JWT interface {
	GenerateToken(username string, userID int64) (string, error)
	ValidateToken() (string, error)
}

// Implement the JWT interface here
type JWTImpl struct {
	publicKey  string
	privateKey string
}

// NewJWT returns a new JWTImpl
func NewJWT() JWT {
	return &JWTImpl{publicKey: "", privateKey: ""}
}

type CustomClaims struct {
	Username  string `json:"username"`
	UserID    int64  `json:"user_id"`
	TokenUUID string `json:"token_uuid"`
}

func (c CustomClaims) GetExpirationTime() (*jwt.NumericDate, error) {
	//TODO implement me
	panic("implement me")
}

func (c CustomClaims) GetIssuedAt() (*jwt.NumericDate, error) {
	//TODO implement me
	panic("implement me")
}

func (c CustomClaims) GetNotBefore() (*jwt.NumericDate, error) {
	//TODO implement me
	panic("implement me")
}

func (c CustomClaims) GetIssuer() (string, error) {
	//TODO implement me
	panic("implement me")
}

func (c CustomClaims) GetSubject() (string, error) {
	//TODO implement me
	panic("implement me")
}

func (c CustomClaims) GetAudience() (jwt.ClaimStrings, error) {
	//TODO implement me
	panic("implement me")
}

func (c CustomClaims) Valid() error {
	return nil
}

// JWT implementation, which generate a token and validate with customClaims
func (j *JWTImpl) GenerateToken(username string, userID int64) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, CustomClaims{
		Username:  username,
		UserID:    userID,
		TokenUUID: uuid.New().String(),
	})
	return token.Raw, nil
}

func (j *JWTImpl) ValidateToken() (string, error) {
	return "", nil
}
