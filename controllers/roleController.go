package controllers

import (
	"github.com/gofiber/fiber/v2"
	"go-admin/database"
	"go-admin/models"
	"strconv"
)

// RoleCreateDTO - структура для создания роли (DTO - Data Transfer Object)
type RoleCreateDTO struct {
	name        string
	permissions []string // id разрешений

}

// AllRoles - возвращает ВСЕХ пользователей
// Например: GET http://localhost:3000/api/roles
func AllRoles(c *fiber.Ctx) error {
	var roles []models.Role  // создаем слайс
	database.DB.Find(&roles) // поиск всех данных в БД
	return c.JSON(roles)     // возвращаем JSON с данными
}

/*
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
*/

// CreateRole - создание роли в БД. Вариант с таблицей разрешений
// Например: POST http://localhost:3000/api/roles
func CreateRole(c *fiber.Ctx) error {

	var roleDTO fiber.Map

	// парсим данные. Если данные не подходят - выходим с ошибкой
	if err := c.BodyParser(&roleDTO); err != nil {
		return err
	}
	// создаем список любых интерфейсов, в нужные значения преобразуем потом
	list := roleDTO["permissions"].([]interface{})

	// необходимо преобразовать разрешения из строк в id
	//создаем слайс разрешений нужной длины
	permissions := make([]models.Permission, len(list)) //так должно работать

	// бежим по полученным в запросе разрешениям
	for i, permissionId := range list {
		var id int
		switch permissionId.(type) {
		case float64:
			id = int(permissionId.(float64))
		case string:
			id, _ = strconv.Atoi(permissionId.(string))
		case int:
			id = permissionId.(int)
		default:
			id = int(permissionId.(float64))
		}

		// это работает
		//id := int(permissionId.(float64))

		//	id, err := strconv.Atoi(permissionId.(string))
		//	if err != nil {

		//panic(err)
		//	}
		permissions[i] = models.Permission{
			Id: uint(id),
		}
	}

	role := models.Role{
		Name:        roleDTO["name"].(string),
		Permissions: permissions, //пихаем сюда разрешения, полученные в цикле
	}

	// создаем запись в БД
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

	// ищем роль в базе, предварительно загружаем таблицу разрешений ("permissions")
	database.DB.Preload("Permissions").Find(&role)

	//выводим данные в виде JSON
	return c.JSON(role)
}

// UpdateRole - обновление имеющихся данных роли по ее ИД (из параметров URL),
// например: PUT http://localhost:3000/api/roles/2
func UpdateRole(c *fiber.Ctx) error {
	// берем параметр из URL
	id, _ := strconv.Atoi(c.Params("id"))

	//создаем новую структуру с заполненным Id, по которому ниже найдем запись в базе
	/*	role := models.Role{
			Id: uint(id),
		}
	*/

	var roleDTO fiber.Map

	// парсим данные, которые ввел пользователь. Если данные не подходят - выходим с ошибкой
	if err := c.BodyParser(&roleDTO); err != nil {
		return err
	}

	// создаем список любых интерфейсов, в нужные значения преобразуем потом
	list := roleDTO["permissions"].([]interface{})

	// необходимо преобразовать разрешения из строк в id
	//создаем слайс разрешений нужной длины
	permissions := make([]models.Permission, len(list)) //так должно работать

	// бежим по полученным в запросе разрешениям
	for i, permissionId := range list {
		id, _ := strconv.Atoi(permissionId.(string))
		permissions[i] = models.Permission{
			Id: uint(id),
		}
	}

	// сперва удалим старые разрешения
	var result models.Role
	database.DB.Table("role_permissions").Where("role_id = ?", id).Delete(&result)

	// теперь создаем роль с нужными разрешениями
	role := models.Role{
		Id:          uint(id),
		Name:        roleDTO["name"].(string),
		Permissions: permissions, //пихаем сюда разрешения, полученные в цикле
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
