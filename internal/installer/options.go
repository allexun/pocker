package installer

type Options struct {
	ProjectPath     string
	ComposerVersion int
	UseSsh          bool
	NoAutoRemove    bool
	Cmd             string
}
