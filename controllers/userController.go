package controllers

import (
	"github.com/gofiber/fiber/v2"
	"go-admin/database"
	"go-admin/models"
	"golang.org/x/crypto/bcrypt"
)

// AllUsers - возвращает ВСЕХ пользователей
func AllUsers(c *fiber.Ctx) error {
	var users []models.User  // создаем слайс с пользователями
	database.DB.Find(&users) // поиск всех пользователей в БД
	return c.JSON(users)     // возвращаем JSON с данными
}

// CreateUser - создание пользователя в БД. Не путать с регистрацией пользователя controllers.Register!!!
func CreateUser(c *fiber.Ctx) error {
	var user models.User

	// парсим данные. Если данных о пользователе нет - выходим с ошибкой
	if err := c.BodyParser(&user); err != nil {
		return err
	}

	// TODO: поменять пароль
	password, _ := bcrypt.GenerateFromPassword([]byte("1234"), 14)

	user.Password = password

	database.DB.Create(&user)

	return c.JSON(user)
}
