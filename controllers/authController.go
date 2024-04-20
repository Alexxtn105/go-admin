package controllers

import (
	"github.com/gofiber/fiber/v2"
	"go-admin/database"
	"go-admin/models"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *fiber.Ctx) error {
	//парсим данные из запроса
	//создаем мапу
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	//проверяем, совпадают ли пароли
	if data["password"] != data["password_confirm"] {
		// если пароли не совпадают, устанавливаем статус 400
		c.Status(400)
		//и выдаем сообщение клиенту (в виде JSON)
		return c.JSON(fiber.Map{"message": "passwords do not match"})
	}

	// хешируем пароль (
	password, _ := bcrypt.GenerateFromPassword([]byte(data["password"]), 14)

	// создаем структуру пользователя данными из принятого JSON
	user := models.User{
		FirstName: data["first_name"],
		LastName:  data["last_name"],
		Email:     data["email"],
		Password:  password,
	}

	//добавляем пользователя в базу
	database.DB.Create(&user)
	// Отправляем ответ клиенту (структуру в виде JSON)
	return c.JSON(user)
}
