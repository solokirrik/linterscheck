# linterscheck

Dead-simple (basically script) to:

1. Aggregate current linters issues
    ```bash
    golangci-lint run --print-issued-lines=false --out-format=json > old-report.json
    go run ./main.go --pre-analise-file=old-report.json 
    ```
2. Compare two golangci-lint reports by categories (old vs new state)
    ```bash
    golangci-lint run --print-issued-lines=false --out-format=json > new-report.json
    go run ./main.go --pre-analise-file=old-report.json --post-analise-file=new-report.json
    ```
