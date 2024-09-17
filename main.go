package main

import (
	"github.com/AzamatAbdranbayev/goplugin/goplugin"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(goplugin.EmptyLineBeforeIfAnalyzer)
}
