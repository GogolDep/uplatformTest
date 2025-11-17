package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type ExchangeClient interface {
	Latest(ctx context.Context) (map[string]float64, error)
	Historical(ctx context.Context, yyyymmdd string) (map[string]float64, error)
}

type exchangeClient struct {
	baseURL string
	appID   string
	client  *http.Client
}

type oxrResponse struct {
	Rates map[string]float64 `json:"rates"`
}

func NewExchangeClient(baseURL, appID string, timeout time.Duration) ExchangeClient {
	return &exchangeClient{
		baseURL: baseURL,
		appID:   appID,
		client:  &http.Client{Timeout: timeout},
	}
}

func (e *exchangeClient) Latest(ctx context.Context) (map[string]float64, error) {
	url := fmt.Sprintf("%s/latest.json?app_id=%s", e.baseURL, e.appID)
	return e.fetch(ctx, url)
}

func (e *exchangeClient) Historical(ctx context.Context, yyyymmdd string) (map[string]float64, error) {
	url := fmt.Sprintf("%s/historical/%s.json?app_id=%s", e.baseURL, yyyymmdd, e.appID)
	return e.fetch(ctx, url)
}

func (e *exchangeClient) fetch(ctx context.Context, url string) (map[string]float64, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := e.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("oxr status %d", resp.StatusCode)
	}
	var data oxrResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}
	return data.Rates, nil
}
