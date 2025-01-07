package installer

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"regexp"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"gitlab.com/kritskov/pocker/internal/common"
	"gitlab.com/kritskov/pocker/internal/composer"
)

func Install(ctx context.Context, options *Options) error {
	composerFile, err := composer.Parse(options.ProjectPath)
	if err != nil {
		return err
	}

	version, err := composerFile.GetPhpVersion()
	if err != nil {
		return err
	}
	fmt.Printf("Using PHP version '%s' with composer %d\n", version, options.ComposerVersion)

	image := fmt.Sprintf("%s:%s-%d", common.ImageBaseName, version, options.ComposerVersion)
	docker, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return err
	}
	fmt.Printf("Using image: %s\n", image)

	cmd := []string{"composer", "install", "--optimize-autoloader", "--ignore-platform-reqs", "--no-interaction"}
	if options.Cmd != "" {
		cmd = regexp.MustCompile(`\s+`).Split(options.Cmd, -1)
	}

	binds, err := getBinds(options)
	if err != nil {
		return err
	}

	c, err := docker.ContainerCreate(ctx, &container.Config{
		Image:      image,
		Cmd:        cmd,
		WorkingDir: "/app",
		Env: []string{
			"GIT_SSH_COMMAND=ssh -o StrictHostKeyChecking=accept-new",
		},
	}, &container.HostConfig{
		AutoRemove: !options.NoAutoRemove,
		Binds:      binds,
	}, nil, nil, "")
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

func getBinds(options *Options) ([]string, error) {
	binds := []string{
		fmt.Sprintf("%s:/app", options.ProjectPath),
	}
	if options.UseSsh {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, err
		}
		binds = append(binds, fmt.Sprintf("%s/.ssh/id_rsa:/root/.ssh/id_rsa:ro", home), fmt.Sprintf("%s/.ssh/id_rsa.pub:/root/.ssh/id_rsa.pub:ro", home))
	}

	return binds, nil
}
