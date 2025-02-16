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
		Usage:   "Pocker is a tool for working with PHP projects using Docker",
		Version: "0.2.0",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "image",
				Usage: "Base image name to use",
				Value: common.ImageBaseName,
			},
			&cli.StringFlag{
				Name:    "php-version",
				Usage:   "PHP version for the image (default: based on composer.json)",
				Value:   "",
				Aliases: []string{"p"},
			},
			&cli.IntFlag{
				Name:    "composer-version",
				Usage:   "Composer version to use",
				Aliases: []string{"c"},
				Value:   2,
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
