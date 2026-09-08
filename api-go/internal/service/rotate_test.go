package service

import (
	"reflect"
	"testing"

	"github.com/SergioCallanPerez/backend-reto-tecnico/api-go/internal/model"
)

func TestRotate_HappyPath(t *testing.T) {
	svc := NewRotateService()

	cases := []struct {
		name     string
		matrix   [][]float64
		expected [][]float64
	}{
		{
			name:     "cuadrada 2x2",
			matrix:   [][]float64{{1, 2}, {3, 4}},
			expected: [][]float64{{3, 1}, {4, 2}},
		},
		{
			name:     "rectangular 2x3",
			matrix:   [][]float64{{1, 2, 3}, {4, 5, 6}},
			expected: [][]float64{{4, 1}, {5, 2}, {6, 3}},
		},
		{
			name:     "1x1 trivial",
			matrix:   [][]float64{{7}},
			expected: [][]float64{{7}},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			rotated, err := svc.Rotate(tc.matrix)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !reflect.DeepEqual(rotated, tc.expected) {
				t.Fatalf("got %v, want %v", rotated, tc.expected)
			}
		})
	}
}

func TestRotate_DimensionsSwap(t *testing.T) {
	svc := NewRotateService()
	// m x n (2x3) rotada debe quedar n x m (3x2).
	rotated, err := svc.Rotate([][]float64{{1, 2, 3}, {4, 5, 6}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rotated) != 3 || len(rotated[0]) != 2 {
		t.Fatalf("dimensiones incorrectas: got %dx%d, want 3x2", len(rotated), len(rotated[0]))
	}
}

func TestRotate_ValidationErrors(t *testing.T) {
	svc := NewRotateService()

	cases := []struct {
		name   string
		matrix [][]float64
	}{
		{name: "matriz vacia", matrix: [][]float64{}},
		{name: "filas vacias", matrix: [][]float64{{}}},
		{name: "no rectangular", matrix: [][]float64{{1, 2}, {3}}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := svc.Rotate(tc.matrix)
			if err == nil {
				t.Fatal("se esperaba un error, se obtuvo nil")
			}
			appErr, ok := err.(*model.AppError)
			if !ok || appErr.Code != model.ErrCodeValidation {
				t.Fatalf("se esperaba VALIDATION_ERROR, se obtuvo %v", err)
			}
		})
	}
}
