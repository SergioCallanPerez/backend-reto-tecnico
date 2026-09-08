package handler

import (
	"github.com/gofiber/fiber/v2"

	"github.com/SergioCallanPerez/backend-reto-tecnico/api-go/internal/model"
)

type matrixRotator interface {
	Rotate(matrix [][]float64) ([][]float64, error)
}

// RotateHandler como segundo endpoint para rotacion
type RotateHandler struct {
	rotateService matrixRotator
}

func NewRotateHandler(rotateService matrixRotator) *RotateHandler {
	return &RotateHandler{rotateService: rotateService}
}

func (h *RotateHandler) RegisterRoutes(router fiber.Router) {
	router.Post("/matrix/rotate", h.PostMatrixRotate)
}

func (h *RotateHandler) PostMatrixRotate(c *fiber.Ctx) error {
	var req model.MatrixRequest
	if err := c.BodyParser(&req); err != nil {
		return writeError(c, model.NewValidationError("invalid request body", err.Error()))
	}
	if appErr := validateMatrixShape(req.Matrix); appErr != nil {
		return writeError(c, appErr)
	}
	rotated, err := h.rotateService.Rotate(req.Matrix)
	if err != nil {
		return writeError(c, err)
	}
	return c.JSON(model.MatrixRotateResponse{Matrix: roundMatrix(rotated)})
}
