// This must be package main
package main

import (
	"bufio"
	"fmt"
	"go/ast"
	"log"
	"os"
	"strings"

	"golang.org/x/tools/go/analysis"
)

func New(conf any) ([]*analysis.Analyzer, error) {
	// TODO: This must be implemented
	fmt.Println("aza9")
	fmt.Printf("My configuration (%[1]T): %#[1]v\n", conf)
	fmt.Println("************************")
	// Вывод настроек, переданных из settings
	fmt.Printf("Настройки (%[1]T): %#[1]v\n", conf)

	// Пример: вы можете обрабатывать настройки
	if configMap, ok := conf.(map[string]interface{}); ok {
		if value, exists := configMap["one"]; exists {
			fmt.Println("Настройка 'one':", value)
		}
	}
	// The configuration type will be map[string]any or []interface, it depends on your configuration.
	// You can use https://github.com/go-viper/mapstructure to convert map to struct.

	return []*analysis.Analyzer{EmptyLineBeforeIfAnalyzer}, nil
}

// Analyzer to check for empty lines before "if" statements
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
				os.Stderr.Write([]byte("==== Logging from plugin ====\n"))
				// Trim spaces and check if the previous line is not a comment or empty
				trimmedCurrentLine := strings.TrimSpace(currentLine)

				fmt.Fprintf(os.Stderr, "%s: integer addition found: ------------------", prevLine)

				fmt.Println("++++++++++")
				os.Stdout.Write([]byte("++++++++++++++"))
				pass.Reportf(pass.Files[0].Pos(), "Debug: анализатор начал работу")

				if !strings.Contains(prevLine, "//") && strings.TrimSpace(prevLine) != "" && strings.HasPrefix(trimmedCurrentLine, "if") {
					pass.Reportf(ifStmt.Pos(), "missing empty line before 'if' at %s:%d", pos.Filename, pos.Line)
					fmt.Println("==]]]]]]]]]]]]]")
				}
			}
			return true
		})
	}
	return nil, nil
}
