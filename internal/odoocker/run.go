package odoocker

import (
	"io"
	"os"
	"slices"

	"github.com/abichinger/odoocker/flags"
	"github.com/leaanthony/clir"
)

type Run struct {
	*Odoocker
	options *flags.Run
}

func NewRun(options *flags.Run) (*Run, error) {
	r := &Run{
		Odoocker: NewOdoocker(options.Common),
		options:  options,
	}

	docker, err := NewDockerClient(r)
	if err != nil {
		return r, err
	}
	r.docker = docker
	return r, nil
}

func (r *Run) cmdOut() io.Writer {
	return os.Stdout
}

func (r *Run) Run() error {
	if err := r.Setup(); err != nil {
		return err
	}
	defer r.Teardown()

	i := slices.IndexFunc(os.Args, func(arg string) bool { return arg == "run" })
	args := append([]string{web}, os.Args[i+1:]...)
	r.docker.ExecAs("odoo", args...)

	return nil

}

func AddRunCommand(cli *clir.Cli) {

	cmd := cli.NewSubCommand("run", "Runs a one-time command inside the odoo container")
	common := &flags.Common{}
	cmd.AddFlags(common)
	flags := &flags.Run{
		Common: common,
	}
	cmd.AddFlags(flags)

	cmd.Action(func() error {
		flags.Apply()
		action, err := NewRun(flags)
		if err != nil {
			return err
		}
		return action.Run()
	})

}
