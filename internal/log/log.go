package log

import (
	"regexp"
	"strings"
)

type matcher interface {
	match(log string) error
}

type baseMatcher struct {
	substrs []string
	onMatch func(log string) error
}

func NewMatcher(substrs []string, onMatch func(log string) error) *baseMatcher {
	return &baseMatcher{
		substrs: substrs,
		onMatch: onMatch,
	}
}

func (m *baseMatcher) test(log string) bool {
	i := 0
	for _, substr := range m.substrs {
		i = strings.Index(log[i:], substr)
		if i == -1 {
			return false
		}
	}
	return true
}

func (m *baseMatcher) match(log string) error {
	if m.test(log) {
		return m.onMatch(log)
	}
	return nil
}

type analyzer struct {
	current  string
	matchers []matcher
}

func (a *analyzer) processLog() {
	for _, m := range a.matchers {
		m.match(a.current)
	}
	a.current = ""
}

var logExp, _ = regexp.Compile(`^(\S+\s){3}(INFO|WARNING|ERROR)`)

func (a *analyzer) isStart(line string) bool {
	return logExp.MatchString(line)
}

func (a *analyzer) addLine(line string) {
	if len(a.current) > 0 && a.isStart(line) {
		a.processLog()
	}
	a.current += line + "\n"
}

func (a *analyzer) addLines(lines string) {
	for _, line := range strings.Split(lines, "\n") {
		a.addLine(line)
	}
}

func (a *analyzer) Write(p []byte) (n int, err error) {
	a.addLines(strings.TrimSpace(string(p)))
	return len(p), nil
}

func NewAnalyer(matchers ...matcher) *analyzer {
	return &analyzer{
		current:  "",
		matchers: matchers,
	}
}

type OdooAnalyzer struct {
	*analyzer
	Result string
	Errors []string
}

const resultSubStr = "odoo.tests.result:"

func (a *OdooAnalyzer) TestResult() string {
	i := strings.Index(a.Result, resultSubStr)
	if i == -1 {
		return a.Result
	}
	return strings.TrimSpace(a.Result[i+len(resultSubStr):])
}

func (a *OdooAnalyzer) Passed() bool {
	return strings.Contains(a.Result, "0 failed") && len(a.Errors) == 0
}

func (a *OdooAnalyzer) String() string {
	res := a.TestResult() + "\n"
	if len(a.Errors) > 0 {
		res += strings.Join(a.Errors, "\n")
	}
	return res
}

func NewOdooAnalyzer() *OdooAnalyzer {
	a := &OdooAnalyzer{}

	a.analyzer = NewAnalyer(
		NewMatcher([]string{resultSubStr}, func(log string) error {
			a.Result = log
			return nil
		}),
		NewMatcher([]string{" ERROR "}, func(log string) error {
			a.Errors = append(a.Errors, log)
			return nil
		}),
	)

	return a
}
