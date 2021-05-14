package main

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"kmodules.xyz/resource-metadata/apis/meta/v1alpha1"
)

var myflow = &Flow{
	Actions:  []Action{
		{
			ReleaseName:       "first",
			ChartRef:          ChartRef{
				URL:  "https://raw.githubusercontent.com/tamalsaha/hell-flow/master/stable/",
				Name: "first",
			},
			Version:           "0.1.0",
			ValuesFile:        "",
			ValuesPatch:       nil,
			Prerequisites:     Prerequisites{
				RequiredResources: []ResourceID{
					{Group: "apps", Version: "v1", Resource: "deployments"},
				},
			},
			OverrideValues:    nil,
			ReadinessCriteria: &ReadinessCriteria{
				HelmWait:          true,
				WaitForReconciled: true,
				// check for installed crd
				RequiredResources: nil,
				// Wait until LB IP is set
				WaitFors:          []WaitFlags{
					//{
					//	Resource:     GroupResource{
					//		Group: "",
					//		Name:  "",
					//	},
					//	Labels:       nil,
					//	All:          false,
					//	Timeout:      v1.Duration{},
					//	ForCondition: "",
					//},
				},
			},
		},
	},
	EdgeList: []DirectedEdge{
		{
			Name:       "",
			Src:        v1.TypeMeta{},
			Dst:        v1.TypeMeta{},
			Connection: v1alpha1.ResourceConnectionSpec{},
		},
	},
}
