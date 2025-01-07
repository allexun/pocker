package builder

import (
	"archive/tar"
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"text/template"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"gitlab.com/kritskov/pocker/internal/common"
	"gitlab.com/kritskov/pocker/internal/docker"
)

func Build(ctx context.Context, options *Options) error {
	tmpl, err := template.New(common.Dockerfile).Parse(docker.FileTemplate)
	if err != nil {
		return err
	}

	dockerfile := &bytes.Buffer{}
	err = tmpl.Execute(dockerfile, &docker.TemplateOptions{
		PhpVersion:      options.PhpVersion,
		ComposerVersion: options.ComposerVersion,
	})
	if err != nil {
		return err
	}

	dockerClient, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return err
	}

	dockerContext, err := CreateDockerContext(dockerfile)
	if err != nil {
		return err
	}

	img := fmt.Sprintf("%s:%s-%d", common.ImageBaseName, options.PhpVersion, options.ComposerVersion)
	tags := []string{img}
	for _, t := range options.Tags {
		tags = append(tags, fmt.Sprintf("%s:%s", common.ImageBaseName, t))
	}

	buildResponse, err := dockerClient.ImageBuild(ctx, dockerContext, types.ImageBuildOptions{
		Tags:    tags,
		Version: types.BuilderBuildKit,
	})
	if err != nil {
		return err
	}
	defer buildResponse.Body.Close()
	if _, err = io.Copy(os.Stdout, buildResponse.Body); err != nil {
		return err
	}

	if options.Push {
		for _, t := range tags {
			if err = docker.ExecPush(ctx, t); err != nil {
				return err
			}
		}
	}

	return nil
}

func CreateDockerContext(dockerfile *bytes.Buffer) (io.Reader, error) {
	buf := new(bytes.Buffer)
	tw := tar.NewWriter(buf)
	defer tw.Close()

	tarHeader := &tar.Header{
		Name: common.Dockerfile,
		Size: int64(dockerfile.Len()),
	}
	if err := tw.WriteHeader(tarHeader); err != nil {
		return nil, err
	}
	if _, err := tw.Write(dockerfile.Bytes()); err != nil {
		return nil, err
	}

	return bytes.NewReader(buf.Bytes()), nil
}
