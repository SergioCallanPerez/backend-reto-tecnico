package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/SergioCallanPerez/backend-reto-tecnico/api-go/internal/model"
)

const (
	maxAttempts = 5
	retryWait   = 5 * time.Second
)

// Endpoint en Node
type StatsClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewStatsClient(baseURL string, timeout time.Duration) *StatsClient {
	return &StatsClient{baseURL: baseURL, httpClient: &http.Client{Timeout: timeout}}
}

// Envío de Q y R a Node, retorna estadísticas
func (c *StatsClient) FetchStats(q, r [][]float64) (model.MatrixStats, error) {
	body, err := json.Marshal(model.QRResult{Q: q, R: r})
	if err != nil {
		return model.MatrixStats{}, model.NewInternalError("no se pudo codificar la solicitud al servicio de estadísticas")
	}

	url := c.baseURL + "/api/v1/matrix/stats"

	var lastErr error
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		stats, retryable, err := c.doRequest(url, body)
		if err == nil {
			return stats, nil
		}
		lastErr = err
		if !retryable || attempt == maxAttempts {
			break
		}
		slog.Warn("stats_service_retry", "attempt", attempt, "error", err.Error())
		time.Sleep(retryWait)
	}
	return model.MatrixStats{}, lastErr
}

func (c *StatsClient) doRequest(url string, body []byte) (model.MatrixStats, bool, error) {
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return model.MatrixStats{}, false, model.NewInternalError("no se pudo construir la solicitud al servicio de estadísticas")
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return model.MatrixStats{}, true, model.NewUpstreamError("no se pudo contactar al servicio de estadísticas", err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 500 {
		return model.MatrixStats{}, true, model.NewUpstreamError(
			"el servicio de estadísticas respondió con un error de servidor",
			fmt.Sprintf("status %d", resp.StatusCode),
		)
	}
	if resp.StatusCode >= 400 {
		return model.MatrixStats{}, false, model.NewUpstreamError(
			"el servicio de estadísticas rechazó la solicitud",
			fmt.Sprintf("status %d", resp.StatusCode),
		)
	}

	var stats model.MatrixStats
	if err := json.NewDecoder(resp.Body).Decode(&stats); err != nil {
		return model.MatrixStats{}, false, model.NewInternalError("no se pudo decodificar la respuesta del servicio de estadísticas")
	}
	return stats, false, nil
}
