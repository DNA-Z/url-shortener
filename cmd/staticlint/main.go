package main

import (
	"github.com/timakin/bodyclose/passes/bodyclose"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/assign"
	"golang.org/x/tools/go/analysis/passes/atomic"
	"golang.org/x/tools/go/analysis/passes/bools"
	"golang.org/x/tools/go/analysis/passes/copylock"
	"golang.org/x/tools/go/analysis/passes/httpresponse"
	"golang.org/x/tools/go/analysis/passes/ifaceassert"
	"golang.org/x/tools/go/analysis/passes/loopclosure"
	"golang.org/x/tools/go/analysis/passes/lostcancel"
	"golang.org/x/tools/go/analysis/passes/nilfunc"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/shift"
	"golang.org/x/tools/go/analysis/passes/stringintconv"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"golang.org/x/tools/go/analysis/passes/tests"
	"golang.org/x/tools/go/analysis/passes/timeformat"
	"golang.org/x/tools/go/analysis/passes/unmarshal"
	"golang.org/x/tools/go/analysis/passes/unreachable"
	"honnef.co/go/tools/quickfix"
	"honnef.co/go/tools/simple"
	"honnef.co/go/tools/staticcheck"
	"honnef.co/go/tools/stylecheck"

	"github.com/breml/errchkjson"
)

func main() {
	analyzers := []*analysis.Analyzer{
		assign.Analyzer,
		atomic.Analyzer,
		bools.Analyzer,
		copylock.Analyzer,
		httpresponse.Analyzer,
		ifaceassert.Analyzer,
		loopclosure.Analyzer,
		lostcancel.Analyzer,
		nilfunc.Analyzer,
		printf.Analyzer,
		shadow.Analyzer,
		shift.Analyzer,
		stringintconv.Analyzer,
		structtag.Analyzer,
		tests.Analyzer,
		timeformat.Analyzer,
		unmarshal.Analyzer,
		unreachable.Analyzer,
	}

	// Анализаторы класса SA из staticcheck
	for _, a := range staticcheck.Analyzers {
		analyzers = append(analyzers, a.Analyzer)
	}

	// Анализаторы класса S (simple)
	for _, a := range simple.Analyzers {
		analyzers = append(analyzers, a.Analyzer)
	}

	// Анализаторы класса ST (stylecheck)
	for _, a := range stylecheck.Analyzers {
		analyzers = append(analyzers, a.Analyzer)
	}

	// Анализаторы класса QF (quickfix)
	for _, a := range quickfix.Analyzers {
		analyzers = append(analyzers, a.Analyzer)
	}

	// Дополнительные анализаторы
	analyzers = append(analyzers,
		errchkjson.NewAnalyzer(),
		bodyclose.Analyzer,
	)

	// Запуск multichecker
	multichecker.Main(analyzers...)
}
