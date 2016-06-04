package zrange

import (
	"github.com/madhusudhancs/redis/cache"
	"github.com/madhusudhancs/redis/store"
	"strings"
	"testing"
)

func TestInvalidArgs(t *testing.T) {
	c := cache.NewCache(2)
	cmds := []string{"", "key", "key a b", "key 0 max"}

	for _, cmd := range cmds {
		options := strings.Fields(cmd)

		if resp, _ := Run(options, c); strings.Contains(string(resp), "$") {
			t.Fatalf("expected err msg found: %s", resp)
		}
	}
}

func TestValidArgs(t *testing.T) {
	c := cache.NewCache(2)
	input := cache.ZSetInput{
		Key: "key",
		Scores: []store.ScoreMember{
			{
				Score:  "10",
				Member: "Java",
			},
		},
	}

	c.ZAdd(input)

	cmds := []string{"key 0 10", "key 0 -1", "key -2 -1"}
	for _, cmd := range cmds {
		options := strings.Fields(cmd)

		if resp, _ := Run(options, c); !strings.Contains(string(resp), "*") {
			t.Fatalf("expected success msg found: %s", resp)
		}
	}
}
