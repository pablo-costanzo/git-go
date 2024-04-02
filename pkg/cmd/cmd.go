package cmd

type Executable interface {
	Name() string
	Exec() error
	Help() string
}
