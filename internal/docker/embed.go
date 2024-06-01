package docker

import (
	"embed"
	"os"
)

//go:embed docker-compose.yml
//go:embed Dockerfile
var fs embed.FS

// FsCopy writes the embedded file to dest
func FsCopy(src string, dest string) error {
	data, err := fs.ReadFile(src)
	if err != nil {
		return err
	}

	return os.WriteFile(dest, data, os.ModePerm)
}
