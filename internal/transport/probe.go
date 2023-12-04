package transport

import "github.com/gofiber/fiber/v2"

func Live(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(nil)
}

func Ready(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(nil)
}
