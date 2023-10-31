package analysis

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

func Run(reportFile string) error {
	curRep := LinterOut{}

	f, err := readFile(reportFile)
	if err != nil {
		return fmt.Errorf("read initial report file: %w", err)
	}

	if err := json.Unmarshal(f, &curRep); err != nil {
		return fmt.Errorf("unmarshal report file: %w", err)
	}

	cur := report{}
	cur.initFree2GoWith()

	for i := range curRep.Issues {
		issue := curRep.Issues[i]

		cur.linterProblemFaced(linterName(issue.FromLinter))
		cur.addProblem(linterName(issue.FromLinter), issue.Pos.Filename)
	}

	cur.countUniqueFiles()

	cur.OutProblems = make([]LinterProblem, 0, len(cur.Problems))
	for i := range cur.Problems {
		problem := cur.Problems[i]

		cur.OutProblems = append(cur.OutProblems, problem)
	}

	cur.FreeToGoWith = make([]linterName, 0, len(cur.Free2GoWith))
	for linter := range cur.Free2GoWith {
		cur.FreeToGoWith = append(cur.FreeToGoWith, linter)
	}

	if err := marshalBack(cur); err != nil {
		return fmt.Errorf("marshal output: %w", err)
	}

	return nil
}

func readFile(reportName string) ([]byte, error) {
	f, err := os.ReadFile(reportName)
	if err != nil {
		return nil, fmt.Errorf("read report file: %w", err)
	}

	return f, nil
}

type linterName string

type report struct {
	Free2GoWith  map[linterName]struct{}      `json:"-"`
	Problems     map[linterName]LinterProblem `json:"-"`
	FreeToGoWith []linterName                 `json:"free_to_go_with"`
	OutProblems  []LinterProblem              `json:"problems"`
}

type LinterProblem struct {
	LinterName         linterName
	Count              int
	UniqueFiles        int
	UniqueFilesWOTests int
	Files              []string `json:"-"`
}

func (r *report) initFree2GoWith() {
	r.Problems = make(map[linterName]LinterProblem)
	r.Free2GoWith = make(map[linterName]struct{})

	for linter := range lintersToAnalize {
		r.Free2GoWith[linterName(linter)] = struct{}{}
	}
}

func (r *report) linterProblemFaced(linter linterName) {
	delete(r.Free2GoWith, linter)
}

func (r *report) addProblem(linter linterName, filename string) {
	if _, ok := r.Problems[linter]; !ok {
		r.Problems[linter] = LinterProblem{
			LinterName: linter,
			Count:      1,
			Files:      []string{filename},
		}

		return
	}

	linterState := r.Problems[linter]
	linterState.Count++
	linterState.Files = append(linterState.Files, filename)
	r.Problems[linter] = linterState
}

func (r *report) countUniqueFiles() {
	for linter, linterState := range r.Problems {
		testFiles := 0
		uniqueFiles := make(map[string]struct{})

		for _, file := range linterState.Files {
			uniqueFiles[file] = struct{}{}
		}

		for file := range uniqueFiles {
			if strings.Contains(file, "_test.go") {
				testFiles++
			}
		}

		linterState.UniqueFiles = len(uniqueFiles)
		linterState.UniqueFilesWOTests = len(uniqueFiles) - testFiles
		r.Problems[linter] = linterState
	}
}

func marshalBack(tf report) error {
	out, err := json.MarshalIndent(&tf, "", "\t")
	if err != nil {
		return fmt.Errorf("marshal results: %w", err)
	}

	fmt.Println(string(out))

	return nil
}

var lintersToAnalize = map[string]struct{}{
	"asasalint":                 {},
	"asciicheck":                {},
	"bidichk":                   {},
	"bodyclose":                 {},
	"containedctx":              {},
	"contextcheck":              {},
	"dogsled":                   {},
	"dupl":                      {},
	"dupword":                   {},
	"durationcheck":             {},
	"errcheck":                  {},
	"errchkjson":                {},
	"errname":                   {},
	"errorlint":                 {},
	"exhaustive":                {},
	"exportloopref":             {},
	"forcetypeassert":           {},
	"funlen":                    {},
	"gci":                       {},
	"gocheckcompilerdirectives": {},
	"gochecknoinits":            {},
	"gocognit":                  {},
	"goconst":                   {},
	"gocritic":                  {},
	"gocyclo":                   {},
	"godox":                     {},
	"goerr113":                  {},
	"gofmt":                     {},
	"gofumpt":                   {},
	"goimports":                 {},
	"gomnd":                     {},
	"goprintffuncname":          {},
	"gosec":                     {},
	"gosimple":                  {},
	"govet":                     {},
	"ineffassign":               {},
	"interfacebloat":            {},
	"lll":                       {},
	"makezero":                  {},
	"misspell":                  {},
	"musttag":                   {},
	"nakedret":                  {},
	"nestif":                    {},
	"nilerr":                    {},
	"nilnil":                    {},
	"noctx":                     {},
	"nolintlint":                {},
	"nonamedreturns":            {},
	"paralleltest":              {},
	"prealloc":                  {},
	"predeclared":               {},
	"reassign":                  {},
	"revive":                    {},
	"rowserrcheck":              {},
	"sqlclosecheck":             {},
	"staticcheck":               {},
	"stylecheck":                {},
	"tenv":                      {},
	"thelper":                   {},
	"tparallel":                 {},
	"typecheck":                 {},
	"unconvert":                 {},
	"unparam":                   {},
	"unused":                    {},
	"usestdlibvars":             {},
	"varnamelen":                {},
	"wastedassign":              {},
	"whitespace":                {},
	"wrapcheck":                 {},
	"wsl":                       {},
}
