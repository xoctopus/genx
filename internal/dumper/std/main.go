package main

import (
	"bytes"
	_ "embed"
	"io"
	"os"
	"slices"
	"sort"
	"strings"

	"github.com/xoctopus/x/misc/must"
	gopkg "golang.org/x/tools/go/packages"
)

func main() {
	pkgs, err := gopkg.Load(nil, "std")
	must.NoErrorF(err, "failed to scan std library packages")

	sort.Slice(pkgs, func(i, j int) bool {
		return pkgs[i].PkgPath < pkgs[j].PkgPath
	})

	b := bytes.NewBuffer(nil)
	for _, p := range pkgs {
		if slices.Contains(strings.Split(p.PkgPath, "/"), "vendor") ||
			slices.Contains(strings.Split(p.PkgPath, "/"), "internal") {
			continue
		}
		b.WriteString(p.PkgPath)
		b.WriteByte('\n')
	}

	f, err := os.OpenFile("std.list", os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	must.NoErrorF(err, "failed to open std package list file")
	defer f.Close()

	_, err = io.Copy(f, b)
	must.NoErrorF(err, "failed to write std package list file")
}
