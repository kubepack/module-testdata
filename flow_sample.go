package main

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	rsapi "kmodules.xyz/resource-metadata/apis/meta/v1alpha1"
	pkgapi "kubepack.dev/module/apis/pkg/v1alpha1"
)

var myflow = &pkgapi.Flow{
	TypeMeta: metav1.TypeMeta{
		APIVersion: "pkg.kubepack.com/v1alpha1",
		Kind:       "Flow",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name:      "myflow",
		Namespace: "demo",
	},
	Spec: pkgapi.FlowSpec{
		Actions: []pkgapi.Action{
			{
				ReleaseName: "first",
				ChartRepoRef: rsapi.ChartRepoRef{
					URL:     "https://raw.githubusercontent.com/kubepack/module-testdata/master/stable/",
					Name:    "first",
					Version: "0.1.0",
				},
				ValuesFile:     "",
				ValuesPatch:    nil,
				ValueOverrides: nil,
				Prerequisites: pkgapi.Prerequisites{
					RequiredResources: []metav1.GroupVersionResource{
						{Group: "apps", Version: "v1", Resource: "deployments"},
					},
				},
				ReadinessCriteria: pkgapi.ReadinessCriteria{
					WaitForReconciled: true,
					// check for installed crd
					ResourcesExist: nil,
					// Wait until LB IP is set
					WaitFors: []pkgapi.WaitFlags{
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
				ChartRepoRef: rsapi.ChartRepoRef{
					URL:     "https://raw.githubusercontent.com/kubepack/module-testdata/master/stable/",
					Name:    "third",
					Version: "0.1.0",
				},
				ValuesFile:  "",
				ValuesPatch: nil,
				/*
				  export POD_NAME=$(kubectl get pods --namespace default -l "app.kubernetes.io/name=first,app.kubernetes.io/instance=first" -o jsonpath="{.items[0].metadata.name}")
				  export CONTAINER_PORT=$(kubectl get pod --namespace default $POD_NAME -o jsonpath="{.spec.containers[0].ports[0].containerPort}")
				  echo "Visit http://127.0.0.1:8080 to use your application"
				  kubectl --namespace default port-forward $POD_NAME 8080:$CONTAINER_PORT
				*/
				ValueOverrides: []pkgapi.LoadValue{
					{
						From: pkgapi.ObjectLocator{
							UseRelease: "first",
							Src: pkgapi.ObjectRef{
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
								Name:         "",
								NameTemplate: "",
							},
							Paths: nil,
						},
						Values: []pkgapi.KV{
							{
								Key:          "first.name",
								Type:         "string",
								PathTemplate: ``,
								Path:         ".metadata.name",
							},
							{
								Key:          "first.port",
								Type:         "string",
								PathTemplate: `{{ jp "{.spec.containers[0].ports[0].containerPort}" . }}`,
								Path:         "",
							},
						},
					},
				},
				Prerequisites: pkgapi.Prerequisites{
					RequiredResources: []metav1.GroupVersionResource{
						{Group: "apps", Version: "v1", Resource: "deployments"},
					},
				},
				ReadinessCriteria: pkgapi.ReadinessCriteria{
					WaitForReconciled: true,
					// check for installed crd
					ResourcesExist: nil,
					// Wait until LB IP is set
					WaitFors: []pkgapi.WaitFlags{
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
		EdgeList: []rsapi.NamedEdge{
			{
				Name:       "",
				Src:        metav1.TypeMeta{},
				Dst:        metav1.TypeMeta{},
				Connection: rsapi.ResourceConnectionSpec{},
			},
		},
	},
}
