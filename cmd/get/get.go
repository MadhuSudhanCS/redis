package get

import (
	"github.com/madhusudhancs/redis/cache"
	"github.com/madhusudhancs/redis/cmd"
	"github.com/madhusudhancs/redis/utils/log"
)

const (
	CmdName = "GET"
)

func init() {
	if err := cmd.Register(CmdName, Run); err != nil {
		log.Errorf("GetCmd: failed to register command. err: %v", err)
	}
}

func Run(options []string, c *cache.Cache) ([]byte, bool) {
	if len(options) != 1 {
		return cmd.GetErrMsg("ERR wrong number of arguments for 'get' command"), false
	}

	resp, err := c.Get(options[0])
	if err != nil {
		log.Errorf("GetCmd: err: %v", err)
		return cmd.GetErrMsg(err), false
	}

	if resp == "" {
		return cmd.GetNullReply(), false
	}

	return cmd.GetBulkString(resp), false
}
