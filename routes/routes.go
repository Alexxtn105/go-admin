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

	// CRUD для пользователей
	app.Get("/api/users", controllers.AllUsers)
	app.Post("/api/users", controllers.CreateUser)
	app.Get("/api/users/:id", controllers.GetUser)
	app.Put("/api/users/:id", controllers.UpdateUser)
	app.Delete("/api/users/:id", controllers.DeleteUser)

	// CRUD для ролей
	app.Get("/api/roles", controllers.AllRoles)
	app.Post("/api/roles", controllers.CreateRole)
	app.Get("/api/roles/:id", controllers.GetRole)
	app.Put("/api/roles/:id", controllers.UpdateRole)
	app.Delete("/api/roles/:id", controllers.DeleteRole)

	// CRUD для разрешений
	app.Get("/api/permissions", controllers.AllPermissions)
	app.Post("/api/permissions", controllers.CreatePermission)
	app.Get("/api/permissions/:id", controllers.GetPermission)
	app.Put("/api/permissions/:id", controllers.UpdatePermission)
	app.Delete("/api/permissions/:id", controllers.DeletePermission)
}
