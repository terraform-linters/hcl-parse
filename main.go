package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
)

var (
	inFile       = flag.String("f", "", "file to parse")
	exprMode     = flag.String("e", "", "expression to parse")
	templateMode = flag.String("t", "", "template to parse")
)

func main() {
	flag.Usage = usage
	flag.Parse()

	var node hclsyntax.Node
	var bytes []byte

	// TODO:
	// - add support file argument
	// - add support stdin
	// - improve diagnostics writer

	if inFile != nil && *inFile != "" {
		f, err := os.ReadFile(*inFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error reading file: %s\n", err)
			os.Exit(1)
		}
		file, diags := hclsyntax.ParseConfig(f, *inFile, hcl.InitialPos)
		if diags.HasErrors() {
			fmt.Fprintf(os.Stderr, "error parsing file: %s\n", diags)
			os.Exit(1)
		}
		node = file.Body.(*hclsyntax.Body)
		bytes = file.Bytes
	} else if exprMode != nil && *exprMode != "" {
		expr, diags := hclsyntax.ParseExpression([]byte(*exprMode), "<expr>", hcl.InitialPos)
		if diags.HasErrors() {
			fmt.Fprintf(os.Stderr, "error parsing expression: %s\n", diags)
			os.Exit(1)
		}
		node = expr
		bytes = []byte(*exprMode)
	} else if templateMode != nil && *templateMode != "" {
		expr, diags := hclsyntax.ParseTemplate([]byte(*templateMode), "<template>", hcl.InitialPos)
		if diags.HasErrors() {
			fmt.Fprintf(os.Stderr, "error parsing template: %s\n", diags)
			os.Exit(1)
		}
		node = expr
		bytes = []byte(*templateMode)
	} else {
		usage()
	}

	hclsyntax.Walk(node, &walker{file: bytes})
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: hclparse [options]\n")
	flag.PrintDefaults()
	os.Exit(2)
}

type walker struct {
	indent int
	leaf   bool
	file   []byte
}

var _ hclsyntax.Walker = (*walker)(nil)

func (w *walker) Enter(node hclsyntax.Node) hcl.Diagnostics {
	if w.leaf {
		panic("leaf node should not have children")
	}

	fmt.Print(strings.Repeat(" ", w.indent))

	switch node := node.(type) {
	case *hclsyntax.Attribute:
		fmt.Printf(`(%T "%s"`, node, node.Name)
	case *hclsyntax.Block:
		fmt.Printf(`(%T "%s" %s`, node, node.Type, node.Labels)
	case *hclsyntax.LiteralValueExpr:
		fmt.Printf(`(%T "%s")`, node, node.SrcRange.SliceBytes(w.file))
		w.leaf = true
	case *hclsyntax.ScopeTraversalExpr:
		fmt.Printf(`(%T "%s")`, node, node.SrcRange.SliceBytes(w.file))
		w.leaf = true
	case *hclsyntax.RelativeTraversalExpr:
		fmt.Printf(`(%T "%s"`, node, node.Traversal.SourceRange().SliceBytes(w.file))
	case *hclsyntax.FunctionCallExpr:
		fmt.Printf(`(%T "%s"`, node, node.Name)
	case *hclsyntax.ForExpr:
		fmt.Printf(`(%T`, node)
		if node.KeyVar != "" {
			fmt.Printf(` key="%s"`, node.KeyVar)
		}
		if node.ValVar != "" {
			fmt.Printf(` val="%s"`, node.ValVar)
		}
	case *hclsyntax.AnonSymbolExpr:
		fmt.Printf(`(%T)`, node)
		w.leaf = true
	case *hclsyntax.BinaryOpExpr:
		fmt.Printf(`(%T "%s"`, node, opAsString(node.Op))
	case *hclsyntax.UnaryOpExpr:
		fmt.Printf(`(%T "%s"`, node, opAsString(node.Op))
	default:
		fmt.Printf("(%T", node)
	}

	fmt.Print("\n")
	w.indent += 2
	return nil
}

func (w *walker) Exit(node hclsyntax.Node) hcl.Diagnostics {
	w.indent -= 2

	if w.leaf {
		w.leaf = false
		return nil
	}

	fmt.Print(strings.Repeat(" ", w.indent))
	fmt.Printf(")\n")
	return nil
}

func opAsString(op *hclsyntax.Operation) string {
	switch op {
	case hclsyntax.OpLogicalOr:
		return "||"
	case hclsyntax.OpLogicalAnd:
		return "&&"
	case hclsyntax.OpLogicalNot:
		return "!"
	case hclsyntax.OpEqual:
		return "=="
	case hclsyntax.OpNotEqual:
		return "!="
	case hclsyntax.OpGreaterThan:
		return ">"
	case hclsyntax.OpGreaterThanOrEqual:
		return ">="
	case hclsyntax.OpLessThan:
		return "<"
	case hclsyntax.OpLessThanOrEqual:
		return "<="
	case hclsyntax.OpAdd:
		return "+"
	case hclsyntax.OpSubtract:
		return "-"
	case hclsyntax.OpMultiply:
		return "*"
	case hclsyntax.OpDivide:
		return "/"
	case hclsyntax.OpModulo:
		return "%"
	case hclsyntax.OpNegate:
		return "-"
	default:
		panic(fmt.Sprintf("unknown operation type: %T", op))
	}
}
