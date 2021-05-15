package main

import (
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chartutil"
	"helm.sh/helm/v3/pkg/engine"
)

type FlowState struct {
	ReleaseName string
	Chrt        *chart.Chart
	Values      chartutil.Values // final values used for rendering
	Engine      engine.EngineInstance
}
