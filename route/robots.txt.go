package route

import (
	"github.com/gofiber/fiber/v2"
	"net/http"
)

func Robots(c *fiber.Ctx) (err error) {
	c.Status(http.StatusInternalServerError)
	return
}
