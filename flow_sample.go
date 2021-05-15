package main

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"kmodules.xyz/resource-metadata/apis/meta/v1alpha1"
)

var myflow = &Flow{
	Actions: []Action{
		{
			ReleaseName: "first",
			ChartRef: ChartRef{
				URL:  "https://raw.githubusercontent.com/tamalsaha/hell-flow/master/stable/",
				Name: "first",
			},
			Version:        "0.1.0",
			ValuesFile:     "",
			ValuesPatch:    nil,
			ValueOverrides: nil,
			Prerequisites: Prerequisites{
				RequiredResources: []ResourceID{
					{Group: "apps", Version: "v1", Resource: "deployments"},
				},
			},
			ReadinessCriteria: &ReadinessCriteria{
				WaitForReconciled: true,
				// check for installed crd
				ResourcesExist: nil,
				// Wait until LB IP is set
				WaitFors: []WaitFlags{
					//{
					//	Resource:     GroupResource{
					//		Group: "",
					//		Name:  "",
					//	},
					//	Labels:       nil,
					//	All:          false,
					//	Timeout:      metav1.Duration{},
					//	ForCondition: "",
					//},
				},
			},
		},
		{
			ReleaseName: "third",
			ChartRef: ChartRef{
				URL:  "https://raw.githubusercontent.com/tamalsaha/hell-flow/master/stable/",
				Name: "third",
			},
			Version:     "0.1.0",
			ValuesFile:  "",
			ValuesPatch: nil,
			/*
			  export POD_NAME=$(kubectl get pods --namespace default -l "app.kubernetes.io/name=first,app.kubernetes.io/instance=first" -o jsonpath="{.items[0].metadata.name}")
			  export CONTAINER_PORT=$(kubectl get pod --namespace default $POD_NAME -o jsonpath="{.spec.containers[0].ports[0].containerPort}")
			  echo "Visit http://127.0.0.1:8080 to use your application"
			  kubectl --namespace default port-forward $POD_NAME 8080:$CONTAINER_PORT
			*/
			ValueOverrides: []LoadValue{
				{
					From: ObjectLocator{
						UseRelease: "first",
						Src: ObjectRef{
							Target: metav1.TypeMeta{
								Kind:       "Pod",
								APIVersion: "v1",
							},
							Selector: &metav1.LabelSelector{
								MatchLabels: map[string]string{
									"app.kubernetes.io/name":     "{{ .Release.Name }}",
									"app.kubernetes.io/instance": "{{ .Release.Name }}",
								},
							},
							Name:         nil,
							NameTemplate: nil,
						},
						Paths: nil,
					},
					Values: []KV{
						{
							Key:          "first.name",
							Type:         "string",
							Format:       "",
							PathTemplate: ``,
							Path:         ".metadata.name",
						},
						{
							Key:          "first.port",
							Type:         "string",
							Format:       "",
							PathTemplate: `{{ jp "{.spec.containers[0].ports[0].containerPort}" . }}`,
							Path:         "",
						},
					},
				},
			},
			Prerequisites: Prerequisites{
				RequiredResources: []ResourceID{
					{Group: "apps", Version: "v1", Resource: "pods"},
				},
			},
			ReadinessCriteria: &ReadinessCriteria{
				WaitForReconciled: true,
				// check for installed crd
				ResourcesExist: nil,
				// Wait until LB IP is set
				WaitFors: []WaitFlags{
					//{
					//	Resource:     GroupResource{
					//		Group: "",
					//		Name:  "",
					//	},
					//	Labels:       nil,
					//	All:          false,
					//	Timeout:      metav1.Duration{},
					//	ForCondition: "",
					//},
				},
			},
		},
	},
	EdgeList: []DirectedEdge{
		{
			Name:       "",
			Src:        metav1.TypeMeta{},
			Dst:        metav1.TypeMeta{},
			Connection: v1alpha1.ResourceConnectionSpec{},
		},
	},
}
