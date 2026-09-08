package service

import "github.com/SergioCallanPerez/backend-reto-tecnico/api-go/internal/model"

type RotateService struct{}

func NewRotateService() *RotateService {
	return &RotateService{}
}

func (s *RotateService) Rotate(a [][]float64) ([][]float64, error) {
	m := len(a)
	if m == 0 {
		return nil, model.NewValidationError("la matriz de entrada está vacía", "")
	}
	n := len(a[0])
	if n == 0 {
		return nil, model.NewValidationError("la matriz de entrada tiene filas vacías", "")
	}
	for _, row := range a {
		if len(row) != n {
			return nil, model.NewValidationError("la matriz de entrada no es rectangular", "todas las filas deben tener el mismo largo")
		}
	}

	rotated := make([][]float64, n)
	for i := 0; i < n; i++ {
		rotated[i] = make([]float64, m)
		for j := 0; j < m; j++ {
			rotated[i][j] = a[m-1-j][i]
		}
	}
	return rotated, nil
}
