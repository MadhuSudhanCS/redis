package zcard

import (
	"github.com/redis/cache"
	"github.com/redis/store"
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

	cmds := []string{"key"}
	for _, cmd := range cmds {
		options := strings.Fields(cmd)

		if resp, _ := Run(options, c); strings.Compare(string(resp), ":1\r\n") != 0 {
			t.Fatalf("expected success msg found: %s", resp)
		}
	}
}
