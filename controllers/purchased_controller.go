package controllers

import (
	"context"
	"fiber-mongo-api/configs"
	"fiber-mongo-api/models"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var purchasedCollection *mongo.Collection = configs.GetCollection(configs.DB, "products")

func CreatePurchased(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var pruchased models.HistoryPurchased
	defer cancel()

	//validate the request body
	if err := c.BodyParser(&pruchased); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"status":  http.StatusBadRequest,
			"message": "error",
			"result":  err.Error(),
		})
	}

	//use the validator library to validate required fields
	if validationErr := validate.Struct(&pruchased); validationErr != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"status":  http.StatusBadRequest,
			"message": "error",
			"result":  validationErr.Error(),
		})
	}

	newPurchased := models.HistoryPurchased{
		Id: primitive.NewObjectID(),

		Created_At: time.Now(),
		Updated_At: time.Now(),
	}

	result, err := productCollection.InsertOne(ctx, newPurchased)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status":  http.StatusInternalServerError,
			"message": "error",
			"result":  err.Error(),
		})
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"status":  http.StatusCreated,
		"message": "success",
		"result":  result,
	})
}

func GetAPurchased(c *fiber.Ctx) error {

	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"status":  http.StatusCreated,
		"message": "success",
		"result":  "result",
	})
}

func EditAPurchased(c *fiber.Ctx) error {

	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"status":  http.StatusCreated,
		"message": "success",
		"result":  "result",
	})
}

func DeleteAPurchased(c *fiber.Ctx) error {

	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"status":  http.StatusCreated,
		"message": "success",
		"result":  "result",
	})
}

func GetAllPurchased(c *fiber.Ctx) error {

	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"status":  http.StatusCreated,
		"message": "success",
		"result":  "result",
	})
}
