package examples

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_CreateUser(t *testing.T) {
	options, err := HandlerWithOptions(&Unimplemented{}, ChiServerOptions{})
	if err != nil {
		t.Fatalf("HandlerWithOptions() failed: %v", err)
	}
	if options == nil {
		t.Fatal("HandlerWithOptions() returned a nil handler")
	}

	req := httptest.NewRequest(http.MethodPost, "/users", nil)
	rec := httptest.NewRecorder()

	options.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotImplemented {
		t.Errorf("expected status %d, got %d", http.StatusNotImplemented, rec.Code)
	}
}
