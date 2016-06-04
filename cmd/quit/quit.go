package quit

import (
	"github.com/madhusudhancs/redis/cache"
	"github.com/madhusudhancs/redis/cmd"
	"github.com/madhusudhancs/redis/utils/log"
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
