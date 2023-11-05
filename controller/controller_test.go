package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	config2 "gomarket-loyalty/config"
	"gomarket-loyalty/repository"
	"gomarket-loyalty/service"
)

func createTestApp() *fiber.App {
	var app = fiber.New(config2.NewFiberConfig())
	app.Use(recover.New())
	productController.Route(app)
	return app
}

var configuration = config2.New("../.env.test")

var database = config2.NewMongoDatabase(configuration)
var productRepository = repository.NewProductRepository(database)
var productService = service.NewProductService(&productRepository)

var productController = NewProductController(&productService)

var app = createTestApp()
