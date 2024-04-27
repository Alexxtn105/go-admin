package main

import (
	"go-admin/database"
	"go-admin/routes"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	//подключаемся к БД
	database.Connect()

	// Инициализируем новое приложение Fiber
	app := fiber.New()

	// настраиваем CORS (в middleware, для кореектной работы фроненда)
	// ВНИМАНИЕ! Это крайне важно, потому что иначе фроненд не получит куки
	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowOrigins:     "http://127.0.0.1, http://localhost:8080", // явно указываем, с какого сайта можно сделать запрос
	}))

	// настраиваем маршруты
	routes.Setup(app)

	// Start the server on port 3000
	log.Fatal(app.Listen(":3000"))
}
