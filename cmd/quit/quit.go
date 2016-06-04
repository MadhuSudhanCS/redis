package quit

import (
	"github.com/redis/cache"
	"github.com/redis/cmd"
	"github.com/redis/utils/log"
)

const (
	CmdName = "QUIT"
)

func init() {
	if err := cmd.Register(CmdName, Run); err != nil {
		log.Errorf("QuitCmd: failed to register command. err: %v", err)
	}
}

func Run(options []string, cache *cache.Cache) ([]byte, bool) {
	return cmd.GetSimpleString("OK"), true
}
