package zcount

import (
	"github.com/madhusudhancs/redis/cache"
	"github.com/madhusudhancs/redis/cmd"
	"strconv"
	"github.com/madhusudhancs/redis/utils/log"
)

const (
	CmdName = "ZCOUNT"
)

func init() {
	if err := cmd.Register(CmdName, Run); err != nil {
		log.Errorf("ZCountCmd: failed to register command. err: %v", err)
	}
}

func Run(options []string, c *cache.Cache) ([]byte, bool) {
	if len(options) != 3 {
		return cmd.GetErrMsg("ERR wrong number of arguments for 'zcount' command"), false
	}

	if _, err := strconv.ParseFloat(options[1], 64); err != nil {
		return cmd.GetErrMsg("ERR min or max is not a float"), false
	}

	if _, err := strconv.ParseFloat(options[2], 64); err != nil {
		return cmd.GetErrMsg("ERR min or max is not a float"), false
	}

	resp, err := c.ZCount(options[0], options[1], options[2])
	if err != nil {
		log.Errorf("ZCountCmd: err: %v", err)
		return cmd.GetErrMsg(err), false
	}

	return cmd.GetIntegerMsg(resp), false
}
