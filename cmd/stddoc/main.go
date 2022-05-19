// Prints a list of strings representing symbols from the standard library.
//
// PS:		go run . | Set-Content $HOME/goapi.txt
// BASH:	go run . >> ~/goapi.txt
//
// Each item in the list is a valid argument for the 'go doc' tool. And now
// a fuzzy finder (e.g. fzf) can be used to access the documentation.
//
// PS:		cat $HOME/goapi.txt | fzf.exe --preview 'go doc {}' --preview-window 'right:67%'
// BASH:	cat ~/goapi.txt | fzf --preview 'go doc {}' --preview-window 'right:67%'
//
// PS:
// function hh
// {
//   $name = Get-Content $HOME/goapi.txt | . $HOME/portable/fzf.exe --preview 'go doc {}' --preview-window 'right:67%'
//   go doc $name
// }
//
// BASH:
// hh() {
//     local entry
//
//     entry="$(cat ~/goapi.txt | fzf --preview 'go doc {}' --preview-window 'right:67%')"
//     go doc $entry
// }

package main

import (
	"bufio"
	"fmt"
	"go/build"
	"go/doc"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
)

// Returns collection of things defined in the given package
func processPackage(pkg *build.Package) []string {
	include := func(info fs.FileInfo) bool {
		for _, name := range pkg.GoFiles {
			if name == info.Name() {
				return true
			}
		}
		for _, name := range pkg.CgoFiles {
			if name == info.Name() {
				return true
			}
		}
		return false
	}

	var res []string

	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, pkg.Dir, include, parser.ParseComments)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't parse the directory %s: %s\n", pkg.Dir, err)
		os.Exit(1)
	}

	if len(pkgs) != 1 {
		return res
	}
	astPkg := pkgs[pkg.Name]

	docPkg := doc.New(astPkg, pkg.ImportPath, doc.Mode(0))
	res = append(res, docPkg.ImportPath)

	for _, value := range docPkg.Consts {
		for _, name := range value.Names {
			res = append(res, fmt.Sprintf("%s.%s", docPkg.ImportPath, name))
			break
		}
	}
	for _, value := range docPkg.Vars {
		for _, name := range value.Names {
			res = append(res, fmt.Sprintf("%s.%s", docPkg.ImportPath, name))
			break
		}
	}

	for _, fun := range docPkg.Funcs {
		res = append(res, fmt.Sprintf("%s.%s", docPkg.ImportPath, fun.Name))
	}

	for _, typ := range docPkg.Types {
		res = append(res, fmt.Sprintf("%s.%s", docPkg.ImportPath, typ.Name))

		for _, value := range typ.Consts {
			for _, name := range value.Names {
				res = append(res, fmt.Sprintf("%s.%s", docPkg.ImportPath, name))
				break
			}
		}
		for _, value := range typ.Vars {
			for _, name := range value.Names {
				res = append(res, fmt.Sprintf("%s.%s", docPkg.ImportPath, name))
				break
			}
		}
		for _, fun := range typ.Funcs {
			res = append(res, fmt.Sprintf("%s.%s", docPkg.ImportPath, fun.Name))
		}
		for _, fun := range typ.Methods {
			res = append(res, fmt.Sprintf("%s.%s.%s", docPkg.ImportPath, typ.Name, fun.Name))
		}
	}

	return res
}

func main() {
	var root = filepath.Join(runtime.GOROOT(), "src")
	var names []string

	err := filepath.WalkDir(root, func(path string, info fs.DirEntry, err error) error {
		if !info.IsDir() {
			return nil
		}

		base := info.Name()
		if base[0] == '.' || base[0] == '_' || base == "testdata" || base == "vendor" || base == "internal" {
			return fs.SkipDir
		}

		// process only packages that can be imported
		buildPkg, importErr := build.ImportDir(path, build.ImportComment)
		if importErr == nil {
			x := processPackage(buildPkg)
			names = append(names, x...)
		}
		return nil
	})

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	stdout := bufio.NewWriter(os.Stdout)
	for _, name := range names {
		fmt.Fprintf(stdout, "%s\n", name)
	}
	stdout.Flush()
}
