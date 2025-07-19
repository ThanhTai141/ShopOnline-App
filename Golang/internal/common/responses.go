package common

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type PaginatedResponse struct {
	Data   any   `json:"data"`
	Total  int64 `json:"total"`
	Limit  int   `json:"limit"`
	Offset int   `json:"offset"`
}

func RespondError(c *fiber.Ctx, code int, msg string) error {
	return c.Status(code).JSON(fiber.Map{
		"error":   http.StatusText(code),
		"message": msg,
	})
}

func RespondCreated(c *fiber.Ctx, id uint) error {
	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"id": id,
	})
}

func RespondseSuccess(c *fiber.Ctx, data any) error {
	if data == nil {
		return c.SendStatus(http.StatusNoContent)
	}
	return c.Status(http.StatusOK).JSON(data)
}

func RespondNoContent(c *fiber.Ctx) error {
	return c.SendStatus(http.StatusNoContent)
}
