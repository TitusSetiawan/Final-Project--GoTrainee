package svr

import (
	"newfinalproject/entity"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func SignIn(uname string) (string, bool) {
	expiredTime := time.Now().Add(30 * time.Minute)
	claims := &entity.Claims{
		Username: uname,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiredTime.Unix(),
		},
	}
	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Create the JWT string
	tokenString, err := token.SignedString(entity.JwtKey)
	if err != nil {
		// If there is an error in creating the JWT return an internal server error
		// w.WriteHeader(http.StatusInternalServerError)
		return "", false
	}
	return tokenString, true
}
