package routes

import (
	"github.com/gofiber/fiber/v2"

	"product-crud/backend/handlers"
)

func SetupProductRoutes(app fiber.Router, handler *handlers.ProductHandler) {
	products := app.Group("/products")
	products.Get("/", handler.GetAll)
	products.Get("/:id", handler.GetByID)
	products.Post("/", handler.Create)
	products.Put("/:id", handler.Update)
	products.Delete("/:id", handler.Delete)
}
