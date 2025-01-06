package builder

import (
	"archive/tar"
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"text/template"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"gitlab.com/kritskov/pocker/internal/common"
	"gitlab.com/kritskov/pocker/internal/docker"
)

func Build(ctx context.Context, opts *Options) error {
	tmpl, err := template.New(common.Dockerfile).Parse(docker.FileTemplate)
	if err != nil {
		return err
	}

	dockerfile := &bytes.Buffer{}
	err = tmpl.Execute(dockerfile, &docker.TemplateOptions{
		PhpVersion:      opts.PhpVersion,
		ComposerVersion: opts.ComposerVersion,
	})
	if err != nil {
		return err
	}

	docker, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return err
	}

	dockerContext, err := CreateDockerContext(dockerfile)
	if err != nil {
		return err
	}

	res, err := docker.ImageBuild(ctx, dockerContext, types.ImageBuildOptions{
		Tags: []string{fmt.Sprintf("%s:%s-%d", common.ImageBaseName, opts.PhpVersion, opts.ComposerVersion)},
	})
	if err != nil {
		return err
	}
	defer res.Body.Close()
	reader := bufio.NewReader(res.Body)
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
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
