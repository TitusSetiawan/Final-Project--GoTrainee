package entity

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Users struct {
	Id        int
	Age       int
	Email     string
	Password  string
	Username  string
	Update_at time.Time
}
type UsersResult struct {
	Id       int
	Age      int
	Email    string
	Username string
}

type UsersPostLogin struct {
	Email    string
	Password string
}

type UsersPostLogin1 struct {
	Email    string
	Password string
	Username string
}

type UserToken struct {
	// Username string
	Token string
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

type UserUpdateIn struct {
	Email        string
	Username     string
	UserUpdateAt time.Time
}

type UserUpdateRes struct {
	Id           int
	Email        string
	Username     string
	Age          int
	UserUpdateAt time.Time
}

var MapUser = map[int]Users{}
var JwtKey = []byte("idontknow")
