package route

import (
	"github.com/gofiber/fiber/v2"
	"net/http"
)

func Wiki(c *fiber.Ctx) error {
	return c.Redirect("https://thftgr.stoplight.io/docs/beatmap-mirror", http.StatusPermanentRedirect)
}
