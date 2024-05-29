package odoocker

import (
	"os"
	"strings"

	"github.com/abichinger/odoocker/flags"
	olog "github.com/abichinger/odoocker/internal/log"
	"github.com/rs/zerolog/log"
)

const db = "db"
const web = "web"

type Odoocker struct {
	options  *flags.Options
	analyzer *olog.OdooAnalyzer
}

func NewOdoocker(options *flags.Options) *Odoocker {
	return &Odoocker{
		options:  options,
		analyzer: olog.NewOdooAnalyzer(),
	}
}

func (o *Odoocker) Run() error {
	err := o.setup()
	if err != nil {
		return err
	}

	err = o.startDb()
	if err != nil {
		return err
	}
	defer o.stopDb()

	err = o.install()
	if err != nil {
		return err
	}

	err = o.test()
	if err != nil {
		return err
	}

	result := o.analyzer.String()
	if o.analyzer.Passed() {
		log.Info().Msg(result)
	} else {
		log.Error().Msg(result)
		os.Exit(1)
	}

	return nil
}

func (o *Odoocker) startDb() error {
	return o.dockerUp(true, db)
}

func (o *Odoocker) stopDb() error {
	return o.dockerDown(db)
}

func (o *Odoocker) install() error {
	modules := strings.Join(o.options.Modules, ",")
	return o.dockerRun(web, "odoo", "--stop-after-init", "-d", "odoo", "-i", modules)
}

func (o *Odoocker) test() error {
	testTags := strings.Join(o.options.TestTags, ",")
	return o.dockerRun(web, "odoo", "--stop-after-init", "-d", "odoo", "--test-tags", testTags)
}
