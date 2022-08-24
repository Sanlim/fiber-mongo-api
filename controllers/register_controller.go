package controllers

import (
	"context"
	"fiber-mongo-api/configs"
	"fiber-mongo-api/models"
	"fiber-mongo-api/security"

	// "fmt"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var userPassCollection *mongo.Collection = configs.GetCollection(configs.DB, "user_pass")

func Register(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var user models.UserPass
	defer cancel()

	//validate the request body
	if err := c.BodyParser(&user); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"status":  http.StatusBadRequest,
			"message": "error",
			"error":   err.Error(),
		})
	}

	err := userPass.FindOne(ctx, bson.M{"username": user.Username}).Decode(&user)
	if err == nil {
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "user already exists",
		})
		return nil
	}

	//use the validator library to validate required fields
	if validationErr := validate.Struct(&user); validationErr != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"status":  http.StatusBadRequest,
			"message": "error",
			"error":   validationErr.Error(),
		})
	}

	password, err := security.EncryptPassword(user.Password)
	if err != nil {
		return c.Status(http.StatusCreated).JSON(fiber.Map{
			"status": "hash password failed",
			"error":  "fail",
		})
	}

	results, err := bizCollection.Find(ctx, bson.M{})
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}
	var biz []models.Business
	defer results.Close(ctx)
	for results.Next(ctx) {
		var singleBusiness models.Business
		if err = results.Decode(&singleBusiness); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		biz = append(biz, singleBusiness)
	}

	newUser := models.UserPass{
		Id:            primitive.NewObjectID(),
		Username:      user.Username,
		Password:      password,
		Business_Name: biz[0].Business_Name,
		Business_Id:   biz[0].Business_Id,
		Created_At:    time.Now().UTC(),
		Updated_At:    time.Now().UTC(),
	}

	result, err := userPassCollection.InsertOne(ctx, newUser)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status":  http.StatusInternalServerError,
			"message": "error",
			"error":   err.Error(),
		})
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"status":  http.StatusCreated,
		"message": "success",
		"data":    result,
	})
}
