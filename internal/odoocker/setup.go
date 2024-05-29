package odoocker

import (
	"os"
	"path/filepath"

	"github.com/abichinger/odoocker/internal/docker"
)

func (o *Odoocker) setup() error {
	if err := o.setupTempDir(); err != nil {
		return err
	}
	tmpDir, err := o.tempDir()
	if err != nil {
		return err
	}
	return docker.WriteDockerCompose(tmpDir)
}

func (o *Odoocker) tempDir() (string, error) {
	return filepath.Abs(o.options.TempDir)
}

func (o *Odoocker) setupTempDir() error {
	path, err := o.tempDir()
	if err != nil {
		return err
	}

	return os.MkdirAll(path, os.ModePerm)
}
