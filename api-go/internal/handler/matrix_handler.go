package handler

import (
	"log/slog"
	"math"

	"github.com/gofiber/fiber/v2"

	"github.com/SergioCallanPerez/backend-reto-tecnico/api-go/internal/model"
)

// Maximo para evitar payloads demasiado grandes.
const maxMatrixDimension = 500

// Decimales visibles
const outputPrecision = 6

type qrFactorizer interface {
	Factorize(matrix [][]float64) (q [][]float64, r [][]float64, err error)
}

type statsFetcher interface {
	FetchStats(q, r [][]float64) (model.MatrixStats, error)
}

type MatrixHandler struct {
	qrService   qrFactorizer
	statsClient statsFetcher
}

func NewMatrixHandler(qrService qrFactorizer, statsClient statsFetcher) *MatrixHandler {
	return &MatrixHandler{qrService: qrService, statsClient: statsClient}
}

func (h *MatrixHandler) RegisterRoutes(router fiber.Router) {
	router.Post("/matrix/qr", h.PostMatrixQR)
}

// Para POST /api/v1/matrix/qr.
func (h *MatrixHandler) PostMatrixQR(c *fiber.Ctx) error {
	var req model.MatrixRequest
	if err := c.BodyParser(&req); err != nil {
		return writeError(c, model.NewValidationError("el cuerpo de la solicitud es inválido", err.Error()))
	}
	if appErr := validateMatrixShape(req.Matrix); appErr != nil {
		return writeError(c, appErr)
	}
	q, r, err := h.qrService.Factorize(req.Matrix)
	if err != nil {
		return writeError(c, err)
	}
	stats, err := h.statsClient.FetchStats(q, r)
	if err != nil {
		return writeError(c, err)
	}
	return c.JSON(model.MatrixQRResponse{
		QR:    model.QRResult{Q: roundMatrix(q), R: roundMatrix(r)},
		Stats: stats,
	})
}

// Redondeo de la salida para visualización
func roundMatrix(matrix [][]float64) [][]float64 {
	factor := math.Pow(10, outputPrecision)
	rounded := make([][]float64, len(matrix))
	for i, row := range matrix {
		rounded[i] = make([]float64, len(row))
		for j, v := range row {
			rounded[i][j] = math.Round(v*factor) / factor
		}
	}
	return rounded
}

func validateMatrixShape(matrix [][]float64) *model.AppError {
	if len(matrix) == 0 {
		return model.NewValidationError("la matriz de entrada está vacía", "")
	}
	if len(matrix) > maxMatrixDimension || len(matrix[0]) > maxMatrixDimension {
		return model.NewValidationError(
			"la matriz de entrada excede el tamaño máximo permitido",
			"las matrices están limitadas a 500x500",
		)
	}
	return nil
}

func writeError(c *fiber.Ctx, err error) error {
	appErr, ok := err.(*model.AppError)
	if !ok {
		appErr = model.NewInternalError(err.Error())
	}
	logError(appErr)
	return c.Status(appErr.HTTPStatus()).JSON(appErr)
}

func logError(appErr *model.AppError) {
	attrs := []any{"code", appErr.Code, "message", appErr.Message, "detail", appErr.Detail}
	if appErr.Code == model.ErrCodeValidation {
		slog.Warn("request_rejected", attrs...)
	} else {
		slog.Error("request_failed", attrs...)
	}
}
