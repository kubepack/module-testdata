package main

import (
	"fmt"
	"github.com/kubepack/module-testdata/charts"
	"io/fs"
	"k8s.io/klog/v2"
)

func main() {
	//fi, err := os.Lstat("/home/tamal/go/src/kubepack.dev/module-testdata/charts/first/templates/service.yaml")
	//if err != nil {
	//	klog.Fatalln(err)
	//}
	//fmt.Println(fi.Name())

	first, err := fs.Sub(charts.FS, "first")
	if err != nil {
		klog.Fatalln(err)
	}
	fmt.Println(first)

	fs.WalkDir(charts.FS, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		fmt.Println(path)
		// fmt.Println(d.Name(), "dir =", d.IsDir())
		return nil
	})
}
