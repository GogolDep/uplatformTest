package services

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"time"

	"uplatform/pkg/randutil"
)

type GiphyClient interface {
	RandomFromSearch(ctx context.Context, query string) (string, error)
}

type giphyClient struct {
	baseURL string
	apiKey  string
	client  *http.Client
}

type giphySearchResp struct {
	Data []struct {
		Images struct {
			Original struct {
				URL string `json:"url"`
			} `json:"original"`
		} `json:"images"`
	} `json:"data"`
}

func NewGiphyClient(baseURL, apiKey string, timeout time.Duration) GiphyClient {
	return &giphyClient{
		baseURL: baseURL,
		apiKey:  apiKey,
		client:  &http.Client{Timeout: timeout},
	}
}

func (g *giphyClient) RandomFromSearch(ctx context.Context, query string) (string, error) {
	u := fmt.Sprintf("%s/gifs/search?api_key=%s&q=%s&limit=50&rating=g",
		g.baseURL, url.QueryEscape(g.apiKey), url.QueryEscape(query))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return "", err
	}
	resp, err := g.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("giphy status %d", resp.StatusCode)
	}
	var data giphySearchResp
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", err
	}
	if len(data.Data) == 0 {
		return "", fmt.Errorf("no gifs found for %q", query)
	}
	i := randutil.Intn(len(data.Data))
	return data.Data[i].Images.Original.URL, nil
}

func init() { rand.Seed(time.Now().UnixNano()) }
