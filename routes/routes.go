package routes

// Зависимости:
// Fiber (web framework):
// go get github.com/gofiber/fiber/v3

import (
	"github.com/gofiber/fiber/v2"
	"go-admin/controllers"
)

func Setup(app *fiber.App) {
	// Define a route for the GET method on the root path '/'
	app.Post("/api/register", controllers.Register)

}
