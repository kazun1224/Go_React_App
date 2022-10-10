package routes

import (
	"admin/src/controllers"
	"admin/src/middleware"

	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	// Group
	api := app.Group("api")

	// admin
	admin := api.Group("admin")
	admin.Post("register", controllers.Register)
	admin.Post("login", controllers.Login)

	adminAuthenticated := admin.Use(middleware.IsAuthenticate)
	adminAuthenticated.Get("user", controllers.User)
	adminAuthenticated.Post("logout", controllers.Logout)
	adminAuthenticated.Put("info", controllers.UpdateInfo)
	adminAuthenticated.Put("password", controllers.UpdatePassword)
	adminAuthenticated.Get("ambassadors", controllers.Ambassadors)
	adminAuthenticated.Get("products", controllers.Products)
	adminAuthenticated.Post("products", controllers.CreateProducts)
	adminAuthenticated.Get("products/:id", controllers.GetProduct)
	adminAuthenticated.Put("products/:id", controllers.UpdateProduct)
	adminAuthenticated.Delete("products/:id", controllers.DeleteProduct)
	adminAuthenticated.Get("users/:id/links", controllers.Link)
	adminAuthenticated.Get("orders", controllers.Orders)

	// Ambassador
	ambassador := api.Group("ambassador")
	ambassador.Post("register", controllers.Register)
	ambassador.Post("login", controllers.Login)
	ambassador.Get("products/frontend", controllers.ProductFrontend)
	ambassador.Get("products/backend", controllers.ProductBackend)

	ambassadorAuthentication := ambassador.Use(middleware.IsAuthenticate)
	ambassadorAuthentication.Get("user", controllers.User)
	ambassadorAuthentication.Post("logout", controllers.Logout)
	ambassadorAuthentication.Put("users/info", controllers.UpdateInfo)
	ambassadorAuthentication.Put("users/password", controllers.UpdatePassword)
	ambassadorAuthentication.Post("links", controllers.CreateLink)
	ambassadorAuthentication.Get("stats", controllers.Stats)
	ambassadorAuthentication.Get("rankings", controllers.Ranking)

	// Checkout
	checkout := api.Group("checkout")
	checkout.Get("links/:code", controllers.GetLink)
	checkout.Post("orders", controllers.CreateOrder)
	checkout.Post("orders/confirm", controllers.CompleteOrder)
}
