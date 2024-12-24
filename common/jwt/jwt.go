package jwt

import (
	"fmt"
	"time"

	"github.com/google/uuid"

	jwtgo "github.com/golang-jwt/jwt/v5"
	// JWT
)

// JWT interface for JWTImpl
type JWT interface {
	GenerateToken(username string, userID int64) (string, string, error)
	GenerateRefreshToken(userID int64, tokenUID string) (string, error)
	ValidateToken(token string) (CustomClaims, error)
}

// JWTImpl Implement the JWT interface here
type JWTImpl struct {
	publicKey  string
	privateKey string
	secret     []byte
}

// NewJWT returns a new JWTImpl
func NewJWT(secret string) JWT {
	secretByes := []byte(secret)
	return &JWTImpl{publicKey: "", privateKey: "", secret: secretByes}
}

type CustomClaims struct {
	Username  string `json:"username"`
	UserID    int64  `json:"user_id"`
	TokenUUID string `json:"token_uuid"`
	jwtgo.RegisteredClaims
}

// GenerateToken JWT implementation, which generate a token and validate with customClaims
func (j *JWTImpl) GenerateToken(username string, userID int64) (string, string, error) {
	tokenUID := uuid.New().String()
	token := jwtgo.NewWithClaims(jwtgo.SigningMethodHS256, CustomClaims{
		Username:  username,
		UserID:    userID,
		TokenUUID: tokenUID,
		RegisteredClaims: jwtgo.RegisteredClaims{
			ExpiresAt: jwtgo.NewNumericDate(time.Now().Add(time.Minute * 7)),
			Issuer:    "kifiApp",
			IssuedAt:  jwtgo.NewNumericDate(time.Now()),
			NotBefore: jwtgo.NewNumericDate(time.Now()),
			Subject:   "access",
			Audience:  jwtgo.ClaimStrings{"kifiApp"},
			ID:        tokenUID,
		},
	})
	s, err := token.SignedString(j.secret)
	if err != nil {
		return "", "", err
	}
	return s, tokenUID, nil
}

func (j *JWTImpl) GenerateRefreshToken(userID int64, tokenUID string) (string, error) {
	token := jwtgo.NewWithClaims(jwtgo.SigningMethodHS256, CustomClaims{
		UserID:    userID,
		TokenUUID: tokenUID,
		RegisteredClaims: jwtgo.RegisteredClaims{
			ExpiresAt: jwtgo.NewNumericDate(time.Now().Add(time.Hour * 1)),
			Issuer:    "kifiApp",
			Subject:   "refresh",
			ID:        tokenUID,
		},
	})
	s, err := token.SignedString(j.secret)
	if err != nil {
		return "", err
	}
	return s, nil
}

func (j *JWTImpl) ValidateToken(token string) (CustomClaims, error) {
	t, err := jwtgo.ParseWithClaims(token, &CustomClaims{}, func(token *jwtgo.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwtgo.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return j.secret, nil
	})
	if err != nil {
		return CustomClaims{}, err
	}
	if claims, ok := t.Claims.(*CustomClaims); ok && t.Valid {
		return *claims, nil
	} else {
		return CustomClaims{}, fmt.Errorf("invalid token")
	}
}
