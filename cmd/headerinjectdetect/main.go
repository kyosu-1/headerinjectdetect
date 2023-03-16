package main

import (
	"github.com/kyosu-1/headerinjectdetect"
	"golang.org/x/tools/go/analysis/unitchecker"
)

func main() { unitchecker.Main(headerinjectdetect.Analyzer) }
