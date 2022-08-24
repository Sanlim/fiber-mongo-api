package helper

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
)

var jwtSecret = os.Getenv("JWT_SECRET")
var paramsSecret = os.Getenv("PARAMS_SECRET")
var iv = []byte{47, 46, 57, 24, 85, 47, 24, 74, 82, 35, 88, 98, 66, 32, 48, 63}

func GetAccessTokenFromHeader(c *fiber.Ctx) string {
	token := strings.Split(c.Get("Authorization"), " ")

	return token[1]
}

func CheckOwnerToken(tokenString, biz_id string) error {
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})
	if err != nil {
		return err
	}
	fmt.Println("token", token)
	var biz interface{}
	// do something with decoded claims
	for key, val := range claims {
		// fmt.Printf("Key: %v, value: %v\n", key, val)
		if key == "biz_id" && val == biz_id {
			biz = val
		}
	}

	if biz == nil {
		return errors.New("this user isn't owner access_token")
	}

	return nil
}

func EncryptParams(params string) (string, error) {
	block, err := aes.NewCipher([]byte(paramsSecret))
	if err != nil {
		return "", err
	}
	plainText := []byte(params)
	cfb := cipher.NewCFBEncrypter(block, iv)
	cipherText := make([]byte, len(plainText))
	cfb.XORKeyStream(cipherText, plainText)
	return base64.StdEncoding.EncodeToString(cipherText), nil
}

func DecryptParams(params string) (string, error) {
	block, err := aes.NewCipher([]byte(paramsSecret))
	if err != nil {
		return "", err
	}
	cipherText, err := base64.StdEncoding.DecodeString(params)
	if err != nil {
		panic(err)
	}
	cfb := cipher.NewCFBDecrypter(block, iv)
	plainText := make([]byte, len(cipherText))
	cfb.XORKeyStream(plainText, cipherText)
	return string(plainText), nil
}
