package controllers

import (
	"github.com/gofiber/fiber/v2"
	"go-admin/database"
	"go-admin/models"
	"strconv"
)

// AllRoles - возвращает ВСЕХ пользователей
// Например: GET http://localhost:3000/api/roles
func AllRoles(c *fiber.Ctx) error {
	var roles []models.Role  // создаем слайс
	database.DB.Find(&roles) // поиск всех данных в БД
	return c.JSON(roles)     // возвращаем JSON с данными
}

// CreateRole - создание роли в БД.
// Например: POST http://localhost:3000/api/roles
func CreateRole(c *fiber.Ctx) error {
	var role models.Role

	// парсим данные. Если данные не подходят - выходим с ошибкой
	if err := c.BodyParser(&role); err != nil {
		return err
	}

	// создаем запсиь в БД
	database.DB.Create(&role)

	return c.JSON(role)
}

// GetRole - получение данных роли по ее ИД (из параметров URL),
// например: GET http://localhost:3000/api/roles/2
func GetRole(c *fiber.Ctx) error {
	// берем параметр из URL
	id, _ := strconv.Atoi(c.Params("id"))

	//создаем новую структуру с заполненным Id, по которому ниже найдем пользователя в базе
	role := models.Role{
		Id: uint(id),
	}

	// ищем роль в базе
	database.DB.Find(&role)

	//выводим данные в виде JSON
	return c.JSON(role)
}

// UpdateRole - обновление имеющихся данных роли по ее ИД (из параметров URL),
// например: PUT http://localhost:3000/api/roles/2
func UpdateRole(c *fiber.Ctx) error {
	// берем параметр из URL
	id, _ := strconv.Atoi(c.Params("id"))

	//создаем новую структуру с заполненным Id, по которому ниже найдем запись в базе
	role := models.Role{
		Id: uint(id),
	}

	// парсим данные, которые ввел пользователь. Если данные не подходят - выходим с ошибкой
	if err := c.BodyParser(&role); err != nil {
		return err
	}

	// обновляем данные в базе
	database.DB.Model(&role).Updates(role)

	// возвращаем json  с измененными данными
	return c.JSON(role)
}

// DeleteRole - удаление роли по ее ИД (из параметров URL),
// например: DELETE http://localhost:3000/api/roles/2
func DeleteRole(c *fiber.Ctx) error {
	// берем параметр из URL
	id, _ := strconv.Atoi(c.Params("id"))

	//создаем новую структуру пользователя с заполненным Id, по которому ниже найдем пользователя в базе
	role := models.Role{
		Id: uint(id),
	}

	// удаляем пользователя из базы
	database.DB.Delete(&role)

	// возвращаем json  с измененными данными
	return nil
}
