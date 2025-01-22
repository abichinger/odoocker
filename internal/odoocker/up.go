package odoocker

import (
	"fmt"
	"io"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/abichinger/odoocker/flags"
	"github.com/leaanthony/clir"
)

type Up struct {
	*Odoocker
	options *flags.Up
}

func NewUp(options *flags.Up) (*Up, error) {
	r := &Up{
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

func (u *Up) cmdOut() io.Writer {
	return os.Stdout
}

func (u *Up) start() error {
	modules := strings.Join(u.options.Modules, ",")
	return u.docker.ExecAs("odoo", web, "/entrypoint.sh", "-d", "odoo", "-i", modules)
}

func (u *Up) Run() error {
	if err := u.Setup(); err != nil {
		return err
	}
	defer u.Teardown()

	go u.start()

	// CREDIT: https://jacobtomlinson.dev/posts/2022/golang-block-until-interrupt-with-ctrl-c/
	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)
	fmt.Println("Blocking, press ctrl+c to continue...")
	<-done // Will block here until user hits ctrl+c

	return nil

}

func AddUpCommand(cli *clir.Cli) {

	cmd := cli.NewSubCommand("up", "Starts the odoo container and installs the addons")
	common := &flags.Common{}
	cmd.AddFlags(common)
	flags := &flags.Up{
		Common: common,
	}
	cmd.AddFlags(flags)

	cmd.Action(func() error {
		flags.Apply()
		action, err := NewUp(flags)
		if err != nil {
			return err
		}
		return action.Run()
	})

}
