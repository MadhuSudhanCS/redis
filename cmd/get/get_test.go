package get

import (
	"github.com/redis/cache"
	"strings"
	"testing"
)

func TestInvalidArgs(t *testing.T) {
	c := cache.NewCache(2)
	cmds := []string{"", "key key"}

	for _, cmd := range cmds {
		options := strings.Fields(cmd)

		if resp, _ := Run(options, c); strings.Contains(string(resp), "$") {
			t.Fatalf("expected err msg found: %s", resp)
		}
	}
}

func TestValidArgs(t *testing.T) {
	c := cache.NewCache(2)
	input := cache.SetInput{
		Key:   "key",
		Value: "value",
	}
	c.Set(input)

	cmds := []string{"key"}
	for _, cmd := range cmds {
		options := strings.Fields(cmd)

		if resp, _ := Run(options, c); strings.Compare(string(resp), "$5\r\nvalue\r\n") != 0 {
			t.Fatalf("expected success msg found: %s", resp)
		}
	}
}

func TestNullReply(t *testing.T) {
	c := cache.NewCache(2)

	cmd := "key"
	options := strings.Fields(cmd)

	if resp, _ := Run(options, c); strings.Compare(string(resp), "$-1\r\n") != 0 {
		t.Fatalf("expected null msg found: %s", resp)
	}
}
