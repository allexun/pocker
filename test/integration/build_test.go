package integration

import (
	"context"
	"fmt"
	"gitlab.com/kritskov/pocker/internal/builder"
	"os/exec"
	"testing"
)

func TestBuildCommandWhenImagesNotPresent(t *testing.T) {
	ctx := context.Background()
	phpVersion := "7.4"
	composerVersion := 2
	phpImage := fmt.Sprintf("php:%s-cli", phpVersion)
	composerImage := fmt.Sprintf("composer:%d", composerVersion)

	if err := exec.CommandContext(ctx, "docker", "image", "rm", "-f", phpImage).Run(); err != nil {
		t.Fatal(err)
	}
	if err := exec.CommandContext(ctx, "docker", "image", "rm", "-f", composerImage).Run(); err != nil {
		t.Fatal(err)
	}

	err := builder.Build(ctx, &builder.Options{
		PhpVersion:      phpVersion,
		ComposerVersion: composerVersion,
	})
	if err != nil {
		t.Fatal(err)
	}
}
