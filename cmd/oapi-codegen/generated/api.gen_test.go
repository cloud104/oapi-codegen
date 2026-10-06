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

	// response.Headers200

	_ = response
}
