package models

import "golang.org/x/crypto/bcrypt"

type User struct {
	Id        uint   `json:"id"`                  // переопределяем имена полей в ответах сервера
	FirstName string `json:"first_name"`          // переопределяем имена полей в ответах сервера
	LastName  string `json:"last_name"`           // переопределяем имена полей в ответах сервера
	Email     string `json:"email" gorm:"unique"` // имя, а также уникальный столбец  gorm (для миграции)
	Password  []byte `json:"-"`                   // не показывать пароль в ответе сервера
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
