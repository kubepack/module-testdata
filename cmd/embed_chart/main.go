package main

import (
	"fmt"
	"io/fs"

	"github.com/kubepack/module-testdata/charts"
)

func main() {
	fs.WalkDir(charts.FS, "", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		fmt.Println(path)
		fmt.Println(d.Name(), "dir =", d.IsDir())
		return nil
	})
}
