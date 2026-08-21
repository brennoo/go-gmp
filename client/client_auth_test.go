package client

import (
	"context"
	"testing"

	"github.com/brennoo/go-gmp/commands"
)

func TestAuthenticate(t *testing.T) {
	cli := New(MockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	cmd := &commands.Authenticate{
		Credentials: commands.AuthCredentials{
			Username: "openvas",
			Password: "123",
		},
	}
	_, err := cli.Authenticate(context.Background(), cmd)
	if err != nil {
		t.Fatalf("Unexpected error during Authenticate: %s", err)
	}
}

func TestDescribeAuth(t *testing.T) {
	cli := New(MockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	cmd := &commands.DescribeAuth{}
	resp, err := cli.DescribeAuth(context.Background(), cmd)
	if err != nil {
		t.Fatalf("Unexpected error during DescribeAuth: %s", err)
	}
	if len(resp.Group) != 1 {
		t.Fatalf("Expected 1 group, got %d", len(resp.Group))
	}
	group := resp.Group[0]
	if group.Name != "method:file" {
		t.Fatalf("Unexpected group name. Expected: method:file Got: %s", group.Name)
	}
	if len(group.AuthConfSetting) != 2 {
		t.Fatalf("Expected 2 auth_conf_setting, got %d", len(group.AuthConfSetting))
	}
	if group.AuthConfSetting[0].Key != "enable" || group.AuthConfSetting[0].Value != "true" {
		t.Fatalf("Unexpected first setting: %+v", group.AuthConfSetting[0])
	}
	if group.AuthConfSetting[1].Key != "order" || group.AuthConfSetting[1].Value != "1" {
		t.Fatalf("Unexpected second setting: %+v", group.AuthConfSetting[1])
	}
}

func TestModifyAuth(t *testing.T) {
	cli := New(MockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	cmd := &commands.ModifyAuth{
		Group: commands.ModifyAuthGroup{
			Name: "method:file",
			AuthConfSetting: []commands.AuthGroupAuthConfSetting{
				{Key: "enable", Value: "true"},
				{Key: "order", Value: "1"},
			},
		},
	}
	_, err := cli.ModifyAuth(context.Background(), cmd)
	if err != nil {
		t.Fatalf("Unexpected error during ModifyAuth: %s", err)
	}
}
