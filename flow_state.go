package main

import (
	"helm.sh/helm/v3/pkg/chart"
)

type FlowState struct {
	ReleaseName string
	Chrt        *chart.Chart
	Values      map[string]interface{} // Final Values used
	//
}

// release name -> flow state
var flowStates = map[string]*FlowState{}
