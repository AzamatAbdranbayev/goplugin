package main

import (
	"bufio"
	"fmt"
	"go/ast"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/singlechecker"
	"log"
	"os"
	"strings"
)

func main() {

	singlechecker.Main(EmptyLineBeforeIfAnalyzer)
}

var EmptyLineBeforeIfAnalyzer = &analysis.Analyzer{
	Name: "emptylinebeforeif",
	Doc:  "Checks for an empty line before 'if' statements",
	Run:  runEmptyLineBeforeIfAnalyzer,
}

func runEmptyLineBeforeIfAnalyzer(pass *analysis.Pass) (interface{}, error) {
	log.Printf("=========")
	fmt.Println("++++++++++")
	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			ifStmt, ok := n.(*ast.IfStmt)
			if ok {
				// Get the position of the "if" statement
				pos := pass.Fset.Position(ifStmt.Pos())

				// Open the source file to read the lines
				srcFile, err := os.Open(pos.Filename)
				if err != nil {
					pass.Reportf(ifStmt.Pos(), "could not read file %s: %v", pos.Filename, err)
					return false
				}
				defer srcFile.Close()

				// Read the source file line by line
				scanner := bufio.NewScanner(srcFile)
				var prevLine, currentLine string
				for i := 1; i <= pos.Line; i++ {
					scanner.Scan()
					if i == pos.Line-1 {
						prevLine = scanner.Text() // Previous line
					}
					if i == pos.Line {
						currentLine = scanner.Text() // Current line with "if"
					}
				}

				// Trim spaces and check if the previous line is not a comment or empty
				trimmedCurrentLine := strings.TrimSpace(currentLine)
				log.Printf("=========")
				log.Printf(prevLine)
				fmt.Printf("%d", "hello")
				if !strings.Contains(prevLine, "//") && strings.TrimSpace(prevLine) != "" && strings.HasPrefix(trimmedCurrentLine, "if") {
					pass.Reportf(ifStmt.Pos(), "missing empty line before 'if' at %s:%d", pos.Filename, pos.Line)
				}
			}
			return true
		})
	}
	return nil, nil
}
