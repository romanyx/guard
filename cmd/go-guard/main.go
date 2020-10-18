package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"go/ast"
	"go/format"
	"go/types"
	"html/template"
	"os"
	"strings"

	"golang.org/x/tools/go/packages"
)

const (
	usg = `go-guard [action] <function>
go-guard generates guard function for a constructor
or constuctor call
Examples:
go-guard func NewRepository
go-guard call NewRepository
`

	fn = `// {{.Name}} allows to guard {{.CName}} constructor.
func {{.Name}}({{range .Params}}{{.Name}} {{.Type}}, {{end}}) {
	{{- range $i, $p := .Params}}
	{{ if .Check -}}
	check.MustNotNil({{add $i 1}},"{{.Name}}",{{.Name}})
	{{- end -}}
	{{- end}}
}`

	call = `{{.Name}}({{range .Params}}{{.Name}}, {{end}})`
)

var (
	fnTmpl = template.Must(template.New("fn").Funcs(template.FuncMap{
		"add": func(x, y int) int {
			return x + y
		},
	}).Parse(fn))

	callTmpl = template.Must(template.New("call").Parse(call))
)

func main() {
	flag.Parse()
	pwd, _ := os.Getwd()
	conf := packages.Config{
		Mode:  packages.LoadSyntax,
		Tests: false,
	}

	checkFlags()

	pkgs, err := packages.Load(&conf, pwd)
	if err != nil {
		exit(fmt.Errorf("load package: %w", err))
	}

	var g guard
	g.CName = flag.Arg(1)
	g.Name = generateName(g.CName)
	parse(pkgs, &g)

	if len(g.Params) == 0 {
		exit(errors.New("function not found or has no parameters"))
	}

	var buf bytes.Buffer
	switch flag.Arg(0) {
	case "call":
		callTmpl.Execute(&buf, &g)
	default:
		fnTmpl.Execute(&buf, &g)
	}

	pretty, err := format.Source(buf.Bytes())
	if err != nil {
		exit(fmt.Errorf("format source: %w", err))
	}

	fmt.Println(string(pretty))
}

func checkFlags() {
	if len(flag.Args()) != 2 {
		usage()
	}

	switch flag.Arg(0) {
	case "func", "call":
	default:
		usage()
	}
}

func parse(pkgs []*packages.Package, g *guard) {
	for _, pkg := range pkgs {
		for _, f := range pkg.Syntax {
			ast.Inspect(f, func(n ast.Node) bool {
				switch d := n.(type) {
				case *ast.FuncDecl:
					if d.Name.Name == g.CName {
						for _, p := range d.Type.Params.List {
							t := pkg.TypesInfo.TypeOf(p.Type)
							prm := param{
								Name: p.Names[0].Name,
								Type: removePath(pkg.PkgPath, t.String()),
							}
							switch t.(type) {
							case *types.Pointer, *types.Interface,
								*types.Slice, *types.Array,
								*types.Map, *types.Chan,
								*types.Signature:
								prm.Check = true
							case *types.Named:
								if _, ok := t.Underlying().(*types.Interface); ok {
									prm.Check = true
								}
							}

							g.Params = append(g.Params, prm)
						}
					}
				}

				return true
			})
		}
	}
}

type guard struct {
	Name   string
	CName  string
	Params []param
}

type param struct {
	Name  string
	Type  string
	Check bool
}

func generateName(base string) string {
	return "guard" + base
}

func removePath(pkgPath string, t string) string {
	parts := strings.Split(t, pkgPath+".")

	switch len(parts) {
	case 2:
		return parts[0] + parts[1]
	default:
		return parts[0]
	}
}

func exit(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(2)
}

func usage() {
	fmt.Fprint(os.Stderr, usg)
	os.Exit(2)
}
