package routes

import (
	"admin/src/controllers"

	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	// POST
	app.Post("/api/admin/register", controllers.Register)
}
