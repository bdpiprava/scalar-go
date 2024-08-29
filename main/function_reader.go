package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"os"
	"reflect"
	"runtime"
)

func readFuncBodyIgnoreError(fn reflect.Value) string {
	body, _ := readFuncBody(fn)
	return fmt.Sprintf(`func example() (string, error)%s}`, body)
}

func readFuncBody(fn reflect.Value) (string, error) {
	p := fn.Pointer()
	fc := runtime.FuncForPC(p)
	filename, line := fc.FileLine(p)
	fset := token.NewFileSet()
	// parse file to AST tree
	node, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		return "", err
	}
	// walk and find the function block
	find := &FindBlockByLine{Fset: fset, Line: line}
	ast.Walk(find, node)

	if find.Block != nil {
		fp, err := os.Open(filename)
		if err != nil {
			return "", err
		}
		defer fp.Close()
		fp.Seek(int64(find.Block.Lbrace-1), 0)
		buf := make([]byte, int64(find.Block.Rbrace-find.Block.Lbrace))
		_, err = io.ReadFull(fp, buf)
		if err != nil {
			return "", err
		}

		return string(buf), nil
	}
	return "", nil
}

// FindBlockByLine is a ast.Visitor implementation that finds a block by line.
type FindBlockByLine struct {
	Fset  *token.FileSet
	Line  int
	Block *ast.BlockStmt
}

// Visit implements the ast.Visitor interface.
func (f *FindBlockByLine) Visit(node ast.Node) ast.Visitor {
	if node == nil {
		return nil
	}

	if blockStmt, ok := node.(*ast.BlockStmt); ok {
		stmtStartingPosition := blockStmt.Pos()
		stmtLine := f.Fset.Position(stmtStartingPosition).Line
		if stmtLine == f.Line {
			f.Block = blockStmt
			return nil
		}
	}
	return f
}
