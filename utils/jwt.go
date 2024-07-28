package utils

import (
	"github.com/dgrijalva/jwt-go/v4"
	"github.com/gofiber/fiber/v2/log"
	"time"
)

// TODO: сделать надежный secretKey
const secretKey = "secret"

// GenerateJwt - генерация токена по ID пользователя
func GenerateJwt(issuer string) (string, error) {
	// создаем JWT-токен
	// сперва создадим заявку (claims)
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    issuer,                                                      // определяет приложение, из которого отправляется токен. В нашем случа будем хранить здесь ID пользователя
		ExpiresAt: jwt.NewTime(float64(time.Now().Add(time.Hour * 24).Unix())), // когда истекает (1 день)
		//Subject:   "Auth",                                                    // определяет тему токена. Пока не используется, сокращаем размер токена
	})

	// создаем токен на основе заявки
	// TODO: ДОБАВИТЬ СЕКРЕТНЫЙ КЛЮЧ ВМЕСТО secret!!!
	return claims.SignedString([]byte(secretKey))
}

// ParseJwt - возвращает ID пользователя по его кукам
func ParseJwt(cookie string) (string, error) {
	//получаем токен из кук (проделываем операцию, обратную созданию токена)
	log.Info("Пытаемся получить токен из кук: ", cookie)
	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil || !token.Valid {
		log.Error(err)
		log.Error("Необходимо проверить CORS!!!")
		return "", err
	}

	// получаем Claims (заявку), с преобразованием заявки в тип *jwt.StandardClaims
	claims := token.Claims.(*jwt.StandardClaims)

	// Возвращаем приложение (Issuer) из claims (заявки). Фактически это ID пользователя
	return claims.Issuer, nil
}
