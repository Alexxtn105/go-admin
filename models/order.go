package models

import "gorm.io/gorm"

type Order struct {
	// здесь вместо Id, UpdatedAt и CreatedAt, можно использовать сруктуру gorm.Model,
	// но там время не типа string, а time.Time, так что пропишем все вручную
	//gorm.Model
	Id        uint   `json:"id"`
	FirstName string `json:"-"`
	LastName  string `json:"-"`
	//вычисляемое поле, которое не нужно добавлять в БД - gorm:"-".
	// для его вычисления - метод Take
	Name  string `json:"name" gorm:"-"`
	Email string `json:"email"`

	// еще одно вычисляемое поле -  общая стоимость заказа
	Total     float64 `json:"total" gorm:"-"`
	UpdatedAt string  `json:"updated_at"`
	CreatedAt string  `json:"created_at"`

	OrderItems []OrderItem `json:"order_items" gorm:"foreignKey:OrderId""` // задаем связь со столбцом OrderId таблицы OrderItem

}

type OrderItem struct {
	Id           uint    `json:"id"`
	OrderId      string  `json:"order_id"`
	ProductTitle string  `json:"product_title"`
	Price        float64 `json:"price"`
	Quantity     uint    `json:"quantity"`
}

// Count - возвращает количество страниц для постраничного вывода
func (o *Order) Count(db *gorm.DB) int64 {
	var total int64
	db.Model(&Order{}).Count(&total)
	return total
}

// Take - возвращает слайс для постраничного вывода
func (o *Order) Take(db *gorm.DB, limit int, offset int) any {
	var orders []Order
	db.Preload("OrderItems").Offset(offset).Limit(limit).Find(&orders)

	// бежим по всем заказам
	for i, _ := range orders {
		total := 0.0
		//бежим по всем заказанным позициям, считаем общее количество
		for _, orderItem := range orders[i].OrderItems {
			total += orderItem.Price * float64(orderItem.Quantity)
		}

		orders[i].Name = orders[i].FirstName + " " + orders[i].LastName
		orders[i].Total = total

	}
	return orders
}
