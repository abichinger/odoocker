package flags

import (
	"os"
	"path/filepath"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Common command line options
type Common struct {
	Version         string `name:"v" default:"latest" description:"Odoo version"`
	PostgresVersion string `name:"pg" default:"15" description:"PostgreSQL version`
	TempDir         string `default:".odoocker" description:"temporary directory"`
	Addons          string `default:"." description:"path to odoo addons folder or a single addon"`
	Verbose         bool   `description:"verbose output"`
	IsTerminal      bool   `name:"term" default:"true" description:"is terminal"`
	Tour            bool   `name:"tour" default:"false" description:"installs websocket-client and chrome"`
	EnvFile         string `name:"env-file" description:"path to env file"`
}

func (c *Common) Apply() {
	if c.Verbose {
		log.Logger = log.Level(zerolog.DebugLevel)
	}

	if _, err := os.Stat(filepath.Join(c.Addons, "__manifest__.py")); err == nil {
		c.Addons = filepath.Join(c.Addons, "..")
	}

	log.Debug().Str("Odoo", c.Version).Str("PostgreSQL", c.PostgresVersion).Str("Addons", c.Addons).Msg("Options")
}

type Test struct {
	*Common
	TestTags []string `name:"t" description:"odoo test-tags"`
	Modules  []string `name:"m" description:"List of modules to install"`
	// Watch           bool     `name:"w" default:"false" description:"Watch for file changes"`
}

type Run struct {
	*Common
}
