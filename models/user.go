package models

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	Id        uint   `json:"id"`                             // переопределяем имена полей в ответах сервера
	FirstName string `json:"first_name"`                     // переопределяем имена полей в ответах сервера
	LastName  string `json:"last_name"`                      // переопределяем имена полей в ответах сервера
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

// ComparePassword -  chfdybnm gfhjkm gjkmpjdfntkz c
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
	return db.Offset(offset).Limit(limit).Find(&users)
	return users
}
