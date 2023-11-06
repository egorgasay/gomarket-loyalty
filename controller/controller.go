package controller

import (
	"github.com/gofiber/fiber/v2"
	"gomarket-loyalty/exception"
	"gomarket-loyalty/model"
	"gomarket-loyalty/service"
)

type Controller struct {
	service service.Service
}

func NewController(Service *service.Service) Controller {
	return Controller{service: *Service}
}

func (controller *Controller) Route(app *fiber.App) {
	app.Get("/", controller.Base)
	app.Post("/v1/user", controller.Create)
}

func (controller *Controller) Base(c *fiber.Ctx) error {
	response := controller.service.Base()
	return c.SendString(response)
}

func (controller *Controller) Create(c *fiber.Ctx) error {
	var user model.RegisterRequest
	if err := c.BodyParser(&user); err != nil {
		return exception.ErrorHandler(c, err)
	}
	err := controller.service.Create(user)

	return exception.ErrorHandler(c, err)

}
