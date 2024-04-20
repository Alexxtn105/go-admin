package main

import (
	"fmt"
	"github.com/gofiber/fiber/v3"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Зависимости:
// Fiber (web framework):
// go get github.com/gofiber/fiber/v3

// GORM - для связи с БД mysql

//go get -u gorm.io/gorm
// go get -u gorm.io/driver/mysql

import (
	"log"
)

func main() {
	// подключаемся к БД
	db, err := gorm.Open(mysql.Open("root:Sulubun205!@/go-admin"), &gorm.Config{})
	if err != nil {
		panic("Could not connect to the database")
	}
	fmt.Println("Database: ", db)

	// Инициализируем новое приложение Fiber
	app := fiber.New()

	// Define a route for the GET method on the root path '/'
	app.Get("/", func(c fiber.Ctx) error {
		// Send a string response to the client
		return c.SendString("Hello, World 👋!")
	})

	// Start the server on port 3000
	log.Fatal(app.Listen(":3000"))
}
