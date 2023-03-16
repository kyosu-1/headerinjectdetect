package headerinjectdetect

import (
	"go/ast"
	"go/token"
	"go/types"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
	"strings"
)

var Analyzer = &analysis.Analyzer{
	Name: "headerinjection",
	Doc:  "Checks for possible HTTP header injection in Go code",
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
			if containsUserInput(pass, value) {
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

func containsUserInput(pass *analysis.Pass, n ast.Node) bool {
	switch expr := n.(type) {
	case *ast.Ident:
		return isUserControlledVar(pass, expr)
	case *ast.BinaryExpr:
		return expr.Op == token.ADD && (containsUserInput(pass, expr.X) || containsUserInput(pass, expr.Y))
	case *ast.CallExpr:
		return isStringManipulationFunction(pass, expr) &&
			(anyArgumentContainsUserInput(pass, expr.Args))
	}
	return false
}

func isUserControlledVar(pass *analysis.Pass, ident *ast.Ident) bool {
	if obj := pass.TypesInfo.ObjectOf(ident); obj != nil {
		name := obj.Name()
		return strings.HasPrefix(name, "user") || strings.Contains(name, "Input")
	}
	return false
}

func isStringManipulationFunction(pass *analysis.Pass, call *ast.CallExpr) bool {
	if selector, ok := call.Fun.(*ast.SelectorExpr); ok {
		if funcObj, ok := pass.TypesInfo.ObjectOf(selector.Sel).(*types.Func); ok {
			return funcObj.Pkg().Path() == "strings"
		}
	}

	return false
}

func anyArgumentContainsUserInput(pass *analysis.Pass, args []ast.Expr) bool {
	for _, arg := range args {
		if containsUserInput(pass, arg) {
			return true
		}
	}
	return false
}