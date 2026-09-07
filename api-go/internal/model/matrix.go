package model

type MatrixRequest struct {
	Matrix [][]float64 `json:"matrix"`
}

type QRResult struct {
	Q [][]float64 `json:"q"`
	R [][]float64 `json:"r"`
}

type MatrixQRResponse struct {
	QR    QRResult    `json:"qr"`
	Stats MatrixStats `json:"stats"`
}

type DiagonalCheck struct {
	Q           bool `json:"q"`
	R           bool `json:"r"`
	AnyDiagonal bool `json:"anyDiagonal"`
}

type MatrixStats struct {
	Max           float64       `json:"max"`
	Min           float64       `json:"min"`
	Average       float64       `json:"average"`
	Sum           float64       `json:"sum"`
	DiagonalCheck DiagonalCheck `json:"diagonalCheck"`
}
