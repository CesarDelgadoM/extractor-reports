package branch

import (
	"github.com/CesarDelgadoM/extractor-reports/internal/requests"
	"github.com/CesarDelgadoM/extractor-reports/pkg/httperrors"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	app     *fiber.App
	service IService
}

func NewBranchHandler(app *fiber.App, service IService) {
	handler := Handler{
		app:     app,
		service: service,
	}

	handler.initRouters()
}

func (h *Handler) initRouters() {
	api := h.app.Group("/extract")

	api.Post("/branches", h.producerBranches)
}

func (h *Handler) producerBranches(ctx *fiber.Ctx) error {
	var request requests.RestaurantRequest

	err := ctx.BodyParser(&request)
	if err != nil {
		return httperrors.BodyParserFailed
	}

	err = h.service.ProducerReport(request)
	if err != nil {
		return err
	}

	return ctx.JSON("Processing restaurant report...")
}
