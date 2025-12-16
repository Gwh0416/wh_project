package jwts

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type JwtToken struct {
	AccessToken  string
	RefreshToken string
	AccessExp    int64
	RefreshExp   int64
}

func CreateToken(val string, secret string, refreshSecret string, exp time.Duration, refreshExp time.Duration) *JwtToken {
	aExp := time.Now().Add(exp).Unix()
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"token": val,
		"exp":   aExp,
	})
	atoken, err := accessToken.SignedString([]byte(secret))
	if err != nil {
		log.Println(err)
	}

	rExp := time.Now().Add(refreshExp).Unix()
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"token": val,
		"exp":   rExp,
	})
	rtoken, err := refreshToken.SignedString([]byte(refreshSecret))
	if err != nil {
		log.Println(err)
	}

	return &JwtToken{
		AccessToken:  atoken,
		AccessExp:    aExp,
		RefreshExp:   rExp,
		RefreshToken: rtoken,
	}
}

func ParseToken(tokenString string, secret string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		token := claims["token"].(string)
		var exp = int64(claims["exp"].(float64))
		if exp <= time.Now().Unix() {
			return "", errors.New("token已过期")
		}
		return token, nil

	} else {
		return "", err
	}
}
