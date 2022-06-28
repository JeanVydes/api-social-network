package main

import (
	"encoding/base64"
	"math/rand"
	"regexp"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

var (
	err error
)

type Map map[string]interface{}

func JSON(c *gin.Context, code int, success bool, message string, data interface{}) {
	exitCode := 0
	if !success {
		exitCode = 0
	}

	c.JSON(code, Message{
		ExitCode: exitCode,
		Message:  message,
		Data:     data,
	})
}

func Abort(c *gin.Context, code int, success bool, message string, data interface{}) {
	exitCode := 0
	if !success {
		exitCode = 0
	}

	c.AbortWithStatusJSON(code, Message{
		ExitCode: exitCode,
		Message:  message,
		Data:     data,
	})
}

func RandomAccountID() int {
	v := rand.Intn(9999999999999-1000000000000) + 1000000000000
	return v
}

func GenerateToken(accountID string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(accountID), bcrypt.DefaultCost)
    if err != nil {
        return "", err
    }

    return base64.StdEncoding.EncodeToString(hash), nil
}

func HashPassword(password string) ([]byte, error) {
	passwordBytes := []byte(password)

	hashedPassword, err := bcrypt.GenerateFromPassword(passwordBytes, bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return hashedPassword, nil
}

func CheckPasswordHash(passwordInput string, hashedPassword []byte) (bool) {
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(passwordInput))
	if err != nil {
		return err == nil
	}

	return true
}

func isAValidEmail(email string) bool {
	emailRegex := "(^[a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+\\.[a-zA-Z0-9-.]+$)"
	if matched, err := regexp.MatchString(emailRegex, email); !matched || err != nil {
		return false
	}

	return true
}