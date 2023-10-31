package results

import (
	"encoding/json"
	"fmt"
	"os"
)

func Run(oldStateReport, newStateReport string) error {
	fileOld, err := readFile(oldStateReport)
	if err != nil {
		return fmt.Errorf("read initial report file: %w", err)
	}

	fileNew, err := readFile(newStateReport)
	if err != nil {
		return fmt.Errorf("read results report file: %w", err)
	}

	oldResult, err := processResult(fileOld)
	if err != nil {
		return fmt.Errorf("process old report: %w", err)
	}

	oldState, err := json.MarshalIndent(oldResult, "", "\t")
	if err != nil {
		return fmt.Errorf("marshal old state out: %w", err)
	}

	fmt.Println(string(oldState))

	newResult, err := processResult(fileNew)
	if err != nil {
		return fmt.Errorf("process new report: %w", err)
	}

	newState, err := json.MarshalIndent(newResult, "", "\t")
	if err != nil {
		return fmt.Errorf("marshal new state out: %w", err)
	}

	fmt.Println(string(newState))

	return nil
}

func readFile(reportName string) ([]byte, error) {
	f, err := os.ReadFile(reportName)
	if err != nil {
		return nil, fmt.Errorf("read report file: %w", err)
	}

	return f, nil
}

func processResult(inp []byte) (map[string]int, error) {
	report := LinterOut{}
	if err := json.Unmarshal(inp, &report); err != nil {
		return nil, fmt.Errorf("unmarshal report: %w", err)
	}

	out := make(map[string]int)
	for i := range categories {
		out[categories[i]] = 0
	}

	for i := range report.Issues {
		issue := report.Issues[i]

		lintCats, gotCats := categoriesIndex[issue.FromLinter]
		if !gotCats {
			fmt.Println("unknown linter:", issue.FromLinter)
			continue
		}

		for i := range categories {
			if lintCats[i] == 0 {
				continue
			}

			cat := categories[i]

			existingCat, gotExistingCat := out[cat]
			if !gotExistingCat {
				out[cat] = 1
				continue
			}

			out[cat] = existingCat + 1
		}
	}

	return out, nil
}

var categories = []string{
	"bugs", "error", "metalinter", "style", "module", "format", "comment", "performance", "unused", "complexity",
	"test", "import", "sql",
}

var categoriesIndex = map[string][]int{
	"asasalint":                 {1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	"asciicheck":                {1, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	"bidichk":                   {1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	"bodyclose":                 {1, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0},
	"containedctx":              {0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	"contextcheck":              {1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	"cyclop":                    {0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	"decorder":                  {0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	"depguard":                  {0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	"dogsled":                   {0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	"dupl":                      {0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	"dupword":                   {0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0},
	"durationcheck":             {1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	"errcheck":                  {1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	"errchkjson":                {1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	"errname":                   {0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	"errorlint":                 {1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	"execinquery":               {0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	"exhaustive":                {1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	"exhaustruct":               {0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	"exportloopref":             {1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	"forbidigo":                 {0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	"forcetypeassert":           {0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	"funlen":                    {0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0},
	"gci":                       {0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 1, 0},
	"gocheckcompilerdirectives": {1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	"gochecknoinits":            {0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	"gocognit":                  {0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0},
	"goconst":                   {0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	"gocritic":                  {0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	"gocyclo":                   {0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0},
	"godot":                     {0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	"godox":                     {0, 0, 0, 1, 0, 0, 1, 0, 0, 0, 0, 0, 0},
	"goerr113":                  {0, 1, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	"gofmt":                     {0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0},
	"gofumpt":                   {0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0},
	"goheader":                  {0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	"goimports":                 {0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 1, 0},
	"gomnd":                     {0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	"gomoddirectives":           {0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	"gomodguard":                {0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	"goprintffuncname":          {0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	"gosec":                     {1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	"gosimple":                  {0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	"gosmopolitan":              {0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	"govet":                     {1, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	"grouper":                   {0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	"importas":                  {0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	"ineffassign":               {0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0},
	"interfacebloat":            {0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	"ireturn":                   {0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	"lll":                       {0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	"loggercheck":               {0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	"maintidx":                  {0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	"makezero":                  {1, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	"misspell":                  {0, 0, 0, 1, 0, 0, 1, 0, 0, 0, 0, 0, 0},
	"musttag":                   {1, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	"nakedret":                  {0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	"nestif":                    {0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0},
	"nilerr":                    {1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	"nilnil":                    {0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	"nlreturn":                  {0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	"noctx":                     {1, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0},
	"nolintlint":                {0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	"nonamedreturns":            {0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	"paralleltest":              {0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 1, 0, 0},
	"prealloc":                  {0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0},
	"predeclared":               {0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	"promlinter":                {0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	"reassign":                  {1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	"revive":                    {0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	"rowserrcheck":              {1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	"sqlclosecheck":             {1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	"staticcheck":               {1, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	"stylecheck":                {0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	"tagalign":                  {0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	"tagliatelle":               {0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	"tenv":                      {0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	"testableexamples":          {0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	"testpackage":               {0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	"thelper":                   {0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	"tparallel":                 {0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 1, 0, 0},
	"typecheck":                 {1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	"unconvert":                 {0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	"unparam":                   {0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0},
	"unused":                    {0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0},
	"usestdlibvars":             {0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	"varnamelen":                {0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	"wastedassign":              {0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	"whitespace":                {0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	"wrapcheck":                 {0, 1, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	"wsl":                       {0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0},
}
