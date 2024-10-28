package odoocker

import (
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/abichinger/odoocker/flags"
	"github.com/abichinger/odoocker/internal/docker"
)

type IOdoocker interface {
	environ(all bool) ([]string, error)
	cmdOut() io.Writer
	composeFile() string
}

func NewDockerClient(o IOdoocker) (*docker.Client, error) {
	env, err := o.environ(true)
	if err != nil {
		return nil, err
	}

	return &docker.Client{
		Env:         env,
		ComposeFile: o.composeFile(),
		Writer:      o.cmdOut(),
	}, nil
}

type Odoocker struct {
	options *flags.Common
	docker  *docker.Client
}

func NewOdoocker(options *flags.Common) *Odoocker {
	return &Odoocker{
		options: options,
	}
}

func (o *Odoocker) environ(all bool) ([]string, error) {
	addons, _ := filepath.Abs(o.options.Addons)

	env := []string{}
	if all {
		env = os.Environ()
	}
	env = append(env, "ODK_ODOO_TAG="+o.options.Version)
	env = append(env, "ODK_POSTGRES_TAG="+o.options.PostgresVersion)
	env = append(env, "ODK_ADDONS="+addons)
	if o.options.Tour {
		env = append(env, "ODK_TOUR=true")
	}
	if o.options.EnvFile != "" {
		content, err := os.ReadFile(o.options.EnvFile)
		if err != nil {
			return env, err
		}
		env = append(env, strings.Split(string(content), "\n")...)
	}

	return env, nil
}

func (o *Odoocker) WriteEnv() error {
	name := filepath.Join(o.tempDir(), ".env")
	env, err := o.environ(false)
	if err != nil {
		return err
	}
	return os.WriteFile(name, []byte(strings.Join(env, "\n")), os.ModePerm)
}

func (o *Odoocker) composeFile() string {
	return filepath.Join(o.options.TempDir, "docker-compose.yml")
}

func (o *Odoocker) Up() error {
	return o.docker.Up(true, true)
}

func (o *Odoocker) Down() error {
	return o.docker.Down()
}

func (o *Odoocker) createTempDir() error {
	if err := os.MkdirAll(o.tempDir(), os.ModePerm); err != nil {
		return err
	}

	if err := docker.FsCopy("docker-compose.yml", o.composeFile()); err != nil {
		return err
	}

	if err := docker.FsCopy("Dockerfile", filepath.Join(o.options.TempDir, "Dockerfile")); err != nil {
		return err
	}

	return o.WriteEnv()
}

func (o *Odoocker) tempDir() string {
	return o.options.TempDir
}

func (o *Odoocker) Setup() error {
	if err := o.createTempDir(); err != nil {
		return err
	}
	if err := o.Up(); err != nil {
		return err
	}
	return nil
}

func (o *Odoocker) Teardown() error {
	if err := o.Down(); err != nil {
		return err
	}
	return nil
}

func (o *Odoocker) cmdOut() io.Writer {
	panic("unimplemented")
}
