package common

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

func LoggingFiber(c *fiber.Ctx) error {
	start := time.Now()
	err := c.Next()
	log.Infof("%s %s %d %s",
		c.Method(), c.Path(), c.Response().StatusCode(), time.Since(start))
	return err
}
