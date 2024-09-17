package goplugin

import (
	"bufio"
	"go/ast"
	"golang.org/x/tools/go/analysis"
	"log"
	"os"
	"strings"
)

var EmptyLineBeforeIfAnalyzer = &analysis.Analyzer{
	Name: "emptylinebeforeif",
	Doc:  "checks aza for an empty line before 'if' statements",
	Run:  runEmptyLineBeforeIfAnalyzer,
}

func runEmptyLineBeforeIfAnalyzer(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			ifStmt, ok := n.(*ast.IfStmt)
			if ok {
				pos := pass.Fset.Position(ifStmt.Pos())
				srcFile, err := os.Open(pos.Filename)
				if err != nil {
					log.Printf("could not read file %s: %v\n", pos.Filename, err)
					return false
				}
				defer srcFile.Close()

				scanner := bufio.NewScanner(srcFile)
				var prevLine string
				var currentLine string
				for i := 1; i <= pos.Line; i++ {
					scanner.Scan()
					if i == pos.Line-1 {
						prevLine = scanner.Text()
					}
					if i == pos.Line {
						currentLine = scanner.Text()
					}
				}

				trimmedCurrentLine := strings.TrimSpace(currentLine)
				if !strings.Contains(prevLine, "//") && strings.TrimSpace(prevLine) != "" && strings.HasPrefix(trimmedCurrentLine, "if") {
					pass.Reportf(ifStmt.Pos(), "missing empty line before 'if' at %s:%d", pos.Filename, pos.Line)
				}
			}
			return true
		})
	}
	return nil, nil
}
