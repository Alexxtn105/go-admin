package models

import "gorm.io/gorm"

type Product struct {
	Id          uint    `json:"id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Image       string  `json:"image"`
	Price       float64 `json:"price"`
}

// Count - возвращает количество страниц для постраничного вывода
func (p *Product) Count(db *gorm.DB) int64 {
	var total int64
	db.Model(&Product{}).Count(&total)
	return total
}

// Take - возвращает слайс для постраничного вывода
func (p *Product) Take(db *gorm.DB, limit int, offset int) any {
	var products []Product
	db.Offset(offset).Limit(limit).Find(&products)
	return products
}
