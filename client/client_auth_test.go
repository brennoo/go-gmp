package client

import (
	"testing"

	"github.com/brennoo/go-gmp/commands"
)

func TestAuthenticate(t *testing.T) {
	cli := New(mockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	cmd := &commands.Authenticate{
		Credentials: &commands.AuthenticateCredentials{
			Username: "openvas",
			Password: "123",
		},
	}
	resp, err := cli.Authenticate(cmd)
	if err != nil {
		t.Fatalf("Unexpected error during Authenticate: %s", err)
	}

	if resp.Status != "200" {
		t.Fatalf("Unexpected status. \nExpected: 200 \nGot: %s", resp.Status)
	}
}

func TestDescribeAuth(t *testing.T) {
	cli := New(mockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	cmd := &commands.DescribeAuth{}
	resp, err := cli.DescribeAuth(cmd)
	if err != nil {
		t.Fatalf("Unexpected error during DescribeAuth: %s", err)
	}
	if resp.Status != "200" {
		t.Fatalf("Unexpected status. Expected: 200 Got: %s", resp.Status)
	}
	if resp.StatusText != "OK" {
		t.Fatalf("Unexpected status text. Expected: OK Got: %s", resp.StatusText)
	}
	if len(resp.Groups) != 1 {
		t.Fatalf("Expected 1 group, got %d", len(resp.Groups))
	}
	group := resp.Groups[0]
	if group.Name != "method:file" {
		t.Fatalf("Unexpected group name. Expected: method:file Got: %s", group.Name)
	}
	if len(group.Settings) != 2 {
		t.Fatalf("Expected 2 auth_conf_setting, got %d", len(group.Settings))
	}
	if group.Settings[0].Key != "enable" || group.Settings[0].Value != "true" {
		t.Fatalf("Unexpected first setting: %+v", group.Settings[0])
	}
	if group.Settings[1].Key != "order" || group.Settings[1].Value != "1" {
		t.Fatalf("Unexpected second setting: %+v", group.Settings[1])
	}
}

func TestModifyAuth(t *testing.T) {
	cli := New(mockedConnection())
	if cli == nil {
		t.Fatalf("Client is nil")
	}

	cmd := &commands.ModifyAuth{
		Group: commands.AuthGroup{
			Name: "method:file",
			Settings: []commands.AuthSetting{
				{Key: "enable", Value: "true"},
				{Key: "order", Value: "1"},
			},
		},
	}
	resp, err := cli.ModifyAuth(cmd)
	if err != nil {
		t.Fatalf("Unexpected error during ModifyAuth: %s", err)
	}
	if resp.Status != "200" {
		t.Fatalf("Unexpected status. Expected: 200 Got: %s", resp.Status)
	}
	if resp.StatusText != "OK" {
		t.Fatalf("Unexpected status text. Expected: OK Got: %s", resp.StatusText)
	}
}
