package getbit

import (
	"github.com/redis/cache"
	"strings"
	"testing"
)

func TestInvalidArgs(t *testing.T) {
	c := cache.NewCache(2)
	cmds := []string{"", "key", "key a 0", "key a b"}

	for _, cmd := range cmds {
		options := strings.Fields(cmd)

		if resp, _ := Run(options, c); string(resp) == "+OK\r\n" {
			t.Fatalf("expected err msg found: %s", resp)
		}
	}
}

func TestValidArgs(t *testing.T) {
	c := cache.NewCache(2)
	cmds := []string{"key 7", "key 512 ", "mykey 2000"}
	for _, cmd := range cmds {
		options := strings.Fields(cmd)

		if resp, _ := Run(options, c); !strings.Contains(string(resp), ":") {
			t.Fatalf("expected success msg found: %s", resp)
		}
	}
}
