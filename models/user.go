package models

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	Id        uint   `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email" gorm:"unique"`            // имя, а также уникальный столбец  gorm (для миграции)
	Password  []byte `json:"-"`                              // не показывать пароль в ответе сервера
	RoleId    uint   `json:"role_id"`                        // роль пользователя
	Role      Role   `json:"role" gorm:"foreignKey:RoleId""` // указываем горму внешний ключ
}

// SetPassword - установка пароля для пользователя
func (u *User) SetPassword(password string) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("1234"), 14)
	u.Password = hashedPassword
}

// ComparePassword -  сравнить пароль пользователя с имеющимся в базе
func (u *User) ComparePassword(password string) error {
	return bcrypt.CompareHashAndPassword(u.Password, []byte(password))
}

// Count - возвращает количество страниц для постраничного вывода
func (u *User) Count(db *gorm.DB) int64 {
	var total int64
	db.Model(&User{}).Count(&total)
	return total
}

// Take - возвращает слайс для постраничного вывода
func (u *User) Take(db *gorm.DB, limit int, offset int) any {
	var users []User
	db.Preload("Role").Offset(offset).Limit(limit).Find(&users)
	return users
}
