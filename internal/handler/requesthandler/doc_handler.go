package requesthandler

import (
	"fmt"
	docService "nexbit/internal/service"
	"nexbit/models"

	"github.com/gofiber/fiber/v2"
)

type DocHandler struct {
	docService *docService.DocService
}

func NewDoctHandler(docService *docService.DocService) *DocHandler {
	return &DocHandler{
		docService: docService,
	}
}

func ParseDocHandler(docService docService.DocService) fiber.Handler {
	return func(ctx *fiber.Ctx) error {

		modelID := ctx.Query("modelID")
		if modelID == "" {
			return fmt.Errorf("empty modelID")
		}

		var reqData models.FetchDocumentRequest
		err := ctx.BodyParser(&reqData)
		if err != nil {
			return err
		}

		if reqData.Base64Source == "" {
			return fmt.Errorf("empty base64 string")
		}

		resp, err := docService.ParseDoc(ctx, modelID, reqData)
		if err != nil {
			return err
		}
		return ctx.JSON(
			fiber.Map{
				"response": resp,
			})
	}
}

func HealthCheckHandler(docService docService.DocService) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		fmt.Println("Health check")
		return ctx.JSON(
			fiber.Map{
				"response": "all good",
			})
	}
}
