package docker

import _ "embed"

//go:embed Dockerfile.gotmpl
var FileTemplate string

type TemplateOptions struct {
	PhpVersion      string
	ComposerVersion int
}
