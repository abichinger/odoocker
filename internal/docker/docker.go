package docker

import (
	"io"
	"os/exec"
	"strings"

	"github.com/rs/zerolog/log"
)

type Client struct {
	ComposeFile string
	Writer      io.Writer
	Env         []string
}

// func (c *Client) Run(args ...string) error {
// 	args = append([]string{"run", "--no-deps", "--rm", "--no-TTY"}, args...)
// 	return c.Compose(args)
// }

func (c *Client) Exec(args ...string) error {
	args = append([]string{"exec", "--no-TTY"}, args...)
	return c.Compose(args)
}

func (c *Client) Up(detach bool, build bool, service ...string) error {
	args := []string{"up"}
	if detach {
		args = append(args, "-d")
	}
	if build {
		args = append(args, "--build")
	}
	args = append(args, service...)
	return c.Compose(args)
}

func (c *Client) Down(service ...string) error {
	args := []string{"down"}
	args = append(args, service...)
	return c.Compose(args)
}

func (c *Client) Compose(args []string) error {
	args = append([]string{"compose", "-f", c.ComposeFile}, args...)
	log.Info().Msg("RUN docker " + strings.Join(args, " "))
	cmd := exec.Command("docker", args...)
	cmd.Stdout = c.Writer
	cmd.Stderr = c.Writer

	cmd.Env = c.Env
	err := cmd.Run()
	if err != nil {
		log.Error().Err(err)
		return err
	}
	return nil
}
