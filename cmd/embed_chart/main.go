package main

import (
	"fmt"
	"io/fs"
	"net/url"

	"github.com/kubepack/module-testdata/charts"
	"k8s.io/klog/v2"
)

func main() {
	u2, err := url.Parse("embed:///first")
	if err != nil {
		klog.Fatalln(err)
	}
	fmt.Println(u2.Scheme)
	fmt.Println(u2.Path)

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
