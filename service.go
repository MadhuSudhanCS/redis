package main

import (
	"github.com/codegangsta/cli"
	"github.com/redis/config"
	"github.com/redis/server"
	"github.com/redis/utils/log"
	"os"
)

func start() {
	s, err := server.NewServer()
	if err != nil {
		log.Fatalf("ExoRedisServer: failed to initialize server. err: %v", err)
	}

	s.Start()
}

func main() {
	var fileName string
	app := cli.NewApp()
	app.Name = "ExoRedis Server"
	app.Usage = "In-Memory Data Store"
	app.Version = "1.0.0"
	app.EnableBashCompletion = true
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "dbfile, f",
			Usage:       "File to save data store",
			Destination: &fileName,
		},
	}

	app.Action = func(c *cli.Context) {
		if fileName == "" {
			log.Errorf("DB file name not specified")
			cli.ShowCommandHelp(c, "")
			return
		}

		config.DBFileName = fileName
		start()
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatalf("Failed to run ExoRedis Server. err: %v", err)
	}
}
