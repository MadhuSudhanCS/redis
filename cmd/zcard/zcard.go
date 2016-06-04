package zcard

import (
	"github.com/redis/cache"
	"github.com/redis/cmd"
	"github.com/redis/utils/log"
)

const (
	CmdName = "ZCARD"
)

func init() {
	if err := cmd.Register(CmdName, Run); err != nil {
		log.Errorf("ZCardCmd: failed to register command. err: %v", err)
	}
}

func Run(options []string, c *cache.Cache) ([]byte, bool) {
	if len(options) != 1 {
		return cmd.GetErrMsg("ERR wrong number of arguments for 'zcard' command"), false
	}

	resp, err := c.ZCard(options[0])
	if err != nil {
		log.Errorf("ZCardCmd: err: %v", err)
		return cmd.GetErrMsg(err), false
	}

	return cmd.GetIntegerMsg(resp), false
}
