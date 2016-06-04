package set

import (
	"github.com/madhusudhancs/redis/cache"
	"github.com/madhusudhancs/redis/cmd"
	"strconv"
	"strings"
	"time"
	"github.com/madhusudhancs/redis/utils/log"
)

const (
	CmdName = "SET"
)

func init() {
	if err := cmd.Register(CmdName, Run); err != nil {
		log.Errorf("SetCmd: failed to register command. err: %v", err)
	}
}

func Run(options []string, c *cache.Cache) ([]byte, bool) {
	if len(options) < 2 {
		return cmd.GetErrMsg("ERR wrong number of arguments for 'set' command"), false
	}

	index := 0
	length := len(options)
	foundEX := false
	foundPX := false
	input := cache.SetInput{}

	input.Key = options[index]
	index = index + 1

	input.Value = options[index]
	index = index + 1

	for ; index < length; index++ {
		if foundEX {
			i, err := strconv.Atoi(options[index])
			if err != nil {
				return cmd.GetErrMsg("ERR value is not an integer or out of range"), false
			}

			input.Expiry = time.Duration(i) * time.Second
			foundEX = false
			continue
		}

		if foundPX {
			i, err := strconv.Atoi(options[index])
			if err != nil {
				return cmd.GetErrMsg("ERR value is not an integer or out of range"), false
			}

			input.Expiry = time.Duration(i) * time.Millisecond
			foundPX = false
			continue
		}

		if strings.EqualFold(options[index], "EX") {
			foundEX = true
			continue
		}

		if strings.EqualFold(options[index], "PX") {
			foundPX = true
			continue
		}

		if strings.EqualFold(options[index], "NX") {
			input.NX = true
			continue
		}

		if strings.EqualFold(options[index], "XX") {
			input.XX = true
			continue
		}

		return cmd.GetErrMsg("ERR syntax error"), false
	}

	if foundEX || foundPX {
		return cmd.GetErrMsg("ERR syntax error"), false
	}

	err := c.Set(input)
	if err != nil {
		log.Errorf("SetCmd: err: %v", err)
		return cmd.GetNullReply(), false
	}

	return cmd.GetSimpleString("OK"), false
}
