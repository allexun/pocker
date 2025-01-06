package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"
	"gitlab.com/kritskov/pocker/cmd/build"
	"gitlab.com/kritskov/pocker/cmd/install"
	"gitlab.com/kritskov/pocker/internal/common"
)

func main() {
	app := &cli.App{
		Name:    "pocker",
		Usage:   "Pocker is a tool for installing dependencies in PHP projects using Composer in Docker",
		Version: "0.1.0",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "image",
				Usage: "Base image name to use",
				Value: common.ImageBaseName,
			},
			&cli.IntFlag{
				Name:  "composer-version",
				Usage: "Composer version to use",
				Value: 2,
			},
		},
		Commands: []*cli.Command{
			build.Command(),
			install.Command(),
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
