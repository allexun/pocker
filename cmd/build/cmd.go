package build

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"gitlab.com/kritskov/pocker/internal/builder"
)

func Command() *cli.Command {
	return &cli.Command{
		Name:  "build",
		Usage: "Build docker image",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "push",
				Usage: "Push the built image to the registry",
			},
			&cli.StringSliceFlag{
				Name:    "tag",
				Usage:   "Additional tags for the Docker image",
				Aliases: []string{"t"},
			},
		},
		Action: buildAction,
	}
}

func buildAction(ctx *cli.Context) error {
	phpVersion := ctx.String("php-version")
	if phpVersion == "" {
		return fmt.Errorf("php-version is required")
	}

	return builder.Build(ctx.Context, &builder.Options{
		PhpVersion:      phpVersion,
		ComposerVersion: ctx.Int("composer-version"),
		Push:            ctx.Bool("push"),
		Tags:            ctx.StringSlice("tag"),
	})
}
