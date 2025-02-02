package v1

import (
	requesthandler "nexbit/internal/handler/requesthandler"
	docService "nexbit/internal/service"

	"github.com/gofiber/fiber/v2"
)

func DocRouter(app *fiber.App, docService docService.DocService) {
	api := app.Group("/v1")
	api.Post("/doc/parse", requesthandler.ParseDocHandler(docService))
	api.Get("/health-check", requesthandler.HealthCheckHandler(docService))
}
