package save

import (
	"github.com/redis/cache"
	"github.com/redis/cmd"
	"github.com/redis/config"
	"github.com/redis/utils/log"
)

const (
	CmdName = "SAVE"
)

func init() {
	if err := cmd.Register(CmdName, Run); err != nil {
		log.Errorf("SaveCmd: failed to register command. err: %v", err)
	}
}

func Run(options []string, c *cache.Cache) ([]byte, bool) {
	if _, err := c.Save(config.DBFileName); err != nil {
		log.Errorf("SaveCmd: err: %v", err)
		return cmd.GetErrMsg(err), false
	}

	return cmd.GetSimpleString("OK"), false
}
