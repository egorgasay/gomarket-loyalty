package controller

import (
	"context"
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
	app.Post("/v1/mechanics", controller.RegisterMechanic)
	app.Post("/v1/orders", controller.CreateOrder)
	app.Get("/v1/orders", controller.GetInfoOrders)
}

func (controller *Controller) Base(c *fiber.Ctx) error {
	response := controller.service.Base()
	return c.SendString(response)
}

func (controller *Controller) Create(c *fiber.Ctx) error {
	var user model.RegisterRequest

	ctx, cancel := context.WithCancel(c.Context())
	defer cancel()

	if err := c.BodyParser(&user); err != nil {
		return exception.ErrorHandler(c, err)
	}
	err := controller.service.Create(ctx, user)

	return exception.ErrorHandler(c, err)

}

func (controller *Controller) RegisterMechanic(c *fiber.Ctx) error {
	var bonus model.Mechanic

	ctx, cancel := context.WithCancel(c.Context())
	defer cancel()

	if err := c.BodyParser(&bonus); err != nil {
		return exception.ErrorHandler(c, err)
	}

	err := controller.service.AddMechanic(ctx, bonus)

	return exception.ErrorHandler(c, err)

}

func (controller *Controller) CreateOrder(c *fiber.Ctx) error {
	var order model.Items

	ctx, cancel := context.WithCancel(c.Context())
	defer cancel()

	orderID := c.Query("order_id")
	clientID := c.Query("client_id")
	if clientID == "" || orderID == "" {
		return exception.ErrorHandler(c, exception.ErrEnabledData)
	}

	if err := c.BodyParser(&order); err != nil {
		return exception.ErrorHandler(c, err)
	}

	err := controller.service.CreateOrder(ctx, clientID, orderID, order)
	return exception.ErrorHandler(c, err)

}

func (controller *Controller) GetInfoOrders(c *fiber.Ctx) error {
	ctx, cancel := context.WithCancel(c.Context())
	defer cancel()

	clientID := c.Query("client_id")
	if clientID == "" {
		return exception.ErrorHandler(c, exception.ErrEnabledData)
	}
	orders, err := controller.service.GetInfoOrders(ctx, clientID)
	if err != nil {
		return exception.ErrorHandler(c, err)
	}

	return c.JSON(orders)
}
