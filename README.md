# hcl-parse

`hcl-parse` is a command line tool that dumps HCL native syntax AST.
This is a helper that prints the nodes processed by `hclsyntax.Walk`.

## Installation

If you have a working Go development environment, you can install this tool
with `go install` in the usual way:

```
$ go install github.com/terraform-linters/hcl-parse@latest
```

This will install `hcl-parse` in `$GOPATH/bin`, which usually places it into
your shell `PATH` so you can then run it as `hcl-parse`.

## Usage

```
usage: hcl-parse [options] [file or content]
  -e    parse as expression
  -t    parse as template
  -version
        show the version number and immediately exit
```

## Examples

Parse a configuration file:

```hcl
cond = var.enabled ? (true) : func(1, var.input)[0]
foo "var" "baz" {
  for = [for x in var.foo: x + 1 if x < 10]
  obj = { a = var.bar[*], var.foo = var.baz[var.qux], c = [1, 2] } 
  temp = "%{ for v in [true] }${v}%{ endfor }"
  wrap = "${true}"
}
```

```console
$ hcl-parse main.hcl
(*hclsyntax.Body
  (hclsyntax.Attributes
    (*hclsyntax.Attribute "cond"
      (*hclsyntax.ConditionalExpr
        (*hclsyntax.ScopeTraversalExpr "var.enabled")
        (*hclsyntax.ParenthesesExpr
          (*hclsyntax.LiteralValueExpr "true")
        )
        (*hclsyntax.RelativeTraversalExpr "[0]"
          (*hclsyntax.FunctionCallExpr "func"
            (*hclsyntax.LiteralValueExpr "1")
            (*hclsyntax.ScopeTraversalExpr "var.input")
          )
        )
      )
    )
  )
  (hclsyntax.Blocks
    (*hclsyntax.Block "foo" [var baz]
      (*hclsyntax.Body
        (hclsyntax.Attributes
          (*hclsyntax.Attribute "for"
            (*hclsyntax.ForExpr val="x"
              (*hclsyntax.ScopeTraversalExpr "var.foo")
              (hclsyntax.ChildScope
                (*hclsyntax.BinaryOpExpr "+"
                  (*hclsyntax.ScopeTraversalExpr "x")
                  (*hclsyntax.LiteralValueExpr "1")
                )
              )
              (hclsyntax.ChildScope
                (*hclsyntax.BinaryOpExpr "<"
                  (*hclsyntax.ScopeTraversalExpr "x")
                  (*hclsyntax.LiteralValueExpr "10")
                )
              )
            )
          )
          (*hclsyntax.Attribute "obj"
            (*hclsyntax.ObjectConsExpr
              (*hclsyntax.ObjectConsKeyExpr keyword="a"
              )
              (*hclsyntax.SplatExpr
                (*hclsyntax.ScopeTraversalExpr "var.bar")
                (*hclsyntax.AnonSymbolExpr)
              )
              (*hclsyntax.ObjectConsKeyExpr
                (*hclsyntax.ScopeTraversalExpr "var.foo")
              )
              (*hclsyntax.IndexExpr
                (*hclsyntax.ScopeTraversalExpr "var.baz")
                (*hclsyntax.ScopeTraversalExpr "var.qux")
              )
              (*hclsyntax.ObjectConsKeyExpr keyword="c"
              )
              (*hclsyntax.TupleConsExpr
                (*hclsyntax.LiteralValueExpr "1")
                (*hclsyntax.LiteralValueExpr "2")
              )
            )
          )
          (*hclsyntax.Attribute "temp"
            (*hclsyntax.TemplateExpr
              (*hclsyntax.TemplateJoinExpr
                (*hclsyntax.ForExpr val="v"
                  (*hclsyntax.TupleConsExpr
                    (*hclsyntax.LiteralValueExpr "true")
                  )
                  (hclsyntax.ChildScope
                    (*hclsyntax.TemplateExpr
                      (*hclsyntax.ScopeTraversalExpr "v")
                    )
                  )
                )
              )
            )
          )
          (*hclsyntax.Attribute "wrap"
            (*hclsyntax.TemplateWrapExpr
              (*hclsyntax.LiteralValueExpr "true")
            )
          )
        )
        (hclsyntax.Blocks
        )
      )
    )
  )
)
```

Parse input as an expression:

```console
$ hcl-parse -e "var.enabled ? var.foo : var.bar"
(*hclsyntax.ConditionalExpr
  (*hclsyntax.ScopeTraversalExpr "var.enabled")
  (*hclsyntax.ScopeTraversalExpr "var.foo")
  (*hclsyntax.ScopeTraversalExpr "var.bar")
)
```

Parse input as a template:

```console
$ hcl-parse -t "%{ for v in [true] }${v}%{ endfor }"
(*hclsyntax.TemplateExpr
  (*hclsyntax.TemplateJoinExpr
    (*hclsyntax.ForExpr val="v"
      (*hclsyntax.TupleConsExpr
        (*hclsyntax.LiteralValueExpr "true")
      )
      (hclsyntax.ChildScope
        (*hclsyntax.TemplateExpr
          (*hclsyntax.LiteralValueExpr "")
        )
      )
    )
  )
)
```
