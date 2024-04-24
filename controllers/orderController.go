package controllers

import (
	"encoding/csv"
	"github.com/gofiber/fiber/v2"
	"go-admin/database"
	"go-admin/models"
	"os"
	"strconv"
)

// AllOrders - вывести список заказов (постранично)
func AllOrders(c *fiber.Ctx) error {
	// ПРИКРУТИМ страничный режим
	// берем номер страницы из параметра URL "page", по умолчанию - "1"
	page, _ := strconv.Atoi(c.Query("page", "1"))

	return c.JSON(models.Paginate(database.DB, &models.Order{}, page))
}

// Export - экспорт данных в csv-файл
func Export(c *fiber.Ctx) error {
	filePath := "./csv/orders.csv"

	if err := CreateFile(filePath); err != nil {
		return err
	}

	return c.Download(filePath)
}

// CreateFile - создать файл
func CreateFile(filePath string) error {
	// создаем файл
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}

	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	var orders []models.Order

	database.DB.Preload("OrderItems").Find(&orders)

	writer.Write([]string{
		"ID", "Name", "Email", "Product Title", "Price", "Quantity",
	})

	for _, order := range orders {
		data := []string{
			strconv.Itoa(int(order.Id)),
			order.FirstName + " " + order.LastName,
			order.Email,
			"",
			"",
			"",
		}

		if err := writer.Write(data); err != nil {
			return err
		}

		for _, orderItem := range order.OrderItems {
			data := []string{
				"",
				"",
				"",
				orderItem.ProductTitle,
				strconv.Itoa(int(orderItem.Price)), // TODO! Почему цена типа int!!!! а как же дробная часть?
				strconv.Itoa(int(orderItem.Quantity)),
			}

			if err := writer.Write(data); err != nil {
				return err
			}
		}
	}

	return nil
}

type Sales struct {
	Date string `json:"date"`
	Sum  string `json:"sum"`
}

// Chart - выполнение произвольного запроса SQL
func Chart(c *fiber.Ctx) error {
	var sales []Sales

	database.DB.Raw(`
	SELECT
		DATE_FORMAT( o.created_at, '%Y-%m-%d') as date,
		SUM(oi.price*oi.quantity) AS sum
	FROM orders o
	JOIN order_items oi ON o.id = oi.order_id
	GROUP BY date
	`).Scan(&sales) //помещаем данные в sales

	return c.JSON(sales)
}
