package tokens

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	JWT_SECRET = "GU9GH9342HT9849EFJW"
	toketTTL   = time.Hour * 12
)

func GenerateAccessToken(id int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, GetClaims(id, toketTTL))

	accessToken, err := token.SignedString([]byte(JWT_SECRET))
	if err != nil {
		return "", fmt.Errorf("error while generating access token")
	}

	return accessToken, nil
}

func ParseJWT(token string) (jwt.Claims, error) {
	parsedToken, _ := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(JWT_SECRET), nil
	})

	return parsedToken.Claims, nil
}

func GetClaims(id int, ttl time.Duration) *jwt.MapClaims {
	return &jwt.MapClaims{
		"sub": fmt.Sprint(id),
		"exp": &jwt.NumericDate{Time: time.Now().Add(ttl)},
	}
}
