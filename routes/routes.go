package routes

// Зависимости:
// Fiber (web framework):
// go get github.com/gofiber/fiber/v3

import (
	"github.com/gofiber/fiber/v2"
	"go-admin/controllers"
	"go-admin/middleware"
)

// Setup - определение маршрутов
func Setup(app *fiber.App) {
	// Определение маршрутов
	// публичные маршруты (доступны неаутентифицированным пользователям)
	app.Post("/api/register", controllers.Register)
	app.Post("/api/login", controllers.Login)

	// маршруты для аутентифицированных пользователей
	// поэтому будем использовать middleware - app.Use()
	app.Use(middleware.IsAuthenticated)

	app.Get("/api/user", controllers.User)
	app.Post("/api/logout", controllers.Logout)
}
