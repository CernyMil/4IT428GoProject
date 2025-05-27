package pkg

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte(os.Getenv("NEWSLETTER_JWT_SECRET"))

// GenerateJWT creates a JWT with arbitrary claims
func GenerateJWT(claimsMap map[string]interface{}, expiresIn time.Duration) (string, error) {
	claims := jwt.MapClaims{}
	for k, v := range claimsMap {
		claims[k] = v
	}
	if expiresIn > 0 {
		claims["exp"] = time.Now().Add(expiresIn).Unix()
	}
	claims["exp"] = time.Now().Add(expiresIn).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// ParseJWT parses and returns all claims
func ParseJWT(tokenString string) (map[string]interface{}, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return nil, errors.New("invalid or expired token")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid claims")
	}
	return claims, nil
}
