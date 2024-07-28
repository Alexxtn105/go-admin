package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"go-admin/utils"
)

// IsAuthenticated - показывает, аутентифицирован ли пользователь
func IsAuthenticated(c *fiber.Ctx) error {
	log.Info("Middleware - getting cookies")
	// берем текущие куки
	cookie := c.Cookies("jwt")
	log.Info("Current cookies len: ", len(cookie))

	//парсим из них токен
	if _, err := utils.ParseJwt(cookie); err != nil {
		log.Error("Unauthorized (401). Error: ", err)
		c.Status(fiber.StatusUnauthorized) //401
		return c.JSON(fiber.Map{
			"message": "unauthenticated",
		})
	}

	// поскольку текущая функция это middleware,
	// возвращаемое значение это переход к следующему маршруту - c.Next()
	return c.Next()
}
