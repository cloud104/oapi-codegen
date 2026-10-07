package api

import (
	"net/http"

	cleanhttp "github.com/hashicorp/go-cleanhttp"
)

func NewTraceableClient() *http.Client {
	return &http.Client{
		Transport: &RequestLogger{
			DefaultTransport: cleanhttp.DefaultTransport(),
		},
	}
}
