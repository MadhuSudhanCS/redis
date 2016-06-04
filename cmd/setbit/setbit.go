package setbit

import (
	"github.com/redis/cache"
	"github.com/redis/cmd"
	"strconv"
	"github.com/redis/utils/log"
)

const (
	CmdName = "SETBIT"
)

func init() {
	if err := cmd.Register(CmdName, Run); err != nil {
		log.Errorf("SetBitCmd: failed to register command. err: %v", err)
	}
}

func Run(options []string, c *cache.Cache) ([]byte, bool) {
	if len(options) != 3 {
		return cmd.GetErrMsg("ERR wrong number of arguments for 'setbit' command"), false
	}

	key := options[0]
	offset, err := strconv.Atoi(options[1])
	if err != nil {
		return cmd.GetErrMsg("ERR bit offset is not an integer or out of range"), false
	}

	value, err := strconv.Atoi(options[2])
	if err != nil || (value != 0 && value != 1) {
		return cmd.GetErrMsg("ERR bit is not an integer or out of range"), false
	}

	resp, err := c.SetBit(key, offset, value)
	if err != nil {
		log.Errorf("SetBitCmd: err: %v", err)
		return cmd.GetErrMsg(err), false
	}

	return cmd.GetIntegerMsg(resp), false

}
