package models

type User struct {
	Id        uint
	FirstName string
	LastName  string
	Email     string `gorm:"unique"` // уникальный (для миграции)
	Password  []byte
}
