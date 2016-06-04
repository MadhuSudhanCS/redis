package zadd

import (
	"github.com/madhusudhancs/redis/cache"
	"github.com/madhusudhancs/redis/cmd"
	"github.com/madhusudhancs/redis/store"
	"fmt"
	"strconv"
	"strings"
	"github.com/madhusudhancs/redis/utils/log"
)

const (
	CmdName = "ZADD"
)

func init() {
	if err := cmd.Register(CmdName, Run); err != nil {
		log.Errorf("ZAddCmd: failed to register command. err: %v", err)
	}
}

func Run(options []string, c *cache.Cache) ([]byte, bool) {
	if len(options) < 3 {
		return cmd.GetErrMsg("ERR wrong number of arguments for 'zadd' command"), false
	}

	input := cache.ZSetInput{}
	length := len(options)
	index := 0
	foundNX := false
	foundXX := false
	input.Key = options[0]
	index = index + 1

	for ; index < length; index++ {
		if strings.EqualFold(options[index], "NX") {
			input.NX = true
			foundNX = true
			continue
		}

		if strings.EqualFold(options[index], "XX") {
			input.XX = true
			foundXX = true
			continue
		}

		if strings.EqualFold(options[index], "INCR") {
			input.INCR = true
			continue
		}

		if strings.EqualFold(options[index], "CH") {
			input.CH = true
			continue
		}

		break
	}

	if foundNX && foundXX {
		return cmd.GetErrMsg("ERR XX and NX options at the same time are not compatible"), false
	}

	input.Scores = []store.ScoreMember{}

	for ; index < length; index++ {
		score := options[index]

		if _, err := strconv.ParseFloat(score, 64); err != nil {
			return cmd.GetErrMsg("ERR value is not a valid float"), false
		}

		index = index + 1
		if index >= length {
			return cmd.GetErrMsg("ERR wrong number of arguments for 'zadd' command"), false
		}

		scoreMember := store.ScoreMember{
			Score:  score,
			Member: options[index],
		}

		input.Scores = append(input.Scores, scoreMember)
	}

	if len(input.Scores) == 0 {
		return cmd.GetErrMsg("ERR wrong number of arguments for 'zadd' command"), false
	}

	if len(input.Scores) > 1 && input.INCR {
		return cmd.GetErrMsg("ERR INCR option supports a single increment-element pair"), false
	}

	resp, err := c.ZAdd(input)
	if input.INCR && err != nil {
		log.Errorf("ZAddCmd: err: %v", err)
		return cmd.GetNullReply(), false
	}

	if err != nil {
		log.Errorf("ZAddCmd: err: %v", err)
		return cmd.GetErrMsg(err), false
	}

	if input.INCR {
		return cmd.GetBulkString(fmt.Sprintf("%d", resp)), false
	}

	return cmd.GetIntegerMsg(resp), false
}
