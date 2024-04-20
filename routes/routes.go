package routes

// Зависимости:
// Fiber (web framework):
// go get github.com/gofiber/fiber/v3

import (
	"github.com/gofiber/fiber/v2"
	"go-admin/controllers"
)

func Setup(app *fiber.App) {
	// Определение маршрутов
	app.Post("/api/register", controllers.Register)
	app.Post("/api/login", controllers.Login)

}
