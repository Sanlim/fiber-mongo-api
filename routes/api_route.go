package routes

import (
	"fiber-mongo-api/controllers"
	"fiber-mongo-api/middleware"

	// "fmt"

	"github.com/gofiber/fiber/v2"
)

func InitApiRouter(app *fiber.App) {
	appRoute := app.Group("/api")
	// app.Use(middleware.AuthorizationRequired())
	checkAuth := middleware.AuthorizationRequired()

	// Api
	v1Api := appRoute.Group("/v1")

	// Register
	v1RegisterApi := v1Api.Group("/register") // {URL}/api/v1/register
	v1RegisterApi.Post("", controllers.Register)

	// Login
	v1LoginApi := v1Api.Group("/login") // {URL}/api/v1/login
	v1LoginApi.Post("", controllers.Login)
	v1LoginApi.Get("/token", checkAuth, controllers.GetToken)

	// User
	v1UserApi := v1Api.Group("/user", checkAuth) // {URL}/api/v1/user
	v1UserApi.Post("", controllers.CreateUser)
	v1UserApi.Get("", controllers.GetAUser)
	v1UserApi.Put("", controllers.EditAUser)
	v1UserApi.Delete("", controllers.DeleteAUser)
	v1UserApi.Get("/all", controllers.GetAllUsers)

	// Business
	v1BusinessApi := v1Api.Group("/biz") // {URL}/api/v1/biz
	v1BusinessApi.Post("", controllers.CreateBusiness)
}
