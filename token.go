package main

import (
	"log"
	"time"
	"math/rand"

	"github.com/dgrijalva/jwt-go"
)

const UserTokenHeader = "X-User-Token-ID"

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
	log.Println("Initialized seed")
}

func createToken(user Credentials) (string, time.Time, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	expires := time.Now().Add(time.Hour * 1)

	claims := make(jwt.MapClaims)
	claims["username"] = user.Username;
	claims["password"] = user.Password;
	claims["exp"] = expires.Unix()
	claims["iat"] = time.Now().Unix()
	token.Claims = claims

	tokenString, err := token.SignedString([]byte(randomString(10)))

	return tokenString, expires, err
}

func randomString(len int) string {
	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		bytes[i] = byte(65 + rand.Intn(90-65))
	}
	return string(bytes)
}
