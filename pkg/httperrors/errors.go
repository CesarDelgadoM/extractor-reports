package httperrors

import "github.com/gofiber/fiber/v2"

const (
	ErrRestaurantNotFound = "Restaurant not found: "
)

var (
	BodyParserFailed         = fiber.NewError(fiber.StatusBadRequest, "failed to parser body")
	RequestAlreadyGenerating = fiber.NewError(fiber.StatusTooManyRequests, "request is already generating")

	RestaurantNotFound = fiber.NewError(fiber.StatusNotFound, "restaurant not found")
)
