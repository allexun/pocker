package installer

import (
	"bufio"
	"context"
	"fmt"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"gitlab.com/kritskov/pocker/internal/common"
	"gitlab.com/kritskov/pocker/internal/composer"
)

func Install(ctx context.Context, projectPath string, composerVersion int) error {
	composerFile, err := composer.Parse(projectPath)
	if err != nil {
		return err
	}

	version, err := composerFile.GetPhpVersion()
	if err != nil {
		return err
	}

	image := fmt.Sprintf("%s:%s-%d", common.ImageBaseName, version, composerVersion)
	docker, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return err
	}

	c, err := docker.ContainerCreate(ctx, &container.Config{
		Image: image,
		Volumes: map[string]struct{}{
			"/app": {},
		},
		Cmd:        []string{"composer", "install", "--optimize-autoloader", "--ignore-platform-reqs", "--no-interaction"},
		WorkingDir: "/app",
	}, &container.HostConfig{
		AutoRemove: true,
		Binds: []string{
			fmt.Sprintf("%s:/app", projectPath),
		},
	}, nil, nil, "php-dev")
	if err != nil {
		return err
	}
	fmt.Println(c.ID)

	if err = docker.ContainerStart(ctx, c.ID, container.StartOptions{}); err != nil {
		return err
	}
	hijack, err := docker.ContainerAttach(ctx, c.ID, container.AttachOptions{
		Stream: true,
		Stdout: true,
		Stderr: true,
	})
	if err != nil {
		return err
	}
	defer hijack.Close()

	scanner := bufio.NewScanner(hijack.Reader)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}

	return nil
}
