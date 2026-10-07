package api_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"testing"

	api "github.com/cloud104/oapi-codegen/cmd/oapi-codegen/generated"
)

func Test_ListBeers(t *testing.T) {
	ctx := t.Context()

	httpClient := api.NewTraceableClient()

	client, err := api.NewClient(
		"http://localhost:8080/",
		api.WithAPIKeyAuth("token"),
		api.WithBearerAuth(&tokenProvider{httpClient: httpClient}),
		api.WithHTTPClient(httpClient),
	)
	if err != nil {
		t.Fatal(err)
	}

	response, err := client.Ping(ctx)
	if err != nil {
		t.Fatal(err)
	}

	pong, err := response.GetPong()
	if err != nil {
		t.Fatal(err)
	}

	_ = pong
}

type tokenProvider struct {
	mu         sync.Mutex
	token      string
	httpClient *http.Client
}

func (p *tokenProvider) GetToken(ctx context.Context) (string, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.token != "" {
		return p.token, nil
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost:8080/token", nil)
	if err != nil {
		return "", err
	}

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("token request failed: %s", resp.Status)
	}

	var result struct {
		Token string `json:"token"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if result.Token == "" {
		return "", fmt.Errorf("empty token returned")
	}

	p.token = result.Token

	return p.token, nil
}

func (p *tokenProvider) InvalidateToken(ctx context.Context) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.token = ""

	return nil
}
