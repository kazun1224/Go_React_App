package controllers

import "github.com/gofiber/fiber/v2"

func Register(ctx *fiber.Ctx) error {
	return ctx.JSON(fiber.Map{
		"message": "hello",
	})
}
