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
	//бежим по всем фалам, помещаем их в нужную папку
	//ПАПКА ДОЛЖНА СУЩЕСТВОВАТЬ!!!
	for _, file := range files {
		// сохранячем файл
		filename = file.Filename

		if err := c.SaveFile(file, "./uploads/"+filename); err != nil {
			return err
		}

		//
	}

	//возвращаем URL загруженного файла

	return c.JSON(fiber.Map{
		"url": "http://localhost:8000/uploads/" + filename,
	})
}
