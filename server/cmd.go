package server

import (
	_ "github.com/redis/cmd/get"
	_ "github.com/redis/cmd/getbit"
	_ "github.com/redis/cmd/quit"
	_ "github.com/redis/cmd/save"
	_ "github.com/redis/cmd/set"
	_ "github.com/redis/cmd/setbit"
	_ "github.com/redis/cmd/zadd"
	_ "github.com/redis/cmd/zcard"
	_ "github.com/redis/cmd/zcount"
	_ "github.com/redis/cmd/zrange"
)
