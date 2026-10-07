package api

import (
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"time"
)

type RequestLogger struct {
	DefaultTransport http.RoundTripper
}

func (rl *RequestLogger) RoundTrip(r *http.Request) (*http.Response, error) {
	// Get the path and handle any URL decoding errors
	path, err := url.PathUnescape(r.URL.Path)
	if err != nil {
		path = r.URL.Path
	}

	// Get the raw query and handle any URL decoding errors
	rawQuery, err := url.QueryUnescape(r.URL.RawQuery)
	if err != nil {
		rawQuery = r.URL.RawQuery
	}

	// Combine path and raw query if there is a query string
	if rawQuery != "" {
		path += "?" + rawQuery
	}

	// Prepare the arguments for logging
	args := []any{
		"request", fmt.Sprintf("%s %s", r.Method, path),
	}

	startTime := time.Now()

	// Make the HTTP request using the default transport
	resp, err := rl.DefaultTransport.RoundTrip(r)
	if err != nil {
		return nil, err
	}

	endTime := time.Now()

	args = append(args,
		"response", resp.Status,
		"duration", endTime.Sub(startTime),
	)

	// Log the request details using structured logging
	ctx := r.Context()
	slog.InfoContext(ctx, "call to API gateway", args...)

	return resp, nil
}
