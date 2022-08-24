package main

import (
	"fiber-mongo-api/configs"
	// "fiber-mongo-api/middleware"
	"fiber-mongo-api/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/joho/godotenv"
)

var (
	app = fiber.New()
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		panic("Error loading .env file")
	}
	app.Use(requestid.New())
	app.Use(requestid.New(requestid.Config{
		Header: "Test-Service-Header",
		Generator: func() string {
			return utils.UUID()
		},
	}))
	app.Use(logger.New(logger.Config{
		Format:     "${pid} ${status} - ${method} ${path}\n",
		TimeFormat: "02-Jan-2006",
		TimeZone:   "Asia/Bangkok",
	}))
}

func main() {
	app.Use(cors.New())

	//run database
	configs.ConnectDB()

	//routes
	routes.InitApiRouter(app)

	err := app.Listen(":5000")
	if err != nil {
		panic(err)
	}
}
