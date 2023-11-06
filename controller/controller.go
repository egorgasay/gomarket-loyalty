package controller

import (
	"errors"
	"fmt"
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
	app.Post("/v1/user/register", controller.Register)
}

func (controller *Controller) Base(c *fiber.Ctx) error {
	response := controller.service.Base()
	return c.SendString(response)
}

func (controller *Controller) Register(c *fiber.Ctx) error {
	var user model.RegisterRequest
	if err := c.BodyParser(&user); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": fmt.Errorf("error parsing request body"),
		})
	}
	err := controller.service.Register(user)
	if err != nil {
		if errors.Is(err, exception.ErrLoginAlreadyExists) {
			return c.Status(409).JSON(fiber.Map{
				"message": "login is already exists",
			})
		}
		if errors.Is(err, exception.ErrEnabledData) {
			return c.Status(400).JSON(fiber.Map{
				"message": "enabled data",
			})
		}
		return c.Status(500).JSON(fiber.Map{
			"message": fmt.Errorf("error registering user"),
		})
	}

	return c.Status(200).SendString("success")
}
