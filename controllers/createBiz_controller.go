package controllers

import (
	"context"
	"fiber-mongo-api/configs"
	"fiber-mongo-api/models"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var bizCollection *mongo.Collection = configs.GetCollection(configs.DB, "biz")

func CreateBusiness(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var bizInput struct {
		BusinessName string `json:"biz_name"`
	}
	c.BodyParser(&bizInput)
	defer cancel()

	newBusiness := models.Business{
		Id:           primitive.NewObjectID(),
		BusinessName: bizInput.BusinessName,
		BusinessId:   utils.UUIDv4(),
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
	}

	result, err := bizCollection.InsertOne(ctx, newBusiness)
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
		"result":  result,
	})
}
