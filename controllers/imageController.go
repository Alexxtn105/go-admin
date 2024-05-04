package controllers

import (
	"github.com/gofiber/fiber/v2"
)

func Upload(c *fiber.Ctx) error {
	// когда мы отправляем файл, нам нужна MultipartForm
	// для загрузки файла используется form-data request
	// возьмем форму из запроса:
	form, err := c.MultipartForm()
	if err != nil {
		return err
	}

	// файлов в форме может быть несколько
	files := form.File["image"]
	filename := ""
	//бежим по всем файлам, помещаем их в нужную папку
	//ПАПКА ДОЛЖНА СУЩЕСТВОВАТЬ!!!
	for _, file := range files {
		// сохраняем файл
		filename = file.Filename

		if err := c.SaveFile(file, "./uploads/"+filename); err != nil {
			return err
		}

		//
	}

	//возвращаем URL загруженного файла
	return c.JSON(fiber.Map{
		"url": "http://localhost:3000/uploads/" + filename,
	})
}
