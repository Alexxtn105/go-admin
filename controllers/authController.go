package controllers

import (
	"github.com/dgrijalva/jwt-go/v4"
	"github.com/gofiber/fiber/v2"
	"go-admin/database"
	"go-admin/models"
	"golang.org/x/crypto/bcrypt"
	"strconv"
	"time"
)

// Register - регистрация нового пользователя
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

	// Возвращаем структуру user в виде JSON
	return c.JSON(user)
}

// Login - логин пользователя, возвращает токен в случае успеха
func Login(c *fiber.Ctx) error {
	//парсим данные из запроса
	//создаем мапу
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	// ищем в базе данные о пользователе по его email
	var user models.User
	database.DB.Where("email = ?", data["email"]).First(&user)

	// если пользователь не найден, значит выводим ошибку 404 и сообщение
	if user.Id == 0 {
		c.Status(404)
		return c.JSON(fiber.Map{
			"message": "user not found",
		})
	}

	//сравниваем хеш пароля из базы с тем, что пришел в запросе
	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(data["password"])); err != nil {
		// устанавливаем статус
		c.Status(400)
		// и отправляем сообщение
		return c.JSON(fiber.Map{
			"message": "incorrect password",
		})
	}

	// TODO:
	// Код выше лучше переработать для усложенения задачи мошенникам.
	// Нужно выводить информацию о том, что пользователь или пароль неправильные.
	// Пока оставим как есть

	// создаем JWT-токен
	// сперва создадим "клеймо"
	//var tme int64
	//tme = time.Now().Add(time.Hour * 24).Unix()

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    strconv.Itoa(int(user.Id)),                                  // кому принадлежит (ИД пользователя)
		ExpiresAt: jwt.NewTime(float64(time.Now().Add(time.Hour * 24).Unix())), // когда истекает (1 день = 86400 сенкунд)
		//ExpiresAt: jwt.NewTime(86400),         // когда истекает (1 день = 86400 сенкунд)
	})
	// создаем токен на основе "клейма"
	// TODO: СЕКРЕТНЫЙ КЛЮЧ ВМЕСТО secret!!!
	token, err := claims.SignedString([]byte("secret"))

	// если ошибка
	if err != nil {
		// устанавливаем статус
		return c.SendStatus(fiber.StatusInternalServerError) //fiber.StatusInternalServerError = 500

	}

	// Возвращаем структуру token, но не данные пользователя
	return c.JSON(token)
}
