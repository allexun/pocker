package install

import (
	"os"

	"github.com/urfave/cli/v2"
	"gitlab.com/kritskov/pocker/internal/installer"
)

func Command() *cli.Command {
	return &cli.Command{
		Name:  "install",
		Usage: "Install composer requirements",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "path",
				Usage: "Path to the project",
			},
			&cli.BoolFlag{
				Name:    "ssh",
				Usage:   "Copy user's SSH key to the container",
				Aliases: []string{"s"},
			},
			&cli.BoolFlag{
				Name:  "no-auto-remove",
				Usage: "Do not remove container after execution",
			},
			&cli.StringFlag{
				Name:  "cmd",
				Usage: "Specify custom command to run in the container",
			},
		},
		Action: installAction,
	}
}

func installAction(ctx *cli.Context) error {
	projectPath := ctx.String("path")
	if projectPath == "" {
		wd, err := os.Getwd()
		if err != nil {
			return err
		}
		projectPath = wd
	}

	return installer.Install(ctx.Context, &installer.Options{
		ProjectPath:     projectPath,
		PhpVersion:      ctx.String("php-version"),
		ComposerVersion: ctx.Int("composer-version"),
		UseSsh:          ctx.Bool("ssh"),
		NoAutoRemove:    ctx.Bool("no-auto-remove"),
		Cmd:             ctx.String("cmd"),
	})
}
