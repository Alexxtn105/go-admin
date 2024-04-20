# go-admin
Тестовый проект на Go в связке с Vue

Запуск сервера:
go run main.go

----------------------------------------------------
Зависимости:

Fiber (web framework):
go get github.com/gofiber/fiber/v3

GORM - для связи с БД mysql
go get -u gorm.io/gorm
go get -u gorm.io/driver/mysql

Realize - для автоматического перезапуска сервера после внесения изменений
go get github.com/oxequa/realize/realize

bcrypt - для хеширования паролей
go get golang.org/x/crypto/bcrypt

----------------------------------------------------

