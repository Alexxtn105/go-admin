package models

import "gorm.io/gorm"

// Entity - интерфейс для постраничного вывода
type Entity interface {
	Count(db *gorm.DB) int64
	Take(db *gorm.DB, limit int, offset int) interface{}
}
