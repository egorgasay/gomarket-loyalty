package config

import (
	"github.com/gofiber/fiber/v2"
	"gomarket-loyalty/exception"
)

func NewFiberConfig() fiber.Config {
	return fiber.Config{
		ErrorHandler: exception.ErrorHandler,
	}
}
