package setbit

import (
	"github.com/redis/cache"
	"strings"
	"testing"
)

func TestInvalidArgs(t *testing.T) {
	c := cache.NewCache(2)
	cmds := []string{"", "key", "key 7", "key a 0", "key a b"}

	for _, cmd := range cmds {
		options := strings.Fields(cmd)

		if resp, _ := Run(options, c); string(resp) == "+OK\r\n" {
			t.Fatalf("expected err msg found: %s", resp)
		}
	}
}

func TestValidArgs(t *testing.T) {
	c := cache.NewCache(2)
	cmds := []string{"key 7 0", "key 512 0", "mykey 200 1"}
	for _, cmd := range cmds {
		options := strings.Fields(cmd)

		if resp, _ := Run(options, c); !strings.Contains(string(resp), ":") {
			t.Fatalf("expected success msg found: %s", resp)
		}
	}
}
