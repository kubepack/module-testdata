package action

import (
	"errors"
	"fmt"
	"time"

	actionx "github.com/tamalsaha/hell-flow/pkg/action"
	"github.com/tamalsaha/hell-flow/pkg/values"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chartutil"
	"helm.sh/helm/v3/pkg/release"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	libchart "kubepack.dev/lib-helm/chart"
	"kubepack.dev/lib-helm/repo"
)

type InstallOptions struct {
	ChartURL     string         `json:"chartURL"`
	ChartName    string         `json:"chartName"`
	Version      string         `json:"version"`
	Values       values.Options `json:",inline,omitempty"`
	DryRun       bool           `json:"dryRun"`
	DisableHooks bool           `json:"disableHooks"`
	Replace      bool           `json:"replace"`
	Wait         bool           `json:"wait"`
	Devel        bool           `json:"devel"`
	Timeout      time.Duration  `json:"timeout"`
	Namespace    string         `json:"namespace"`
	ReleaseName  string         `json:"releaseName"`
	Atomic       bool           `json:"atomic"`
	SkipCRDs     bool           `json:"skipCRDs"`
}

type Installer struct {
	cfg *actionx.Configuration

	opts   InstallOptions
	reg    *repo.Registry
	result *release.Release
}

func NewInstaller(getter genericclioptions.RESTClientGetter, namespace string, helmDriver string) (*Installer, error) {
	cfg := new(actionx.Configuration)
	// TODO: Use secret driver for which namespace?
	err := cfg.Init(getter, namespace, helmDriver, debug)
	if err != nil {
		return nil, err
	}
	cfg.Capabilities = chartutil.DefaultCapabilities

	return NewInstallerForConfig(cfg), nil
}

func NewInstallerForConfig(cfg *actionx.Configuration) *Installer {
	return &Installer{
		cfg: cfg,
	}
}

func (x *Installer) WithOptions(opts InstallOptions) *Installer {
	x.opts = opts
	return x
}

func (x *Installer) WithRegistry(reg *repo.Registry) *Installer {
	x.reg = reg
	return x
}

func (x *Installer) Run() (*release.Release, error) {
	if x.opts.Version == "" && x.opts.Devel {
		debug("setting version to >0.0.0-0")
		x.opts.Version = ">0.0.0-0"
	}

	if x.reg == nil {
		return nil, errors.New("x.reg is not set")
	}

	chrt, err := x.reg.GetChart(x.opts.ChartURL, x.opts.ChartName, x.opts.Version)
	if err != nil {
		return nil, err
	}

	cmd := action.NewInstall(&x.cfg.Configuration)
	var extraAPIs []string

	cmd.DryRun = x.opts.DryRun
	cmd.ReleaseName = x.opts.ReleaseName
	cmd.Namespace = x.opts.Namespace
	cmd.Replace = x.opts.Replace // Skip the name check
	cmd.ClientOnly = false
	cmd.APIVersions = chartutil.VersionSet(extraAPIs)
	cmd.Version = x.opts.Version
	cmd.DisableHooks = x.opts.DisableHooks
	cmd.Atomic = x.opts.Atomic
	cmd.Wait = x.opts.Wait
	cmd.Timeout = x.opts.Timeout

	validInstallableChart, err := libchart.IsChartInstallable(chrt.Chart)
	if !validInstallableChart {
		return nil, err
	}

	if chrt.Metadata.Deprecated {
		_, err = fmt.Println("# WARNING: This chart is deprecated")
		if err != nil {
			return nil, err
		}
	}

	if req := chrt.Metadata.Dependencies; req != nil {
		// If CheckDependencies returns an error, we have unfulfilled dependencies.
		// As of Helm 2.4.0, this is treated as a stopping condition:
		// https://github.com/helm/helm/issues/2209
		if err := action.CheckDependencies(chrt.Chart, req); err != nil {
			return nil, err
		}
	}

	vals, err := x.opts.Values.MergeValues(chrt.Chart)
	if err != nil {
		return nil, err
	}
	// chartutil.CoalesceValues(chrt, chrtVals) will use vals to render templates
	chrt.Chart.Values = map[string]interface{}{}

	return cmd.Run(chrt.Chart, vals)
}

func (x *Installer) Do() error {
	var err error
	x.result, err = x.Run()
	return err
}

func (x *Installer) Result() *release.Release {
	return x.result
}
