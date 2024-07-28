package middleware

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"go-admin/database"
	"go-admin/models"
	"go-admin/utils"
	"strconv"
)

func IsAuthorized(c *fiber.Ctx, page string) error {
	log.Info("getting cookies")
	// берем текущие куки
	cookie := c.Cookies("jwt")
	log.Info("Parsing JWT")
	//парсим из них токен
	Id, err := utils.ParseJwt(cookie)

	if err != nil {
		return err
	}

	// сперва найдем пользователя по его Id в куках
	userId, _ := strconv.Atoi(Id)
	user := models.User{
		Id: uint(userId),
	}
	database.DB.Preload("Role").Find(&user)

	// теперь найдем его роль
	role := models.Role{
		Id: user.RoleId,
	}
	database.DB.Preload("Permissions").Find(&role)

	// найдем разрешения для роли, имеется ли доступ к странице (имена ролей: view_users, edit_users и тп.)
	if c.Method() == "GET" {
		for _, permission := range role.Permissions {
			if permission.Name == "view_"+page || permission.Name == "view_"+page {
				return nil
			}
		}
	} else {
		for _, permission := range role.Permissions {
			if permission.Name == "edit_"+page || permission.Name == "edit_"+page {
				return nil
			}
		}
	}

	c.Status(fiber.StatusUnauthorized)
	return errors.New("unauthorized")
}
