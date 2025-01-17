package controllers

import (
	"github.com/dgrijalva/jwt-go/v4"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"go-admin/database"
	"go-admin/models"
	"go-admin/utils"
	"strconv"
	"time"
)

// Register - регистрация нового пользователя
func Register(c *fiber.Ctx) error {
	//парсим данные из запроса
	//создаем мапу для данных запроса
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	//проверяем, совпадают ли введенные пароли
	if data["password"] != data["password_confirm"] {
		// если пароли не совпадают, устанавливаем статус 400
		c.Status(400)
		//и выдаем сообщение клиенту (в виде JSON)
		return c.JSON(fiber.Map{"message": "passwords do not match"})
	}

	// хешируем пароль
	// старый вариант:
	//password, _ := bcrypt.GenerateFromPassword([]byte(data["password"]), 14)

	// создаем структуру пользователя данными из принятого JSON
	user := models.User{
		FirstName: data["first_name"],
		LastName:  data["last_name"],
		Email:     data["email"],
		RoleId:    1, // присваиваем роль по умолчанию
		//Password:  password,    //см. ниже
	}

	// с использованеим метода из модели ползователя:
	user.SetPassword(data["password"])

	// TODO: проверить, есть ли пользователь с указанным email в базе!!!

	//добавляем пользователя в базу
	database.DB.Create(&user)

	// Возвращаем структуру user в виде JSON
	return c.JSON(user)
}

// Login - логин пользователя, возвращает токен в случае успеха
func Login(c *fiber.Ctx) error {
	//парсим данные из запроса
	//создаем мапу
	log.Info("Trying to login")
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {

		log.Errorf("%s", err)

		return err
	}

	// ищем в базе данные о пользователе по его email
	var user models.User
	database.DB.Where("email = ?", data["email"]).First(&user)
	log.Info("user finded. ID: ", user.Id)
	// если пользователь не найден, значит выводим ошибку 404 и сообщение
	if user.Id == 0 {
		c.Status(404)
		return c.JSON(fiber.Map{
			"message": "user not found",
		})
	}

	//сравниваем хеш пароля из базы с тем, что пришел в запросе
	//	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(data["password"])); err != nil {
	if err := user.ComparePassword(data["password"]); err != nil {
		// устанавливаем статус
		c.Status(400)
		log.Info("incorrect password")

		// и отправляем сообщение
		return c.JSON(fiber.Map{
			"message": "incorrect password",
		})
	}

	// TODO: УСЛОЖНИТЬ ЗАДАЧУ
	// Код выше лучше переработать для усложенения задачи злоумышленникам.
	// Не нужно конкретизировать, что именно неправильно ввел пользователь, почту или пароль.
	// Пока оставим как есть

	// генерацию тоекна вынес в middleware (utils.GenerateJwt)
	token, err := utils.GenerateJwt(strconv.Itoa(int(user.Id)))

	// если ошибка
	if err != nil {
		// устанавливаем статус
		log.Info("Internal Server Error (500)")

		return c.SendStatus(fiber.StatusInternalServerError) // 500
	}

	//сохраняем токен в куки
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true, // это необходимо для того, чтобы фронтенд имел доступ к кукам
	}
	log.Info("Сохраняем куки, длина: ", len(cookie.Value))
	// Пихаем куки в контекст.
	// ВНИМАНИЕ! Необходимо настроить CORS в функции main с параметром allowCredentials: true. См.  строку  app.Use(cors.New(...))
	c.Cookie(&cookie)

	log.Info("Success")

	// возвращаем обычное сообщение в формате JSON
	return c.JSON(fiber.Map{
		"message": "success",
	})

	// Возвращаем структуру token, но не данные пользователя
	//return c.JSON(token)
}

// Claims - структура для заявки jwt
type Claims struct {
	jwt.StandardClaims
}

// User - возвращает аутентифицированного пользователя
func User(c *fiber.Ctx) error {
	// сперва получаем куки
	cookie := c.Cookies("jwt") // берем куки по ключу "jwt"

	// вынес парсинг jwt в middleware (utils)
	userId, _ := utils.ParseJwt(cookie)

	//получаем пользователя из БД по его ID (claims.Issuer)
	var user models.User
	database.DB.Where("id = ?", userId).First(&user)

	//возвращаем пользователя
	return c.JSON(user)
}

// Logout - выход пользователя. Фактически, необходимо удалить куки
func Logout(c *fiber.Ctx) error {
	//cookie := c.Cookies("jwt")

	// Нам необхожимо удалить куки аутентифицированного пользователя.
	// Будем использовать небольшой трюк:
	// создаем другие куки, но с пустым токеном и датой истечения где-то в прошлом:
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    "", //token,
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true, // это необходимо для того, чтобы фронтенд имел доступ к кукам
	}

	//устанавливаем текущие куки
	c.Cookie(&cookie)

	return c.JSON(fiber.Map{
		"message": "logout success",
	})
}

// UpdateInfo - обновление информации о текщем залогиненном пользователе
func UpdateInfo(c *fiber.Ctx) error {
	//парсим данные из запроса
	//создаем мапу для данных запроса
	var data map[string]string
	if err := c.BodyParser(&data); err != nil {
		return err
	}

	// получаем куки
	cookie := c.Cookies("jwt") // берем куки по ключу "jwt"

	// вынес парсинг jwt в middleware (utils)
	id, _ := utils.ParseJwt(cookie)
	userId, _ := strconv.Atoi(id)

	//получаем пользователя из БД по его ID (claims.Issuer)
	user := models.User{
		Id:        uint(userId),
		FirstName: data["first_name"],
		LastName:  data["last_name"],
		Email:     data["email"],
	}
	database.DB.Model(&user).Updates(user)

	//возвращаем пользователя
	return c.JSON(user)
}

// UpdatePassword - обновление пароля о текщего залогиненного пользователя
func UpdatePassword(c *fiber.Ctx) error {
	//парсим данные из запроса
	//создаем мапу для данных запроса
	var data map[string]string
	if err := c.BodyParser(&data); err != nil {
		return err
	}

	//проверяем, совпадают ли введенные НОВЫЕ пароли
	if data["password"] != data["password_confirm"] {
		// если пароли не совпадают, устанавливаем статус 400
		c.Status(400)
		//и выдаем сообщение клиенту (в виде JSON)
		return c.JSON(fiber.Map{"message": "passwords do not match"})
	}

	// получаем куки
	cookie := c.Cookies("jwt") // берем куки по ключу "jwt"

	// вынес парсинг jwt в middleware (utils)
	id, _ := utils.ParseJwt(cookie)
	userId, _ := strconv.Atoi(id)

	user := models.User{
		Id: uint(userId),
	}
	user.SetPassword(data["password"])

	database.DB.Model(&user).Updates(user)

	//возвращаем пользователя
	return c.JSON(user)
}
