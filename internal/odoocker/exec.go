package odoocker

import (
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/abichinger/odoocker/internal/term"
	"github.com/rs/zerolog/log"
)

func (o *Odoocker) environ() []string {
	env := os.Environ()
	env = append(env, "ODK_ODOO_TAG="+o.options.Version)
	env = append(env, "ODK_POSTGRES_TAG="+o.options.PostgresVersion)
	env = append(env, "ODK_ADDONS="+o.options.Addons)
	return env
}

func (o *Odoocker) composeFile() (string, error) {
	dir, err := o.tempDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "docker-compose.yml"), nil
}

func (o *Odoocker) dockerRun(args ...string) error {
	args = append([]string{"run", "--no-deps", "--rm", "--no-TTY"}, args...)
	return o.docker(args)
}

func (o *Odoocker) dockerUp(detach bool, service ...string) error {
	args := []string{"up"}
	if detach {
		args = append(args, "-d")
	}
	args = append(args, service...)
	return o.docker(args)
}

func (o *Odoocker) dockerDown(service ...string) error {
	args := []string{"down"}
	args = append(args, service...)
	return o.docker(args)
}

func (o *Odoocker) cmdOut() io.Writer {
	return io.MultiWriter(
		term.NewScroll(os.Stdout, 5),
		o.analyzer,
	)
}

func (o *Odoocker) docker(args []string) error {
	path, err := o.composeFile()
	if err != nil {
		return err
	}
	args = append([]string{"compose", "-f", path}, args...)
	log.Info().Msg("RUN docker " + strings.Join(args, " "))
	cmd := exec.Command("docker", args...)

	cmdOut := o.cmdOut()
	cmd.Stdout = cmdOut
	cmd.Stderr = cmdOut

	cmd.Env = o.environ()
	err = cmd.Run()
	if err != nil {
		log.Error().Err(err)
		return err
	}
	return nil
}
