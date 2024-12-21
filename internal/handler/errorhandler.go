package handler

import (
	"github.com/gofiber/fiber/v2"
)

func CustomErrorHandler(c *fiber.Ctx, statusCode int, msg string) error {
	return c.Status(statusCode).JSON(fiber.Map{
		"status":  false,
		"data":    map[string]interface{}{},
		"message": msg,
	})
}
