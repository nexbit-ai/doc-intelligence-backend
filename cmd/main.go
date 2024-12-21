package main

import (
	"fmt"
	"os"
	"time"

	router "nexbit/internal/router/v1"
	service "nexbit/internal/service"

	external "nexbit/external"

	"github.com/joho/godotenv"

	externalDIClient "nexbit/external/microsoft"

	"github.com/gofiber/fiber/v2"

	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	fmt.Println("Starting the server...")
	app := fiber.New()

	// Add the recover middleware
	app.Use(recover.New(recover.Config{
		EnableStackTrace:  true,
		StackTraceHandler: stackTraceHandler,
	}))

	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error in loading env variables", err)
	}
	// user := os.Getenv("DB_USER")
	// dbname := os.Getenv("DB_NAME")
	// password := os.Getenv("DB_PASSWORD")
	// host := os.Getenv("DB_HOST")
	// port := os.Getenv("DB_PORT")
	// sslmode := os.Getenv("DB_SSLMODE")

	port := os.Getenv("PORT")
	// // connStr := fmt.Sprintf("user=%s dbname=%s password=%s host=%s port=%s sslmode=%s",
	// // 	user, dbname, password, host, port, sslmode)

	// // dbService, err := repo.NewDBService(connStr)
	// // if err != nil {
	// // 	log.Fatalln(err)
	// // }

	// // err = dbService.Ping()
	// // if err != nil {
	// // 	log.Fatalf("Error connecting to the database: %v\n", err)
	// // } else {
	// // 	fmt.Println("Successfully connected to the PostgreSQL database!")
	// // }

	// // defer dbService.Close()

	httpClient := external.NewHTTPClient(50 * time.Second)

	externalDIClient := externalDIClient.NewDIClientClient(httpClient)

	docService := service.NewDocService(externalDIClient)
	router.DocRouter(app, docService)

	if err := app.Listen(":" + port); err != nil {
		fmt.Println("Error starting server:", err)
	}
}

func stackTraceHandler(ctx *fiber.Ctx, err interface{}) {
	errMsg := fmt.Sprintf("Panic: %v", err)
	ctx.Status(fiber.StatusInternalServerError)
	err = ctx.JSON(fiber.Map{
		"error":   errMsg,
		"message": "Internal Server Error. Please try again later.",
	})
}
