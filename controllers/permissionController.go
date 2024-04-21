package controllers

import (
	"github.com/gofiber/fiber/v2"
	"go-admin/database"
	"go-admin/models"
	"strconv"
)

// AllPermissions - возвращает ВСЕ разрешения
// Например: GET http://localhost:3000/api/permissions
func AllPermissions(c *fiber.Ctx) error {
	var permissions []models.Permission // создаем слайс
	database.DB.Find(&permissions)      // поиск всех данных в БД
	return c.JSON(permissions)          // возвращаем JSON с данными
}

// CreatePermission - создание разрешения в БД.
// Например: POST http://localhost:3000/api/permissions
func CreatePermission(c *fiber.Ctx) error {
	var permission models.Permission

	// парсим данные. Если данные не подходят - выходим с ошибкой
	if err := c.BodyParser(&permission); err != nil {
		return err
	}

	// создаем запсиь в БД
	database.DB.Create(&permission)

	return c.JSON(permission)
}

// GetPermission - получение данных разрешения по его ИД (из параметров URL),
// например: GET http://localhost:3000/api/permissions/2
func GetPermission(c *fiber.Ctx) error {
	// берем параметр из URL
	id, _ := strconv.Atoi(c.Params("id"))

	//создаем новую структуру с заполненным Id, по которому ниже найдем пользователя в базе
	permission := models.Permission{
		Id: uint(id),
	}

	// ищем роль в базе
	database.DB.Find(&permission)

	//выводим данные в виде JSON
	return c.JSON(permission)
}

// UpdatePermission - обновление имеющихся данных разрешения по его ИД (из параметров URL),
// например: PUT http://localhost:3000/api/roles/2
func UpdatePermission(c *fiber.Ctx) error {
	// берем параметр из URL
	id, _ := strconv.Atoi(c.Params("id"))

	//создаем новую структуру с заполненным Id, по которому ниже найдем запись в базе
	permission := models.Permission{
		Id: uint(id),
	}

	// парсим данные, которые ввел пользователь. Если данные не подходят - выходим с ошибкой
	if err := c.BodyParser(&permission); err != nil {
		return err
	}

	// обновляем данные в базе
	database.DB.Model(&permission).Updates(permission)

	// возвращаем json  с измененными данными
	return c.JSON(permission)
}

// DeletePermission - удаление роли по ее ИД (из параметров URL),
// например: DELETE http://localhost:3000/api/roles/2
func DeletePermission(c *fiber.Ctx) error {
	// берем параметр из URL
	id, _ := strconv.Atoi(c.Params("id"))

	//создаем новую структуру с заполненным Id, по которому ниже найдем запись в базе
	permission := models.Permission{
		Id: uint(id),
	}

	// удаляем пользователя из базы
	database.DB.Delete(&permission)

	// возвращаем json  с измененными данными
	return nil
}
