package docker

import (
	"context"
	"os"
	"os/exec"
)

func ExecPush(ctx context.Context, image string) error {
	return dockerExecImageOperation(ctx, "push", image)
}

func ExecPull(ctx context.Context, image string) error {
	return dockerExecImageOperation(ctx, "pull", image)
}

func dockerExecImageOperation(ctx context.Context, operation, image string) error {
	cmd := exec.CommandContext(ctx, "docker", operation, image)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()

	return cmd.Run()
}
