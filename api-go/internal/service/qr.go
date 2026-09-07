package service

import (
	"math"

	"github.com/SergioCallanPerez/backend-reto-tecnico/api-go/internal/model"
)

const verificationEpsilon = 1e-9

type QRService struct{}

func NewQRService() *QRService {
	return &QRService{}
}

func (s *QRService) Factorize(a [][]float64) (q [][]float64, r [][]float64, err error) {
	m := len(a)
	if m == 0 {
		return nil, nil, model.NewValidationError("the input matrix is empty", "")
	}
	n := len(a[0])
	if n == 0 {
		return nil, nil, model.NewValidationError("the input matrix has empty rows", "")
	}
	for _, row := range a {
		if len(row) != n {
			return nil, nil, model.NewValidationError(
				"the input matrix is not rectangular",
				"all rows must have the same length",
			)
		}
	}
	if m < n {
		return nil, nil, model.NewValidationError(
			"the input matrix must have at least as many rows as columns",
			"QR factorization requires m >= n",
		)
	}

	// Se usan columnas porque Gram-Schmidt trabaja vector a vector.
	v := make([][]float64, n)
	for k := 0; k < n; k++ {
		v[k] = make([]float64, m)
		for i := 0; i < m; i++ {
			v[k][i] = a[i][k]
		}
	}

	qCols := make([][]float64, n)
	rMat := make([][]float64, n)
	for i := range rMat {
		rMat[i] = make([]float64, n)
	}

	for k := 0; k < n; k++ {
		norm := euclideanNorm(v[k])
		if norm < verificationEpsilon {
			return nil, nil, model.NewValidationError(
				"the input matrix does not have linearly independent columns",
				"QR factorization requires the columns to be linearly independent",
			)
		}
		rMat[k][k] = norm
		qCols[k] = scale(v[k], 1/norm)

		// Las columnas usan la Q ya calculada, no la A original
		for j := k + 1; j < n; j++ {
			proj := dot(qCols[k], v[j])
			rMat[k][j] = proj
			v[j] = subtractScaled(v[j], qCols[k], proj)
		}
	}

	q = make([][]float64, m)
	for i := range q {
		q[i] = make([]float64, n)
	}
	for k := 0; k < n; k++ {
		for i := 0; i < m; i++ {
			q[i][k] = qCols[k][i]
		}
	}
	r = rMat

	if !verify(a, q, r) {
		return nil, nil, model.NewInternalError("QR factorization failed the Q*R ≈ A sanity check")
	}

	return q, r, nil
}

func verify(a, q, r [][]float64) bool {
	m := len(a)
	n := len(a[0])
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			var sum float64
			for k := 0; k < n; k++ {
				sum += q[i][k] * r[k][j]
			}
			if math.Abs(sum-a[i][j]) > verificationEpsilon {
				return false
			}
		}
	}
	return true
}

func euclideanNorm(v []float64) float64 {
	var sumSquares float64
	for _, x := range v {
		sumSquares += x * x
	}
	return math.Sqrt(sumSquares)
}

func dot(a, b []float64) float64 {
	var sum float64
	for i := range a {
		sum += a[i] * b[i]
	}
	return sum
}

func scale(v []float64, factor float64) []float64 {
	result := make([]float64, len(v))
	for i, x := range v {
		result[i] = x * factor
	}
	return result
}

func subtractScaled(v, w []float64, factor float64) []float64 {
	result := make([]float64, len(v))
	for i := range v {
		result[i] = v[i] - factor*w[i]
	}
	return result
}
