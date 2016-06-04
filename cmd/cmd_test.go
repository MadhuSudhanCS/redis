package cmd

import (
	"strings"
	"testing"
)

func TestEmptyReqNewCommand(t *testing.T) {
	req := ""

	if _, err := NewCommand(req); err == nil {
		t.Fatalf("expected to thorw error")
	}
}

func TestNewCommand(t *testing.T) {
	req := "get     mykey_0"

	cmd, err := NewCommand(req)
	if err != nil {
		t.Fatalf("failed to parse command. err: %v", err)
	}

	if strings.Compare(cmd.Name, "GET") != 0 {
		t.Fatalf("expected: %s found: %s", "get", cmd.Name)
	}

	if len(cmd.Options) != 1 {
		t.Fatalf("invalid options: %v", cmd.Options)
	}
}
