package getbit

import (
	"github.com/redis/cache"
	"github.com/redis/cmd"
	"strconv"
	"github.com/redis/utils/log"
)

const (
	CmdName = "GETBIT"
)

func init() {
	if err := cmd.Register(CmdName, Run); err != nil {
		log.Errorf("GetBitCmd: failed to register command. err: %v", err)
	}
}

func Run(options []string, c *cache.Cache) ([]byte, bool) {
	if len(options) != 2 {
		return cmd.GetErrMsg("ERR wrong number of arguments for 'getbit' command"), false
	}

	key := options[0]
	offset, err := strconv.Atoi(options[1])
	if err != nil {
		return cmd.GetErrMsg("ERR bit offset is not an integer or out of range"), false
	}

	resp := c.GetBit(key, offset)

	return cmd.GetIntegerMsg(resp), false
}
