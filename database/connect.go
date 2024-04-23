package database

// GORM - для связи с БД mysql
// go get -u gorm.io/gorm
// go get -u gorm.io/driver/mysql

import (
	"go-admin/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// DB - для доступа к БД
var DB *gorm.DB

func Connect() {
	
	// подключаемся к БД (имя_пользователя:пароль@/адрес сервера (для localhost пусто)
	database, err := gorm.Open(mysql.Open("root:Sulubun205!@/go-admin"), &gorm.Config{})
	if err != nil {

		panic("Could not connect to the database")
	}
	DB = database

	// запускаем автоматическую миграцию для создания таблицы пользователей
	database.AutoMigrate(&models.User{}, &models.Role{}, &models.Permission{}, &models.Product{})
}
