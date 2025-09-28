package types

import (
	"time"

	// "github.com/dgrijalva/jwt-go"÷\
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type RegisterUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
 
type User struct{ 
	Username string `json:"username"`
	PasswordHash string `json:"password"`
}

func NewUser(user RegisterUser) (User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err!= nil {
		return User{}, err
	}

	return User{
		Username: user.Username,
		PasswordHash: string(hashedPassword),
	}, nil
}

func ValidatedPassword (hashedPassword , plainPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	return err == nil
}

func CreateToken (user User ) string {
	now := time.Now()
	validUpto := now.Add(time.Hour * 24).Unix()

	claims := jwt.MapClaims{
		"user": user.Username,
		"exp": validUpto,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := "myPersonalAmitSecret"

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return ""
	}
	return tokenString 
}