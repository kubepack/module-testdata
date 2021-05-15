package main

import (
	"bytes"
	"context"
	"strings"
	"time"

	"github.com/tamalsaha/hell-flow/pkg/values"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/dynamic"
	"k8s.io/klog/v2"
	"kmodules.xyz/client-go/discovery"
	"kubepack.dev/kubepack/apis/kubepack/v1alpha1"
	"kubepack.dev/kubepack/pkg/lib"
)

type ActionRunner struct {
	dc           dynamic.Interface
	ClientGetter genericclioptions.RESTClientGetter
	mapper       discovery.ResourceMapper
	flowstore    map[string]*FlowState

	Namespace string
	action    Action
	err       error
}

func (runner *ActionRunner) Execute() error {
	if runner.MeetsPrerequisites() {
		runner.Apply().WaitUntilReady()
		if runner.Err() != nil {
			// if ae, ok := runner.Err().(AlreadyErrored); ok {
			// action was already errored out, Rest before reusing
		}
	} else {
		if runner.Err() != nil {
			// if ae, ok := runner.Err().(AlreadyErrored); ok {
			// action was already errored out, Rest before reusing
		}
	}
	return nil
}

// Do we need this?
func (e *ActionRunner) ResetError() {
	e.err = nil
}

func (e *ActionRunner) Err() error {
	return e.err
}

func (e *ActionRunner) MeetsPrerequisites() bool {
	if e.err != nil {
		e.err = NewAlreadyErrored(e.err)
		return false
	}

	return e.resourceExists(context.TODO(), e.action.Prerequisites.RequiredResources)
}

func (e *ActionRunner) Apply() *ActionRunner {
	if e.err != nil {
		e.err = NewAlreadyErrored(e.err)
		return e
	}

	ref := ChartLocator{
		URL:     e.action.URL,
		Name:    e.action.Name,
		Version: e.action.Version,
	}

	chrt, err := lib.DefaultRegistry.GetChart(ref.URL, ref.Name, ref.Version)
	if err != nil {
		e.err = err
		return e
	}

	opts := values.Options{
		// ReplaceValues: nil,
		ValuesFile:   e.action.ValuesFile,
		ValuesPatch:  e.action.ValuesPatch,
		StringValues: nil,
		Values:       nil,
		KVPairs:      nil,
	}

	//for _, v := range e.action.ValueOverrides {
	//	// v.From.Src
	//	// v.From.UseRelease
	//	// v.From.Paths
	//}
	//

	vals, err := opts.MergeValues(chrt.Chart)
	if err != nil {
		e.err = err
		return e
	}

	vt, err := InstallOrUpgrade(e.ClientGetter, e.Namespace, ref, e.action.ReleaseName, values.Options{
		ReplaceValues: vals,
	})
	if err != nil {
		e.err = err
	}
	klog.Infoln("chart %+v %s", ref, vt)

	return e
}

func (e *ActionRunner) WaitUntilReady() {
	if e.err != nil {
		e.err = NewAlreadyErrored(e.err)
		return
	}

	if e.action.ReadinessCriteria.Timeout.Duration == 0 {
		e.action.ReadinessCriteria.Timeout = metav1.Duration{Duration: 15 * time.Minute}
	}
	// start := time.Now()
	// calculate timeout

	ctx, cancel := context.WithTimeout(context.TODO(), e.action.ReadinessCriteria.Timeout.Duration)
	defer cancel()

	rready := e.resourceExists(ctx, e.action.Prerequisites.RequiredResources)
	if e.err != nil {
		return
	}
	if !rready {
		return
	}

	// WaitForFlags
	waitflags := make([]v1alpha1.WaitFlags, 0, len(e.action.ReadinessCriteria.WaitFors))
	for _, w := range e.action.ReadinessCriteria.WaitFors {
		waitflags = append(waitflags, v1alpha1.WaitFlags{
			Resource: v1alpha1.GroupResource{
				Group: w.Resource.Group,
				Name:  w.Resource.Name,
			},
			Labels:       w.Labels,
			All:          w.All,
			Timeout:      w.Timeout,
			ForCondition: w.ForCondition,
		})
	}

	var buf bytes.Buffer
	printer := lib.WaitForPrinter{
		Name:      e.action.ReleaseName,
		Namespace: e.Namespace,
		WaitFors:  waitflags,
		W:         &buf,
	}
	err := printer.Do()
	if err != nil {
		e.err = err
		return
	}
	if buf.Len() > 0 {
		klog.Infoln("running commands:")
		for _, line := range strings.Split(buf.String(), "\n") {
			klog.Infoln(line)
		}
	}

	checker := lib.WaitForChecker{
		Namespace:    e.Namespace,
		WaitFors:     waitflags,
		ClientGetter: e.ClientGetter,
	}
	err = checker.Do()
	if err != nil {
		e.err = err
		return
	}
}

func (e *ActionRunner) resourceExists(ctx context.Context, resources []ResourceID) bool {
	for _, r := range resources {
		exists, err := IsResourceExistsAndReady(ctx, e.dc, e.mapper, schema.GroupVersionResource{
			Group:    r.Group,
			Version:  r.Version,
			Resource: r.Resource,
		})
		if err != nil {
			e.err = err
			return false
		}
		if !exists {
			return false
		}
	}
	return true
}

func (e *ActionRunner) IsReady() bool {
	if e.err != nil {
		e.err = NewAlreadyErrored(e.err)
		return false
	}

	return false
}

type AlreadyErrored struct {
	underlying error
}

func NewAlreadyErrored(underlying error) error {
	if _, ok := underlying.(AlreadyErrored); ok {
		return underlying
	}
	return AlreadyErrored{underlying: underlying}
}

func (e AlreadyErrored) Error() string {
	return e.underlying.Error()
}

func (e AlreadyErrored) Underlying() error {
	return e.underlying
}
