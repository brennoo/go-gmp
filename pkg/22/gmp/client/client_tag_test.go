package client

import (
	"testing"

	"github.com/brennoo/go-gmp/pkg/22/gmp"
)

func TestCreateTag(t *testing.T) {
	cli := New(mockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	cmd := &gmp.CreateTagCommand{
		Name: "geo:long",
		Resources: &gmp.Resources{
			Resource: &gmp.Resource{ID: "b493b7a8-7489-11df-a3ec-002264764cea"},
			Type:     "target",
		},
		Value: "52.2788",
	}
	resp, err := cli.CreateTag(cmd)
	if err != nil {
		t.Fatalf("Unexpected error during CreateTag: %s", err)
	}
	if resp.Status != "201" {
		t.Errorf("Expected status 201, got %s", resp.Status)
	}
	if resp.ID != "254cd3ef-bbe1-4d58-859d-21b8d0c046c6" {
		t.Errorf("Expected tag ID '254cd3ef-bbe1-4d58-859d-21b8d0c046c6', got %s", resp.ID)
	}
}
