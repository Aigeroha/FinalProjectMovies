package middleware

import "github.com/gofiber/fiber/v3"

func AdminOnly() fiber.Handler {
	return func(c fiber.Ctx) error {

		secret := c.Get("X-Admin-Secret")

		if secret != "my-super-secret-key" {
			return c.Status(403).JSON(fiber.Map{
				"success": false,
				"error":   "доступ запрещен: требуется секретный ключ администратора",
			})
		}
		return c.Next()
	}
}
