package installer

type Options struct {
	ProjectPath     string
	PhpVersion      string
	ComposerVersion int
	UseSsh          bool
	NoAutoRemove    bool
	Cmd             string
}
