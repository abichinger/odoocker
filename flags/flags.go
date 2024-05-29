package flags

// Options struct to hold all command line options
type Options struct {
	Version         string   `name:"v" default:"latest" description:"Odoo version"`
	PostgresVersion string   `name:"pg" default:"15" description:"PostgreSQL version`
	TestTags        []string `name:"t" description:"odoo test-tags"`
	Modules         []string `name:"m" description:"List of modules to install"`
	TempDir         string   `default:".odoocker" description:"temporary directory"`
	Addons          string   `default:"./addons" description:"path to odoo addons"`
	// Watch           bool     `name:"w" default:"false" description:"Watch for file changes"`
	// Verbose         bool     `name:"verbose"`
}
