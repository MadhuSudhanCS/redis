package zrange

import (
	"github.com/redis/cache"
	"github.com/redis/cmd"
	"strconv"
	"strings"
	"github.com/redis/utils/log"
)

const (
	CmdName = "ZRANGE"
)

func init() {
	if err := cmd.Register(CmdName, Run); err != nil {
		log.Errorf("ZRangeCmd: failed to register command. err: %v", err)
	}
}

func Run(options []string, c *cache.Cache) ([]byte, bool) {
	if len(options) != 3 && len(options) != 4 {
		return cmd.GetErrMsg("ERR wrong number of arguments for 'zrange' command"), false
	}

	start, err := strconv.Atoi(options[1])
	if err != nil {
		return cmd.GetErrMsg("ERR value is not an integer or out of range"), false
	}

	stop, err := strconv.Atoi(options[2])
	if err != nil {
		return cmd.GetErrMsg("ERR value is not an integer or out of range"), false
	}

	withScores := len(options) == 4
	if len(options) == 4 && !strings.EqualFold(options[3], "WITHSCORES") {
		return cmd.GetErrMsg("ERR syntax error"), false
	}

	scores, err := c.ZRange(options[0], start, stop)
	if err != nil {
		log.Errorf("ZRangeCmd: err: %v", err)
		return cmd.GetErrMsg(err), false
	}

	resp := make([]interface{}, 0)

	for _, score := range scores {
		resp = append(resp, score.Member)

		if withScores {
			resp = append(resp, score.Score)
		}
	}

	return cmd.GetArrayMsg(resp), false
}
