package headerinjectdetect

import (
	"go/ast"
	"go/token"
	"go/types"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const doc = "headerinjectdetect is ..."

// Analyzer is ...
var Analyzer = &analysis.Analyzer{
	Name: "headerinjectdetect",
	Doc:  doc,
	Run:  run,
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
	},
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspector := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	nodeFilter := []ast.Node{
		(*ast.CallExpr)(nil),
	}

	inspector.Preorder(nodeFilter, func(n ast.Node) {
		call := n.(*ast.CallExpr)

		// Check if the call expression is a selector expression, and proceed accordingly.
		selectorExpr, ok := call.Fun.(*ast.SelectorExpr)
		if !ok {
			return
		}

		funcObj, ok := pass.TypesInfo.ObjectOf(selectorExpr.Sel).(*types.Func)
		if !ok {
			return
		}

		if isHeaderSetMethod(funcObj) {
			value := call.Args[1]
			if isConcatenation(value) {
				pass.Reportf(call.Pos(), "possible HTTP header injection found")
			}
		}
	})

	return nil, nil
}

func isHeaderSetMethod(f *types.Func) bool {
	return f != nil &&
		f.Pkg().Path() == "net/http" &&
		f.Name() == "Set"
}

func isConcatenation(n ast.Node) bool {
	_, ok := n.(*ast.BinaryExpr)
	return ok && n.(*ast.BinaryExpr).Op == token.ADD
}
