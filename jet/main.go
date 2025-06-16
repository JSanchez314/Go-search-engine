package main

import (
	"log"
	"os"
	"os/signal"
	"pro/jet/db"
	"pro/jet/routes"
	"pro/jet/utils"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/template/handlebars/v2"
	"github.com/joho/godotenv"
)

func main() {
	engine := handlebars.New("./views", ".hbs")

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Could not find the environment variables")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = ":4000"
	} else {
		port = ":" + port
	}

	log.Println("Starting server on port", port)

	app := fiber.New(fiber.Config{
		IdleTimeout: 5 * time.Second,
		Views:       engine,
	})

	app.Use(compress.New())

	log.Println("Initializing database connection")
	db.InitDB()

	log.Println("Setting routes")
	routes.SetRoutes(app)
	utils.StartCronJobs()

	go func() {
		if err := app.Listen(port); err != nil {
			log.Panic(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	<-c

	log.Println("Shutting down server")
	app.Shutdown()
}
