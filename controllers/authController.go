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

	// TODO: УСЛОЖНИТЬ ЗАДАЧУ
	// Код выше лучше переработать для усложенения задачи злоумышленникам.
	// Не нужно конкретизировать, что именно неправильно ввел пользователь, почту или пароль.
	// Пока оставим как есть

	// создаем JWT-токен
	// сперва создадим заявку (claims)
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    strconv.Itoa(int(user.Id)),                                  //  определяет приложение, из которого отправляется токен
		ExpiresAt: jwt.NewTime(float64(time.Now().Add(time.Hour * 24).Unix())), // когда истекает (1 день = 86400 сенкунд)
		//Subject:   "Auth",                                                      // определяет тему токена.
	})

	// создаем токен на основе заявки
	// TODO: ДОБАВИТЬ СЕКРЕТНЫЙ КЛЮЧ ВМЕСТО secret!!!
	token, err := claims.SignedString([]byte("secret"))

	// если ошибка
	if err != nil {
		// устанавливаем статус

		return c.SendStatus(fiber.StatusInternalServerError) //fiber.StatusInternalServerError = 500

	}

	//сохраняем токен в куки
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true, // это необходимо для того, чтобы фронтенд имел доступ к кукам
	}

	// Пихаем куки в контекст.
	// ВНИМАНИЕ! Необходимо настроить CORS в функции main с параметром allowCredentials: true. См.  строку  app.Use(cors.New(...))
	c.Cookie(&cookie)

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

	//получаем токен из кук (проделываем операцию, обратную созданию токена)
	token, err := jwt.ParseWithClaims(cookie, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})

	// добавил это, поскльку если произошел логаут - сервер слетал, из-за того, что токен пуст
	if token == nil {
		return c.JSON(fiber.Map{
			"message": "token is null",
		})
	}

	if err != nil && !token.Valid {
		c.Status(fiber.StatusUnauthorized) //401
		return c.JSON(fiber.Map{
			"message": "unauthenticated",
		})
	}

	// получаем Claims (заявку)
	claims := token.Claims.(*Claims) //здесь осужествляем преобразование завки в наш внутренний тип *Claims,
	// чтобы в дальнейшем можно было получить приложение (claims.Issuer)

	// получаем приложение (Issuer) из claims (заявки). Фактически это ID пользователя
	// user_id := claims.Issuer

	//получаем пользователя из БД по его ID
	var user models.User
	database.DB.Where("id = ?", claims.Issuer).First(&user)

	//возвращаем пользователя
	return c.JSON(user)
}

// Logout - выход пользователя. Фактически, необходимо удалить куки
func Logout(c *fiber.Ctx) error {
	//cookie := c.Cookies("jwt")

	// Нам необхожимо удалить куки аутентифицированного пользователя.
	// Для этого создаем другие куки, но с пустым токеном и датой истечения где-то в прошлом:
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
