package exception

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
)

func ErrorHandler(c *fiber.Ctx, err error) error {
	var unmarshalTypeError *json.UnmarshalTypeError

	if err != nil {
		if errors.Is(err, ErrAlreadyExists) {
			err := fmt.Sprintf("login is already exists %s", err)
			return c.Status(409).JSON(fiber.Map{
				"message": err,
			})
		}
		if errors.Is(err, ErrEnabledData) || errors.As(err, &unmarshalTypeError) {
			err := fmt.Sprintf("enabled data %s", err)
			return c.Status(400).JSON(fiber.Map{
				"message": err,
			})
		}

		if errors.Is(err, ErrNotFound) {
			err := fmt.Sprintf("error not data %s", err)
			return c.Status(204).JSON(fiber.Map{
				"message": err,
			})
		}
		err := fmt.Sprintf("error create user %s", err)
		return c.Status(500).JSON(fiber.Map{
			"message": err,
		})
	}
	return c.Status(200).SendString("success")
}
