package main

import (
	"admin/src/database"
	"admin/src/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	// Connection To Mysql
	database.Connect()
	// Migration
	database.AutoMigrate()
	// Redis
	database.SetupRedis()
	// Cache
	database.SetupCacheChannel()

	// fiber API
	app := fiber.New()

	// CORSの設定
	app.Use(cors.New(cors.Config{
		// 認証にcookieなどの情報を必要とするかどうか
		AllowCredentials: true,
	}))

	// Setup Routes
	routes.Setup(app)

	app.Listen(":3000")
}
