package zadd

import (
	"github.com/madhusudhancs/redis/cache"
	"strings"
	"testing"
)

func TestInvalidArgs(t *testing.T) {
	c := cache.NewCache(2)
	cmds := []string{"", "key", "key incr nx", "incr nx xx ch", "incr nx ch", "incr 0 go a python", "abc def"}

	for _, cmd := range cmds {
		options := strings.Fields(cmd)

		if resp, _ := Run(options, c); string(resp) == "+OK\r\n" {
			t.Fatalf("expected err msg found: %s", resp)
		}
	}
}

func TestValidArgs(t *testing.T) {
	c := cache.NewCache(2)
	cmds := []string{"key 0 python", "key 10 go", "key nx 10 java", "key ch -1 c++"}
	for _, cmd := range cmds {
		options := strings.Fields(cmd)

		if resp, _ := Run(options, c); !strings.Contains(string(resp), ":") {
			t.Fatalf("expected success msg found: %s", resp)
		}
	}
}
