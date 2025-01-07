package build

import (
	"github.com/urfave/cli/v2"
	"gitlab.com/kritskov/pocker/internal/builder"
)

func Command() *cli.Command {
	return &cli.Command{
		Name:  "build",
		Usage: "Build docker image",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "php-version",
				Usage:    "PHP version for the image",
				Value:    "",
				Aliases:  []string{"v"},
				Required: true,
			},
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
	return builder.Build(ctx.Context, &builder.Options{
		PhpVersion:      ctx.String("php-version"),
		ComposerVersion: ctx.Int("composer-version"),
		Push:            ctx.Bool("push"),
		Tags:            ctx.StringSlice("tag"),
	})
}
