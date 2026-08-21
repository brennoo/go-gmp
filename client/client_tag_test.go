package client

import (
	"context"
	"testing"

	"github.com/brennoo/go-gmp/commands"
)

func TestCreateTag(t *testing.T) {
	cli := New(MockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	cmd := &commands.CreateTag{
		Name: "geo:long",
		Resources: commands.TagResources{
			Type: "target",
			Resource: []commands.TagResource{
				{ID: "b493b7a8-7489-11df-a3ec-002264764cea"},
			},
		},
		Value: "52.2788",
	}
	resp, err := cli.CreateTag(context.Background(), cmd)
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

func TestModifyTag(t *testing.T) {
	cli := New(MockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	cmd := &commands.ModifyTag{
		TagID:  "254cd3ef-bbe1-4d58-859d-21b8d0c046c6",
		Active: commands.Bool(false),
	}
	resp, err := cli.ModifyTag(context.Background(), cmd)
	if err != nil {
		t.Fatalf("Unexpected error during ModifyTag: %s", err)
	}
	if resp.Status != "200" {
		t.Errorf("Expected status 200, got %s", resp.Status)
	}
}

func TestGetTags(t *testing.T) {
	cli := New(MockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	ctx := context.Background()
	resp, err := cli.GetTags(ctx, "name=test-tag")
	if err != nil {
		t.Fatalf("Unexpected error during GetTags: %s", err)
	}
	if resp.Status != "200" {
		t.Errorf("Expected status 200, got %s", resp.Status)
	}
	if len(resp.Tag) != 1 || resp.Tag[0].Name != "geo:long" {
		t.Errorf("Expected tag 'geo:long', got %+v", resp.Tag)
	}
}

func TestDeleteTag(t *testing.T) {
	cli := New(MockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	cmd := &commands.DeleteTag{
		TagID:    "254cd3ef-bbe1-4d58-859d-21b8d0c046c6",
		Ultimate: true,
	}
	resp, err := cli.DeleteTag(context.Background(), cmd)
	if err != nil {
		t.Fatalf("Unexpected error during DeleteTag: %s", err)
	}
	if resp.Status != "200" {
		t.Errorf("Expected status 200, got %s", resp.Status)
	}
}
