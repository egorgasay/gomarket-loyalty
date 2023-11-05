package controller

import (
	"github.com/gofiber/fiber/v2"
	"gomarket-loyalty/service"
)

type ProductController struct {
	ProductService service.ProductService
}

func NewProductController(productService *service.ProductService) ProductController {
	return ProductController{ProductService: *productService}
}

func (controller *ProductController) Route(app *fiber.App) {
	app.Get("/", controller.Base)
}

func (controller *ProductController) Base(c *fiber.Ctx) error {
	response := controller.ProductService.Base()

	return c.SendString(response)
}
