package controllers

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"fiber-mongo-api/configs"
	"fiber-mongo-api/models"
	"fiber-mongo-api/security"
	"fiber-mongo-api/types"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
	evp "github.com/walkert/go-evp"
	"github.com/zenazn/pkcs7pad"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	jwtSecret = os.Getenv("JWT_SECRET")
)

type (
	MsgLogin types.Login
	MsgToken types.Token
)

var userPass *mongo.Collection = configs.GetCollection(configs.DB, "user_pass")

func Login(c *fiber.Ctx) error {
	var body MsgLogin
	err := c.BodyParser(&body)
	if err != nil {
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "cannot parse json",
			"msg":   err.Error(),
		})
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var user models.UserPass
	defer cancel()

	err = userPass.FindOne(ctx, bson.M{"username": body.Username}).Decode(&user)
	if err != nil {
		c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Bad Credentials",
		})
		return nil
	}

	err = security.VerifyPassword(user.Password, body.Password)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"status": http.StatusBadRequest,
			"error":  "Incorrect Password",
		})
	}

	token, err := createJwtToken(user)
	if err != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "StatusInternalServerError",
			"msg":   err.Error(),
		})
		return nil
	}
	c.Status(fiber.StatusOK).JSON(fiber.Map{
		"access_token":  token.AccessToken,
		"refresh_token": token.RefreshToken,
		"user":          user.Username,
		"message":       "Login successful!",
	})
	return nil
}

func createJwtToken(user models.UserPass) (MsgToken, error) {
	var msgToken MsgToken
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["iss"] = os.Getenv("APP_NAME")
	claims["sub"] = utils.UUIDv4()
	claims["user"] = user.Username
	claims["biz_id"] = user.Business_Id
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	t, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return msgToken, err
	}
	msgToken.AccessToken = t

	refreshToken := jwt.New(jwt.SigningMethodHS256)
	rtClaims := refreshToken.Claims.(jwt.MapClaims)
	rtClaims["iss"] = os.Getenv("APP_NAME")
	rtClaims["sub"] = utils.UUIDv4()
	rtClaims["exp"] = time.Now().Add(time.Hour * 24 * 7).Unix()
	rt, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return msgToken, err
	}
	msgToken.RefreshToken = rt
	return msgToken, nil
}

// encrypt
func GetToken(c *fiber.Ctx) error {
	text := "Hello World!"

	token := CreateToken(text)

	return c.Status(http.StatusOK).JSON(fiber.Map{"token": token, "text": text})
}

// function
func CreateToken(text string) string {
	rawKey := os.Getenv("RAWKEY")
	data := pkcs7pad.Pad([]byte(text), 16) // 1. Pad the plaintext with PKCS#7
	fmt.Println("padded data: ", hex.EncodeToString(data))

	encryptedData := encrypt(rawKey, data)
	fmt.Println("encrypted data: ", encryptedData)

	return encryptedData
}

func encrypt(rawKey string, plainText []byte) string {
	salt := []byte("ABCDEFGH") // hardcoded at the moment

	// Gets key and IV from raw key.
	key, iv := evp.BytesToKeyAES256CBCMD5([]byte(salt), []byte(rawKey))

	// Create new AES cipher block
	block, err := aes.NewCipher(key)
	if err != nil {
		return err.Error()
	}

	cipherText := make([]byte, len(plainText))

	// Encrypt.
	encryptStream := cipher.NewCTR(block, iv)
	encryptStream.XORKeyStream(cipherText, plainText)

	ivHex := hex.EncodeToString(iv)
	encryptedDataHex := hex.EncodeToString([]byte("Salted__")) + hex.EncodeToString(salt) + hex.EncodeToString(cipherText) // 2. Apply the OpenSSL format, hex encode the result
	return ivHex + ":" + encryptedDataHex                                                                                  // 3. Any value for ivHex can be used here, e.g. "00000000000000000000000000000000"
}
