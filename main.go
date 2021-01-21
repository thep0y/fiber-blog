/*
* @Author: thepoy
* @Email:  thepoy@aliyun.com
* @File Name: main.go
* @Created:   2021-01-17 12:25:04
* @Modified:  2021-01-21 15:25:22
 */

package main

import (
	"blog_api/apis"
	"blog_api/models"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	app := fiber.New(fiber.Config{
		ServerHeader: "fiber",
		ErrorHandler: func(c *fiber.Ctx, e error) error {
			// 状态码默认为500
			code := fiber.StatusInternalServerError

			if err, ok := e.(*fiber.Error); ok {
				code = err.Code
			}

			e = c.Status(code).JSON(&models.ErrorResponse{
				Error: e.Error(),
			})
			if e != nil {
				return c.Status(code).JSON(&models.ErrorResponse{
					Error: "Internal Server Error",
				})
			}

			return nil
		},
	})

	app.Use(logger.New())

	app.Get("/blog/:id", apis.GetBlog)
	app.Post("/register", apis.Register)
	app.Post("/login", apis.Login)
	app.Post("/send/mail", apis.SendValidateCodeViaMail)

	app.Listen(":3000")
}
