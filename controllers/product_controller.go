package controllers

import (
	"context"
	"fiber-mongo-api/configs"
	"fiber-mongo-api/models"
	"fmt"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var productCollection *mongo.Collection = configs.GetCollection(configs.DB, "products")

func CreateProduct(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var product models.Product
	defer cancel()

	//validate the request body
	if err := c.BodyParser(&product); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"status":  http.StatusBadRequest,
			"message": "error",
			"result":  err.Error(),
		})
	}

	//use the validator library to validate required fields
	if validationErr := validate.Struct(&product); validationErr != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"status":  http.StatusBadRequest,
			"message": "error",
			"result":  validationErr.Error(),
		})
	}

	newProduct := models.Product{
		Id:                  primitive.NewObjectID(),
		Product_Code:        product.Product_Code,
		Product_Name:        product.Product_Name,
		Product_Type:        product.Product_Type,
		Product_Unit:        product.Product_Unit,
		Product_Amount:      product.Product_Amount,
		Product_Price:       product.Product_Price,
		Product_Description: product.Product_Description,
		Created_At:          time.Now(),
		Updated_At:          time.Now(),
	}

	result, err := productCollection.InsertOne(ctx, newProduct)
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

func GetAProduct(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	productId := c.Query("productId")
	var product models.Product
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(productId)

	err := productCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&product)

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status":  http.StatusInternalServerError,
			"message": "error",
			"result":  err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"result":  product,
		"status":  http.StatusOK,
		"message": "success",
	})
}

func EditAProduct(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	productId := c.Query("productId")
	var updateProduct models.Product
	var oldProduct models.Product
	defer cancel()

	//validate the request body
	if err := c.BodyParser(&updateProduct); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"status":  http.StatusBadRequest,
			"message": "error",
			"result":  err.Error(),
		})
	}

	//use the validator library to validate required fields
	if validationErr := validate.Struct(&updateProduct); validationErr != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"status":  http.StatusBadRequest,
			"message": "validator error",
			"result":  validationErr.Error(),
		})
	}

	objId, _ := primitive.ObjectIDFromHex(productId)
	err := productCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&oldProduct)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status":  http.StatusInternalServerError,
			"message": "error",
			"result":  err.Error(),
		})
	}

	if updateProduct.Product_Code == "" {
		updateProduct.Product_Code = oldProduct.Product_Code
	}
	if updateProduct.Product_Name == "" {
		updateProduct.Product_Name = oldProduct.Product_Name
	}
	if updateProduct.Product_Type == "" {
		updateProduct.Product_Type = oldProduct.Product_Type
	}
	if updateProduct.Product_Unit == "" {
		updateProduct.Product_Unit = oldProduct.Product_Unit
	}
	if updateProduct.Product_Amount == 0 {
		updateProduct.Product_Amount = oldProduct.Product_Amount
	}
	if updateProduct.Product_Price == 0 {
		updateProduct.Product_Price = oldProduct.Product_Price
	}
	if updateProduct.Product_Description == "" {
		updateProduct.Product_Description = oldProduct.Product_Description
	}

	update := bson.M{
		"product_code":        updateProduct.Product_Code,
		"product_name":        updateProduct.Product_Name,
		"product_type":        updateProduct.Product_Type,
		"product_unit":        updateProduct.Product_Unit,
		"product_amount":      updateProduct.Product_Amount,
		"product_description": updateProduct.Product_Description,
		"updated_at":          time.Now(),
	}

	result, err := productCollection.UpdateOne(ctx, bson.M{"id": objId}, bson.M{"$set": update})

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status":  http.StatusInternalServerError,
			"message": "error",
			"result":  err.Error(),
		})
	}

	//get updated product details
	var updatedproduct models.Product
	if result.MatchedCount == 1 {
		err := productCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&updatedproduct)

		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"status":  http.StatusInternalServerError,
				"message": "error",
				"result":  err.Error(),
			})
		}
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"result":  updatedproduct,
		"status":  http.StatusOK,
		"message": "success",
	})
}

func DeleteAProduct(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	productId := c.Query("productId")
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(productId)

	result, err := productCollection.DeleteOne(ctx, bson.M{"id": objId})
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
			"message": "fail",
			"result":  "Product with specified ID not found!"},
		)
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"status":  http.StatusOK,
		"message": "success",
		"result":  "Product successfully deleted!",
	})
}

func GetAllProducts(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var productes []models.Product
	var count int64
	defer cancel()

	results, err := productCollection.Find(ctx, bson.M{})

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
		var singleBiz models.Product
		if err = results.Decode(&singleBiz); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"status":  http.StatusInternalServerError,
				"message": "error",
				"result":  err.Error(),
			})
		}

		productes = append(productes, singleBiz)
		count++
	}

	if productes == nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"results": "No product available",
			"status":  http.StatusBadRequest,
			"message": "fail",
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"results": productes,
		"status":  http.StatusOK,
		"message": "success",
		"total":   count,
	})
}
