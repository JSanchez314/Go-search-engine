package routes

import (
	//"pro/jet/db"

	"github.com/gofiber/fiber/v2"
)

func SetRoutes(app *fiber.App) {
	app.Get("/", DashboardHandler)
	app.Post("/", DashboarPostdHandler)

	app.Get("/login", LoginHandler)
	app.Post("/login", LoginPostHandler)
	app.Post("/logout", LogOutHandler)

}
