package controllers

import (
	"context"
	"fiber-mongo-api/configs"
	"fiber-mongo-api/helper"
	"fiber-mongo-api/models"
	"fmt"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection = configs.GetCollection(configs.DB, "users")
var validate = validator.New()

func CreateUser(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var user models.User
	defer cancel()

	//validate the request body
	if err := c.BodyParser(&user); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"status":  http.StatusBadRequest,
			"message": "error",
			"result":  &fiber.Map{"data": err.Error()}})
	}

	//use the validator library to validate required fields
	if validationErr := validate.Struct(&user); validationErr != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"status":  http.StatusBadRequest,
			"message": "error",
			"result":  validationErr.Error(),
		})
	}

	newUser := models.User{
		Id:         primitive.NewObjectID(),
		Name:       user.Name,
		Location:   user.Location,
		Title:      user.Title,
		Created_At: time.Now(),
		Updated_At: time.Now(),
	}

	result, err := userCollection.InsertOne(ctx, newUser)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status":  http.StatusInternalServerError,
			"message": "error",
			"result":  err.Error(),
		})

	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"result":  result,
		"status":  http.StatusCreated,
		"message": "success",
	})
}

func GetAUser(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	userId := c.Query("userId")
	var user models.User
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(userId)

	err := userCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&user)

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status":  http.StatusInternalServerError,
			"message": "error",
			"result":  err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"result":  user,
		"status":  http.StatusOK,
		"message": "success",
	})
}

func EditAUser(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	userId := c.Query("userId")
	var user models.User
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(userId)

	//validate the request body
	if err := c.BodyParser(&user); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"status":  http.StatusBadRequest,
			"message": "error",
			"result":  err.Error(),
		})
	}

	//use the validator library to validate required fields
	if validationErr := validate.Struct(&user); validationErr != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"status":  http.StatusBadRequest,
			"message": "error",
			"result":  validationErr.Error(),
		})
	}

	update := bson.M{"name": user.Name, "location": user.Location, "title": user.Title, "updated_at": time.Now()}

	result, err := userCollection.UpdateOne(ctx, bson.M{"id": objId}, bson.M{"$set": update})

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status":  http.StatusInternalServerError,
			"message": "error",
			"result":  err.Error(),
		})
	}

	//get updated user details
	var updatedUser models.User
	if result.MatchedCount == 1 {
		err := userCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&updatedUser)

		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"status":  http.StatusInternalServerError,
				"message": "error",
				"result":  err.Error(),
			})
		}
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"result":  updatedUser,
		"status":  http.StatusOK,
		"message": "success",
	})
}

func DeleteAUser(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	userId := c.Query("userId")
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(userId)

	result, err := userCollection.DeleteOne(ctx, bson.M{"id": objId})
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

func GetAllUsers(c *fiber.Ctx) error {
	header := helper.GetAccessTokenFromHeader(c)
	biz_id := c.Query("biz_id")
	err := helper.CheckOwnerToken(header, biz_id)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var users []models.User
	var count int64
	defer cancel()

	results, err := userCollection.Find(ctx, bson.M{})

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
		var singleUser models.User
		if err = results.Decode(&singleUser); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"status":  http.StatusInternalServerError,
				"message": "error",
				"result":  err.Error(),
			})
		}

		users = append(users, singleUser)
		count++
	}

	if users == nil {
		return c.Status(http.StatusOK).JSON(fiber.Map{
			"results": "No user available",
			"status":  http.StatusOK,
			"message": "success",
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"results": users,
		"status":  http.StatusOK,
		"message": "success",
		"total":   count,
	})
}
