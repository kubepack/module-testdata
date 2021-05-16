package main

import (
	"fmt"
	"log"

	actionx "github.com/tamalsaha/hell-flow/pkg/action"
	"github.com/tamalsaha/hell-flow/pkg/lib/action"
	"github.com/tamalsaha/hell-flow/pkg/values"

	haction "helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chartutil"
	"helm.sh/helm/v3/pkg/storage/driver"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/klog/v2"
	kutil "kmodules.xyz/client-go"
	"kubepack.dev/kubepack/pkg/lib"
)

type ChartLocator struct {
	URL     string
	Name    string
	Version string
}

func InstallOrUpgrade(getter genericclioptions.RESTClientGetter, namespace string, ref ChartLocator, releaseName, partOf string, opts values.Options) (kutil.VerbType, error) {
	helmDriver := "secrets"

	cfg := new(actionx.Configuration)
	// TODO: Use secret driver for which namespace?
	err := cfg.Init(getter, namespace, helmDriver, debug)
	if err != nil {
		return kutil.VerbUnchanged, err
	}
	cfg.Capabilities = chartutil.DefaultCapabilities

	// If a release does not exist, install it.
	histClient := haction.NewHistory(&cfg.Configuration)
	histClient.Max = 1
	if _, err := histClient.Run(releaseName); err == driver.ErrReleaseNotFound {
		i := action.NewInstallerForConfig(cfg)
		i.WithRegistry(lib.DefaultRegistry).
			WithOptions(action.InstallOptions{
				ChartURL:     ref.URL,
				ChartName:    ref.Name,
				Version:      ref.Version,
				Values:       opts,
				DryRun:       false,
				DisableHooks: false,
				Replace:      false,
				Wait:         false,
				Devel:        false,
				Timeout:      0,
				Namespace:    namespace,
				ReleaseName:  releaseName,
				Atomic:       false,
				SkipCRDs:     false,
				PartOf:       partOf,
			})
		rel, err := i.Run()
		if err != nil {
			return kutil.VerbUnchanged, err
		}
		klog.Infoln(rel)
		return kutil.VerbCreated, err // Installed
	} else if err != nil {
		return kutil.VerbUnchanged, err
	}

	i := action.NewUpgraderForConfig(cfg)
	i.WithRegistry(lib.DefaultRegistry).
		WithReleaseName(releaseName).
		WithOptions(action.UpgradeOptions{
			ChartURL:      url,
			ChartName:     name,
			Version:       version,
			Values:        opts,
			Install:       false,
			Devel:         false,
			Namespace:     namespace,
			Timeout:       0,
			Wait:          false,
			DisableHooks:  false,
			DryRun:        false,
			Force:         false,
			ResetValues:   false,
			ReuseValues:   false,
			Recreate:      false,
			MaxHistory:    0,
			Atomic:        false,
			CleanupOnFail: false,
			PartOf:        partOf,
		})
	rel, err := i.Run()
	if err != nil {
		return kutil.VerbUnchanged, err
	}
	klog.Infoln(rel)
	return kutil.VerbUpdated, err // Upgraded
}

func debug(format string, v ...interface{}) {
	format = fmt.Sprintf("[debug] %s\n", format)
	_ = log.Output(2, fmt.Sprintf(format, v...))
}
