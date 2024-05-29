package docker

import (
	"embed"
	"os"
	"path/filepath"
)

//go:embed docker-compose.yml
var fs embed.FS

// WriteDockerCompose writes the embedded docker-compose.yml into the specified output directory.
func WriteDockerCompose(outputDir string) error {
	data, err := fs.ReadFile("docker-compose.yml")
	if err != nil {
		return err
	}

	outputPath := filepath.Join(outputDir, "docker-compose.yml")
	return os.WriteFile(outputPath, data, os.ModePerm)
}
