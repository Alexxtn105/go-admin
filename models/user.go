package models

type User struct {
	Id        uint   `json:"id"`                  // переопределяем имена полей в ответах сервера
	FirstName string `json:"first_name"`          // переопределяем имена полей в ответах сервера
	LastName  string `json:"last_name"`           // переопределяем имена полей в ответах сервера
	Email     string `json:"email" gorm:"unique"` // имя, а также уникальный столбец  gorm (для миграции)
	Password  []byte `json:"-"`                   // не показывать пароль в ответе сервера
}
