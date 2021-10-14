package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"jwtproject/database"
	"jwtproject/routes"
)

func main(){
	database.Connect()
	app:=fiber.New()


	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
	}))

	routes.SetUp(app)
	app.Listen(":8000")
}