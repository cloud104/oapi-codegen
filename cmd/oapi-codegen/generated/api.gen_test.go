package api_test

import (
	"testing"

	api "github.com/cloud104/oapi-codegen/cmd/oapi-codegen/generated"
)

func Test_ListBeers(t *testing.T) {
	ctx := t.Context()

	client, err := api.NewClient("http://localhost:8080/")
	if err != nil {
		t.Fatal(err)
	}

	response, err := client.ListBeers(ctx, &api.ListBeersParams{})
	if err != nil {
		t.Fatal(err)
	}

	link, _ := response.GetHeader[string]("Link")
	total, _ := response.GetHeader[int]("X-Total-Count")
	retryAfter, _ := response.GetHeader[int]("Retry-After")

	_ = response
	_ = link
	_ = total
	_ = retryAfter
}
