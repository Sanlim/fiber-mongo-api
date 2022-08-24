package controllers

import (
	"context"
	"fiber-mongo-api/configs"
	"fiber-mongo-api/helper"
	"fiber-mongo-api/models"
	"fmt"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
	"go.mongodb.org/mongo-driver/bson"
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
		Id:            primitive.NewObjectID(),
		Business_Name: bizInput.BusinessName,
		Business_Id:   utils.UUIDv4(),
		Created_At:    time.Now().UTC(),
		Updated_At:    time.Now().UTC(),
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

func GetABusiness(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	bizId := c.Query("bizId")
	var biz models.Business
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(bizId)

	err := bizCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&biz)

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status":  http.StatusInternalServerError,
			"message": "error",
			"result":  err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"result":  biz,
		"status":  http.StatusOK,
		"message": "success",
	})
}

func EditABusiness(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	bizId := c.Query("bizId")
	var biz models.Business
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(bizId)

	//validate the request body
	if err := c.BodyParser(&biz); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"status":  http.StatusBadRequest,
			"message": "error",
			"result":  err.Error(),
		})
	}

	//use the validator library to validate required fields
	if validationErr := validate.Struct(&biz); validationErr != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"status":  http.StatusBadRequest,
			"message": "error",
			"result":  validationErr.Error(),
		})
	}

	update := bson.M{"name": biz.Business_Name, "updated_at": time.Now()}

	result, err := bizCollection.UpdateOne(ctx, bson.M{"id": objId}, bson.M{"$set": update})

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status":  http.StatusInternalServerError,
			"message": "error",
			"result":  err.Error(),
		})
	}

	//get updated biz details
	var updatedbiz models.Business
	if result.MatchedCount == 1 {
		err := bizCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&updatedbiz)

		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"status":  http.StatusInternalServerError,
				"message": "error",
				"result":  err.Error(),
			})
		}
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"result":  updatedbiz,
		"status":  http.StatusOK,
		"message": "success",
	})
}

func DeleteABusiness(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	bizId := c.Query("bizId")
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(bizId)

	result, err := bizCollection.DeleteOne(ctx, bson.M{"id": objId})
	fmt.Println("result: ", result)

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status":  http.StatusInternalServerError,
			"message": "error",
			"result":  err.Error(),
		})
	}

	if result.DeletedCount < 1 {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"status":  http.StatusNotFound,
			"message": "error",
			"result":  "User with specified ID not found!"},
		)
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"status":  http.StatusOK,
		"message": "success",
		"result":  "User successfully deleted!",
	})
}

func GetAllBusiness(c *fiber.Ctx) error {
	header := helper.GetAccessTokenFromHeader(c)
	biz_id := c.Query("biz_id")
	err := helper.CheckOwnerToken(header, biz_id)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var bizes []models.Business
	var count int64
	defer cancel()

	results, err := bizCollection.Find(ctx, bson.M{})

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status":  http.StatusInternalServerError,
			"message": "error",
			"result":  err.Error(),
		})
	}

	//reading from the db in an optimal way
	defer results.Close(ctx)
	for results.Next(ctx) {
		var singleBiz models.Business
		if err = results.Decode(&singleBiz); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"status":  http.StatusInternalServerError,
				"message": "error",
				"result":  err.Error(),
			})
		}

		bizes = append(bizes, singleBiz)
		count++
	}

	if bizes == nil {
		return c.Status(http.StatusOK).JSON(fiber.Map{
			"results": "No user available",
			"status":  http.StatusOK,
			"message": "success",
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"results": bizes,
		"status":  http.StatusOK,
		"message": "success",
		"total":   count,
	})
}
