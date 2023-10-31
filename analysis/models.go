package analysis

type LinterOut struct {
	Issues []Issues `json:"Issues"`
}
type Pos struct {
	Filename string `json:"Filename"`
	Offset   int    `json:"Offset"`
	Line     int    `json:"Line"`
	Column   int    `json:"Column"`
}
type LineRange struct {
	From int `json:"From"`
	To   int `json:"To"`
}
type Issues struct {
	FromLinter           string    `json:"FromLinter"`
	Text                 string    `json:"Text"`
	Severity             string    `json:"Severity"`
	Replacement          any       `json:"Replacement"`
	Pos                  Pos       `json:"Pos"`
	ExpectNoLint         bool      `json:"ExpectNoLint"`
	ExpectedNoLintLinter string    `json:"ExpectedNoLintLinter"`
	LineRange            LineRange `json:"LineRange,omitempty"`
}
