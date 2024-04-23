package models

import "gorm.io/gorm"

type Order struct {
	// здесь вместо Id, UpdatedAt и CreatedAt, можно использовать сруктуру gorm.Model,
	// но там время не типа string, а time.Time, так что пропишем все вручную
	//gorm.Model

	Id         uint        `json:"id"`
	FirstName  string      `json:"first_name"`
	LastName   string      `json:"last_name"`
	Email      string      `json:"email"`
	UpdatedAt  string      `json:"updated_at"`
	CreatedAt  string      `json:"created_at"`
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
func (p *Order) Count(db *gorm.DB) int64 {
	var total int64
	db.Model(&Order{}).Count(&total)
	return total
}

// Take - возвращает слайс для постраничного вывода
func (p *Order) Take(db *gorm.DB, limit int, offset int) any {
	var products []Order
	db.Preload("OrderItems").Offset(offset).Limit(limit).Find(&products)
	return products
}
