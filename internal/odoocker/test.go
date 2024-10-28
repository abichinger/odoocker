package odoocker

import (
	"io"
	"os"
	"strings"

	"github.com/abichinger/odoocker/flags"
	olog "github.com/abichinger/odoocker/internal/log"
	"github.com/abichinger/odoocker/internal/term"
	"github.com/leaanthony/clir"
	"github.com/rs/zerolog/log"
)

const db = "db"
const web = "web"

type Test struct {
	*Odoocker

	options  *flags.Test
	analyzer *olog.OdooAnalyzer
}

func NewTest(options *flags.Test) (*Test, error) {
	t := &Test{
		Odoocker: NewOdoocker(options.Common),
		options:  options,
		analyzer: olog.NewOdooAnalyzer(),
	}

	docker, err := NewDockerClient(t)
	if err != nil {
		return t, err
	}
	t.docker = docker
	return t, nil
}

func (t *Test) cmdOut() io.Writer {
	var out io.Writer = os.Stdout
	if t.options.IsTerminal && term.IsTerminal(os.Stdout.Fd()) {
		out = term.NewScroll(os.Stdout, 5)
	}

	return io.MultiWriter(
		out,
		t.analyzer,
	)
}

func (t *Test) Run() error {
	if err := t.Setup(); err != nil {
		return err
	}
	defer t.Teardown()

	if err := t.install(); err != nil {
		return err
	}

	if err := t.test(); err != nil {
		return err
	}

	result := t.analyzer.String()
	if t.analyzer.Passed() {
		log.Info().Msg(result)
	} else {
		log.Error().Msg(result)
		os.Exit(1)
	}

	return nil
}

func (t *Test) install() error {
	modules := strings.Join(t.options.Modules, ",")
	return t.docker.ExecAs("odoo", web, "/entrypoint.sh", "--stop-after-init", "-d", "odoo", "-i", modules)
}

func (t *Test) test() error {
	testTags := strings.Join(t.options.TestTags, ",")
	return t.docker.ExecAs("odoo", web, "/entrypoint.sh", "--stop-after-init", "-d", "odoo", "--test-tags", testTags)
}

func AddTestCommand(cli *clir.Cli) {

	cmd := cli.NewSubCommand("test", "Run odoo tests")
	common := &flags.Common{}
	cmd.AddFlags(common)
	flags := &flags.Test{
		Common: common,
	}
	cmd.AddFlags(flags)

	cmd.Action(func() error {

		flags.Apply()

		if len(flags.Modules) == 0 {
			log.Error().Msg("Provide at least one module")
			cmd.PrintHelp()
			os.Exit(1)
		}

		if len(flags.TestTags) == 0 {
			log.Error().Msg("Provide at least one test tag")
			cmd.PrintHelp()
			os.Exit(1)
		}

		action, err := NewTest(flags)
		if err != nil {
			return err
		}
		return action.Run()
	})

}
