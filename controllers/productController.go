package controllers

import (
	"github.com/gofiber/fiber/v2"
	"go-admin/database"
	"go-admin/models"
	"math"
	"strconv"
)

// AllProducts - возвращает ВСЕ продукты
// Например: 				GET http://localhost:3000/api/products
// Например (постранично): 	GET http://localhost:3000/api/products?page=2

func AllProducts(c *fiber.Ctx) error {
	// ПРИКРУТИМ страничный режим
	// берем номер страницы из параметра URL "page", по умолчанию - "1"
	page, _ := strconv.Atoi(c.Query("page", "1"))
	//вводим ограничения для постраничного вывода пользователей, если их много
	limit := 5
	//начальная позиция на выбранной странице
	offset := (page - 1) * limit
	//общее количество
	var total int64

	var products []models.Product // создаем слайс с продуктами

	//Вариант для ролей:
	//делаем предзагрузку таблицы ролей по foreignKey,
	//чтобы корректно отображать данные пользователей и их ролей
	// также вводим ограничение на количество (limit)
	database.DB.Offset(offset).Limit(limit).Find(&products)

	// получаем количество записей
	database.DB.Model(&models.Product{}).Count(&total)

	// Если нужно выводить постранично - код ниже
	// для постраничного вывода.
	// параметр номера страницы - в URL, например, для второй страницы:
	// http://localhost:3000/api/users?page=2
	return c.JSON(fiber.Map{
		"data": products,
		"meta": fiber.Map{
			"page":      page,
			"total":     total,
			"last_page": math.Floor(float64(int(total)/limit)) + 1, //
		}})
}

// CreateProduct - создание продукта в БД.
// Например: POST http://localhost:3000/api/products
func CreateProduct(c *fiber.Ctx) error {
	var product models.Product

	// парсим данные. Если данные не подходят - выходим с ошибкой
	if err := c.BodyParser(&product); err != nil {
		return err
	}

	database.DB.Create(&product)

	return c.JSON(product)
}

// GetProduct - получение данных продукта по его ИД (из параметров URL),
// например: GET http://localhost:3000/api/products/2
func GetProduct(c *fiber.Ctx) error {
	// берем параметр из URL
	id, _ := strconv.Atoi(c.Params("id"))

	//создаем новую структуру с заполненным Id, по которому ниже найдем в базе
	product := models.Product{
		Id: uint(id),
	}

	// ищем в базе
	database.DB.Find(&product)

	//выводим данные о пользователе в виде JSON
	return c.JSON(product)
}

// UpdateProduct - обновление имеющихся данных продукта по его ИД (из параметров URL),
// например: PUT http://localhost:3000/api/products/2
func UpdateProduct(c *fiber.Ctx) error {
	// берем параметр из URL
	id, _ := strconv.Atoi(c.Params("id"))

	//создаем новую структуру с заполненным Id, по которому ниже найдем в базе
	product := models.Product{
		Id: uint(id),
	}

	// парсим данные, которые ввел пользователь. Если данные не подходят - выходим с ошибкой
	if err := c.BodyParser(&product); err != nil {
		return err
	}

	// обновляем данные в базе
	database.DB.Model(&product).Updates(product)

	// возвращаем json  с измененными данными
	return c.JSON(product)
}

// DeleteProduct - удаление продукта по его ИД (из параметров URL),
// например: DELETE http://localhost:3000/api/products/2
func DeleteProduct(c *fiber.Ctx) error {
	// берем параметр из URL
	id, _ := strconv.Atoi(c.Params("id"))

	//создаем новую структуру с заполненным Id, по которому ниже найдем в базе
	product := models.Product{
		Id: uint(id),
	}

	// удаляем пользователя из базы
	database.DB.Delete(&product)

	// возвращаем json  с измененными данными
	return nil
}
