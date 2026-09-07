package service

import (
	"math"
	"testing"

	"github.com/SergioCallanPerez/backend-reto-tecnico/api-go/internal/model"
)

const testEpsilon = 1e-9

func TestFactorize_HappyPath(t *testing.T) {
	cases := []struct {
		name   string
		matrix [][]float64
	}{
		{
			name: "square 3x3",
			matrix: [][]float64{
				{1, 2, 3},
				{4, 5, 6},
				{7, 8, 10},
			},
		},
		{
			name: "rectangular 4x2 (mas filas que columnas)",
			matrix: [][]float64{
				{1, 2},
				{3, 4},
				{5, 6},
				{7, 8},
			},
		},
		{
			name:   "1x1 trivial",
			matrix: [][]float64{{5}},
		},
	}

	svc := NewQRService()

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			q, r, err := svc.Factorize(tc.matrix)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			assertReconstructs(t, tc.matrix, q, r)
			assertOrthonormalColumns(t, q)
			assertUpperTriangular(t, r)
		})
	}
}

func TestFactorize_ValidationErrors(t *testing.T) {
	cases := []struct {
		name   string
		matrix [][]float64
	}{
		{name: "matriz vacia", matrix: [][]float64{}},
		{name: "filas vacias", matrix: [][]float64{{}}},
		{name: "no rectangular", matrix: [][]float64{{1, 2, 3}, {4, 5}}},
		{name: "menos filas que columnas", matrix: [][]float64{{1, 2, 3}, {4, 5, 6}}},
		{name: "columnas linealmente dependientes", matrix: [][]float64{{1, 2}, {2, 4}, {3, 6}}},
	}

	svc := NewQRService()

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, _, err := svc.Factorize(tc.matrix)
			if err == nil {
				t.Fatal("se esperaba un error, se obtuvo nil")
			}
			appErr, ok := err.(*model.AppError)
			if !ok {
				t.Fatalf("se esperaba *model.AppError, se obtuvo %T", err)
			}
			if appErr.Code != model.ErrCodeValidation {
				t.Errorf("se esperaba code %s, se obtuvo %s", model.ErrCodeValidation, appErr.Code)
			}
		})
	}
}

func TestFactorize_LinearIndependenceBoundary(t *testing.T) {
	svc := NewQRService()

	t.Run("justo por encima del epsilon: pasa", func(t *testing.T) {
		matrix := [][]float64{
			{1, 0},
			{0, 1e-8},
		}
		if _, _, err := svc.Factorize(matrix); err != nil {
			t.Fatalf("se esperaba exito, se obtuvo error: %v", err)
		}
	})

	t.Run("por debajo del epsilon: falla como dependiente", func(t *testing.T) {
		matrix := [][]float64{
			{1, 1},
			{0, 1e-12},
		}
		if _, _, err := svc.Factorize(matrix); err == nil {
			t.Fatal("se esperaba un error de validacion, se obtuvo nil")
		}
	})
}

func TestVerify(t *testing.T) {
	a := [][]float64{{4, 0}, {0, 9}}
	q := [][]float64{{1, 0}, {0, 1}}

	t.Run("Q y R correctos: true", func(t *testing.T) {
		r := [][]float64{{4, 0}, {0, 9}}
		if !verify(a, q, r) {
			t.Fatal("se esperaba true para una reconstruccion correcta")
		}
	})

	t.Run("R incorrecto: false", func(t *testing.T) {
		r := [][]float64{{4, 0}, {0, 100}}
		if verify(a, q, r) {
			t.Fatal("se esperaba false para una reconstruccion incorrecta")
		}
	})
}

func assertReconstructs(t *testing.T, a, q, r [][]float64) {
	t.Helper()
	m := len(a)
	n := len(a[0])
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			var sum float64
			for k := 0; k < n; k++ {
				sum += q[i][k] * r[k][j]
			}
			if math.Abs(sum-a[i][j]) > testEpsilon {
				t.Fatalf("Q*R no reconstruye A en [%d][%d]: got %v, want %v", i, j, sum, a[i][j])
			}
		}
	}
}

func assertOrthonormalColumns(t *testing.T, q [][]float64) {
	t.Helper()
	m := len(q)
	n := len(q[0])
	for c1 := 0; c1 < n; c1++ {
		for c2 := 0; c2 < n; c2++ {
			var dot float64
			for i := 0; i < m; i++ {
				dot += q[i][c1] * q[i][c2]
			}
			expected := 0.0
			if c1 == c2 {
				expected = 1.0
			}
			if math.Abs(dot-expected) > testEpsilon {
				t.Fatalf("columnas %d y %d de Q no son ortonormales: producto punto = %v, se esperaba %v", c1, c2, dot, expected)
			}
		}
	}
}

func assertUpperTriangular(t *testing.T, r [][]float64) {
	t.Helper()
	n := len(r)
	for i := 0; i < n; i++ {
		for j := 0; j < i; j++ {
			if math.Abs(r[i][j]) > testEpsilon {
				t.Fatalf("R no es triangular superior: R[%d][%d] = %v, se esperaba 0", i, j, r[i][j])
			}
		}
	}
}
