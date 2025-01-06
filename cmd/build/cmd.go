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
		},
		Action: buildAction,
	}
}

func buildAction(ctx *cli.Context) error {
	return builder.Build(ctx.Context, &builder.Options{
		PhpVersion:      ctx.String("php-version"),
		ComposerVersion: ctx.Int("composer-version"),
	})
}
