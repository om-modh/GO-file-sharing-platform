// utils/jwtUtils.go
package utils

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("JWT_SECRET")

// GenerateJWT generates a new JWT token for a given userID
func GenerateJWT(userID int) (string, error) {
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "user_id": userID,
        "exp":     time.Now().Add(time.Hour * 1).Unix(),
    })

    return token.SignedString(jwtSecret)
}

// VerifyJWT verifies the provided JWT token and returns the userID if valid
func VerifyJWT(tokenString string) (*jwt.Token, error) {
    return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        return jwtSecret, nil
    })
}

// GetUserIDFromToken extracts the userID from the JWT token in the HTTP request
func GetUserIDFromToken(r *http.Request) (int, error) {
    tokenString := r.Header.Get("Authorization")
    token, err := VerifyJWT(tokenString)
    if err != nil {
        return 0, err
    }

    claims, ok := token.Claims.(jwt.MapClaims)
    if !ok || !token.Valid {
        return 0, err
    }

    userID, ok := claims["user_id"].(float64)
    if !ok {
        return 0, err
    }

    return int(userID), nil
}
