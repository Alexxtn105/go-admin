package middleware

import (
	"github.com/gofiber/fiber/v2"
	"go-admin/utils"
)

// IsAuthenticated - показывает, аутентифицирован ли пользователь
func IsAuthenticated(c *fiber.Ctx) error {
	// берем текущие куки
	cookie := c.Cookies("jwt")

	//парсим из них токен
	if _, err := utils.ParseJwt(cookie); err != nil {
		c.Status(fiber.StatusUnauthorized) //401
		return c.JSON(fiber.Map{
			"message": "unauthenticated",
		})
	}

	// поскольку текущая функция это middleware,
	// возвращаемое значение это переход к следующему маршруту - c.Next()
	return c.Next()
}
