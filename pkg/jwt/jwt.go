package app

import (
	"paopao-ce-teaching/internal/core/user"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var (
	JWTSecset = "18a6413dc4fe394c66345ebe501b2f26"
	JWTIssuer = "paopao-api"
	JWTExpire = 1 * time.Hour
)

type Claims struct {
	UID      int64  `json:"uid"`
	Username string `json:"username"`
	jwt.StandardClaims
}

func GetJWTSecret() []byte {
	return []byte(JWTSecset)
}

func GenerateToken(User *user.User) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(JWTExpire)
	claims := Claims{
		UID:      User.ID,
		Username: User.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    JWTIssuer + ":" + User.Salt,
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(GetJWTSecret())
	return token, err
}

func ParseToken(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return GetJWTSecret(), nil
	})
	if err != nil {
		return nil, err
	}
	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}

	return nil, err
}
