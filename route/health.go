package route

import (
	"github.com/gofiber/fiber/v2"
	"net/http"
)

func Health(c *fiber.Ctx) (err error) {
	c.Status(http.StatusOK)
	return
}
