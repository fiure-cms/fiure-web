package middlewares

import "github.com/gofiber/fiber/v2"

func HeaderConf(versionCode string) fiber.Handler {
	return func(c *fiber.Ctx) error {

		// Version
		c.Set("Version", versionCode)

		// Set some security headers:
		c.Set("X-XSS-Protection", "1; mode=block")
		c.Set("X-Content-Type-Options", "nosniff")
		c.Set("X-Download-Options", "noopen")
		c.Set("Strict-Transport-Security", "max-age=5184000")
		c.Set("X-Frame-Options", "SAMEORIGIN")
		c.Set("X-DNS-Prefetch-Control", "off")

		return c.Next()
	}
}
