package term

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/moby/term"
	"github.com/morikuni/aec"
)

type scroll struct {
	fd     *os.File
	size   int
	lines  []string
	prefix string
}

type scrollOption func(s *scroll) error

func WithPrefix(prefix string) scrollOption {
	return func(s *scroll) error {
		s.prefix = prefix
		return nil
	}
}

func NewScroll(fd *os.File, size int, options ...scrollOption) io.Writer {
	if !term.IsTerminal(fd.Fd()) {
		return fd
	}

	s := &scroll{
		fd:   fd,
		size: size,
	}

	for _, opt := range options {
		opt(s)
	}

	return s
}

func (s *scroll) up(n int) (err error) {
	up := aec.Up(1)
	erase := aec.EraseLine(aec.EraseModes.All)
	for i := 0; i < n; i++ {
		_, err = fmt.Fprint(s.fd, up)
		if err != nil {
			return
		}
		_, err = fmt.Fprint(s.fd, erase)
		if err != nil {
			return
		}
	}
	return
}

func (s *scroll) addLine(line string) {
	if len(s.lines) >= s.size {
		s.lines = append(s.lines[1:], line)
	} else {
		s.lines = append(s.lines, line)
	}
}

func (s *scroll) addLines(lines string) {
	for _, line := range strings.Split(lines, "\n") {
		s.addLine(line)
	}
}

func (s *scroll) Write(p []byte) (n int, err error) {
	s.up(len(s.lines))
	s.addLines(strings.TrimSpace(string(p)))

	for _, line := range s.lines {
		_, err = fmt.Fprintln(s.fd, s.prefix, line)
		if err != nil {
			return
		}
	}
	n = len(p)
	return
}
