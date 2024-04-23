package controllers

import (
	"github.com/gofiber/fiber/v2"
	"go-admin/database"
	"go-admin/models"
	"strconv"
)

func AllOrders(c *fiber.Ctx) error {
	// ПРИКРУТИМ страничный режим
	// берем номер страницы из параметра URL "page", по умолчанию - "1"
	page, _ := strconv.Atoi(c.Query("page", "1"))

	return c.JSON(models.Paginate(database.DB, &models.Order{}, page))
}
