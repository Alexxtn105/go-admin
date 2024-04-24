package controllers

import (
	"github.com/gofiber/fiber/v2"
	"go-admin/database"
	"go-admin/middleware"
	"go-admin/models"
	"math"
	"strconv"
)

// AllUsers - возвращает ВСЕХ пользователей
// Например: 				GET http://localhost:3000/api/users
// Например (постранично): 	GET http://localhost:3000/api/users?page=2

func AllUsers(c *fiber.Ctx) error {
	//смотрим разрешения пользователя
	err := middleware.IsAuthorized(c, "users")

	if err != nil {
		return err
	}

	// ПРИКРУТИМ страничный режим
	// берем номер страницы из параметра URL "page", по умолчанию - "1"
	page, _ := strconv.Atoi(c.Query("page", "1"))

	// Используем постраничный вывод.
	// параметр номера страницы - в URL, например, для второй страницы:
	// http://localhost:3000/api/users?page=2
	return c.JSON(models.Paginate(database.DB, &models.User{}, page))
}

func AllUsers_OLD_WORKING(c *fiber.Ctx) error {

	//смотрим разрешения пользователя
	err := middleware.IsAuthorized(c, "users")

	if err != nil {
		return err
	}

	// ПРИКРУТИМ страничный режим
	// берем номер страницы из параметра URL "page", по умолчанию - "1"
	page, _ := strconv.Atoi(c.Query("page", "1"))

	//вводим ограничения для постраничного вывода, если их много
	limit := 5
	//начальная позиция на выбранной странице
	offset := (page - 1) * limit
	//общее количество
	var total int64

	var users []models.User // создаем слайс с данными

	//database.DB.Find(&users) // поиск всех данных в БД

	//Вариант для ролей:
	//делаем предзагрузку таблицы ролей по foreignKey,
	//чтобы корректно отображать данные ролей
	//также вводим ограничение на количество (limit)
	database.DB.Preload("Role").Offset(offset).Limit(limit).Find(&users)

	// получаем количество записей
	database.DB.Model(&models.User{}).Count(&total)

	// Используем постраничный вывод.
	// параметр номера страницы - в URL, например, для второй страницы:
	// http://localhost:3000/api/users?page=2
	return c.JSON(fiber.Map{
		"data": users,
		"meta": fiber.Map{
			"page":      page,
			"total":     total,
			"last_page": math.Floor(float64(int(total)/limit)) + 1, //
		}})

}

// CreateUser - создание пользователя в БД. Не путать с регистрацией пользователя controllers.Register!!!
// Например: POST http://localhost:3000/api/users
func CreateUser(c *fiber.Ctx) error {
	//смотрим разрешения пользователя
	err := middleware.IsAuthorized(c, "users")
	if err != nil {
		return err
	}
	var user models.User

	// парсим данные. Если данные не подходят - выходим с ошибкой
	if err := c.BodyParser(&user); err != nil {
		return err
	}

	// устанавливаем пароль по умолчанию
	// TODO - поменять пароль
	user.SetPassword("1234")

	database.DB.Create(&user)

	return c.JSON(user)
}

// GetUser - получение данных пользователя по его ИД (из параметров URL),
// например: GET http://localhost:3000/api/users/2
func GetUser(c *fiber.Ctx) error {
	//смотрим разрешения пользователя
	err := middleware.IsAuthorized(c, "users")
	if err != nil {
		return err
	}
	// берем параметр из URL
	id, _ := strconv.Atoi(c.Params("id"))

	//создаем новую структуру с заполненным Id, по которому ниже найдем в базе
	user := models.User{
		Id: uint(id),
	}

	// ищем пользователя в базе
	database.DB.Preload("Role").Find(&user)

	//выводим данные в виде JSON
	return c.JSON(user)
}

// UpdateUser - обновление имеющихся данных пользователя по его ИД (из параметров URL),
// например: PUT http://localhost:3000/api/users/2
func UpdateUser(c *fiber.Ctx) error {
	//смотрим разрешения пользователя
	err := middleware.IsAuthorized(c, "users")
	if err != nil {
		return err
	}
	// берем параметр из URL
	id, _ := strconv.Atoi(c.Params("id"))

	//создаем новую структуру с заполненным Id, по которому ниже найдем в базе
	user := models.User{
		Id: uint(id),
	}

	// парсим данные, которые ввел пользователь. Если данные не подходят - выходим с ошибкой
	if err := c.BodyParser(&user); err != nil {
		return err
	}

	// обновляем данные в базе
	database.DB.Model(&user).Updates(user)

	// возвращаем json с измененными данными
	return c.JSON(user)
}

// DeleteUser - удаление пользователя по его ИД (из параметров URL),
// например: DELETE http://localhost:3000/api/users/2
func DeleteUser(c *fiber.Ctx) error {
	//смотрим разрешения пользователя
	err := middleware.IsAuthorized(c, "users")
	if err != nil {
		return err
	}
	// берем параметр из URL
	id, _ := strconv.Atoi(c.Params("id"))

	//создаем новую структуру с заполненным Id, по которому ниже найдем пользователя в базе
	user := models.User{
		Id: uint(id),
	}

	// удаляем пользователя из базы
	database.DB.Delete(&user)

	// возвращаем json с измененными данными
	return nil
}
