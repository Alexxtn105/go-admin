package models

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"math"
)

// Paginate - постраничный вывод данных
func Paginate(db *gorm.DB, entity Entity, page int) fiber.Map {
	//вводим ограничения для постраничного вывода пользователей, если их много
	limit := 15
	//начальная позиция на выбранной странице
	offset := (page - 1) * limit

	// берем данные для постраничного вывода (слайс)
	data := entity.Take(db, limit, offset)

	//общее количество
	total := entity.Count(db)

	return fiber.Map{
		"data": data,
		"meta": fiber.Map{
			"page":      page,
			"total":     total,
			"last_page": math.Floor(float64(int(total)/limit)) + 1,
		}}
}
