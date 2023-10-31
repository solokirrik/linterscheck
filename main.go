package main

import (
	"flag"
	"fmt"

	"github.com/solokirrik/linterscheck/analysis"
	"github.com/solokirrik/linterscheck/results"
)

func main() {
	help := flag.Bool("help", false, "selected report file")
	preAnalyzeFilePath := flag.String("pre-analise-file", "", "selected report file")
	postAnalyzeFilePath := flag.String("post-analise-file", "", "selected report file")

	flag.Parse()

	if help != nil && *help {
		fmt.Println(`Usage:
	--pre-analise-file=<path to report file>
	--post-analise-file=<path to report file>`)

		return
	}

	if postAnalyzeFilePath != nil && preAnalyzeFilePath != nil &&
		*postAnalyzeFilePath != "" && *preAnalyzeFilePath != "" {
		if err := results.Run(*preAnalyzeFilePath, *postAnalyzeFilePath); err != nil {
			fmt.Println(err)
		}

		return
	}

	if preAnalyzeFilePath == nil && postAnalyzeFilePath == nil {
		fmt.Println("report file path is required")
		return
	}

	if preAnalyzeFilePath != nil && *preAnalyzeFilePath != "" {
		if err := analysis.Run(*preAnalyzeFilePath); err != nil {
			fmt.Println(err)
		}

		return
	}
}
