package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"

	"product-crud/backend/handlers"
	"product-crud/backend/routes"
	"product-crud/backend/store"
)

func main() {
	app := fiber.New(fiber.Config{
		AppName: "Product CRUD API",
	})

	prometheus := fiberprometheus.New("product_crud_api")
	prometheus.RegisterAt(app, "/metrics")
	app.Use(prometheus.Middleware)

	app.Use(logger.New(logger.Config{
		Format: "${time} | ${status} | ${latency} | ${ip} | ${method} | ${path} | ${error}\n",
	}))

	corsOrigins := os.Getenv("CORS_ORIGINS")
	if corsOrigins == "" {
		corsOrigins = "http://localhost:5173"
	}

	app.Use(cors.New(cors.Config{
		AllowOrigins: corsOrigins,
		AllowHeaders: "Origin, Content-Type, Accept",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
	}))

	productStore := store.NewProductStore()
	productHandler := handlers.NewProductHandler(productStore)

	api := app.Group("/api")
	routes.SetupProductRoutes(api, productHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Println(port, "port")
	log.Fatal(app.Listen(":" + port))
}
