package app

import (
	"github.com/golang-jwt/jwt/v4"
	"time"
)

type JwtClaims struct {
	UserId   string `json:"user_id"`
	UserName string `json:"user_name"`
	jwt.RegisteredClaims
}

func GenerateJWT(userId string, userName string, jwtSecret string, expire time.Duration) (string, error) {
	claim := JwtClaims{
		UserId:   userId,
		UserName: userName,
		RegisteredClaims: jwt.RegisteredClaims{
			//过期时间
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expire)),
			//签发时间
			IssuedAt: jwt.NewNumericDate(time.Now()),
			//生效时间
			NotBefore: jwt.NewNumericDate(time.Now()),
			//签发人
			Issuer: "gin-scaffold",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	tokenString, err := token.SignedString([]byte(jwtSecret))
	return tokenString, err
}

func ParseJWT(tokenString string, jwtSecret string) (*JwtClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JwtClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})
	if err != nil {
		return nil, err
	}
	if token != nil {
		claims, ok := token.Claims.(*JwtClaims)
		if ok && token.Valid {
			return claims, nil
		}
	}
	return nil, err
}
