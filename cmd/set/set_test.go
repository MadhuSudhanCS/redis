package set

import (
	"github.com/madhusudhancs/redis/cache"
	"strings"
	"testing"
)

func TestInvalidArgs(t *testing.T) {
	c := cache.NewCache(2)
	cmds := []string{"", "key", "key value ex", "key value ex ex", "key value 10 ex", "key value ggg"}

	for _, cmd := range cmds {
		options := strings.Fields(cmd)

		if resp, _ := Run(options, c); string(resp) == "+OK\r\n" {
			t.Fatalf("expected err msg found: %s", resp)
		}
	}
}

func TestValidArgs(t *testing.T) {
	c := cache.NewCache(2)
	cmds := []string{"key value", "key value ex 10", "key value ex 10 ex 100", "key value ex 10 px 3", "key value"}
	for _, cmd := range cmds {
		options := strings.Fields(cmd)

		if resp, _ := Run(options, c); string(resp) != "+OK\r\n" {
			t.Fatalf("expected success msg found: %s", resp)
		}
	}
}

func TestNullReply(t *testing.T) {
	c := cache.NewCache(2)

	cmd := "key value"
	options := strings.Fields(cmd)

	if resp, _ := Run(options, c); string(resp) != "+OK\r\n" {
		t.Fatalf("expected success msg found: %s", resp)
	}

	cmds := []string{"key value nx", "newkey value xx"}
	for _, cmd := range cmds {
		options := strings.Fields(cmd)

		if resp, _ := Run(options, c); string(resp) != "$-1\r\n" {
			t.Fatalf("expected null msg found: %s", resp)
		}
	}
}
