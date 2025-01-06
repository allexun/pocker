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
				Name:    "path",
				Usage:   "Path to the project",
				Value:   "",
				Aliases: []string{"p"},
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

	return installer.Install(ctx.Context, projectPath, ctx.Int("composer-version"))
}
