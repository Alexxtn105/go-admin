package main

import (
	"github.com/gofiber/fiber/v2"
	"go-admin/routes"
)

import (
	"go-admin/database"
	"log"
)

func main() {
	//подключаемся к БД
	database.Connect()

	// Инициализируем новое приложение Fiber
	app := fiber.New()

	// настраиваем маршруты
	routes.Setup(app)

	// Start the server on port 3000
	log.Fatal(app.Listen(":3000"))
}
