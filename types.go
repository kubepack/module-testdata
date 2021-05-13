package main

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type ChartRef struct {
	URL  string `json:"url" protobuf:"bytes,1,opt,name=url"`
	Name string `json:"name" protobuf:"bytes,2,opt,name=name"`
}

type ChartSelection struct {
	ChartRef    `json:",inline" protobuf:"bytes,1,opt,name=chartRef"`
	Version     string `json:"version" protobuf:"bytes,2,opt,name=version"`
	ReleaseName string `json:"releaseName" protobuf:"bytes,3,opt,name=releaseName"`
	Namespace   string `json:"namespace" protobuf:"bytes,4,opt,name=namespace"`

	ValuesFile string `json:"valuesFile,omitempty" protobuf:"bytes,6,opt,name=valuesFile"`
	// RFC 6902 compatible json patch. ref: http://jsonpatch.com
	// +optional
	// +kubebuilder:pruning:PreserveUnknownFields
	ValuesPatch *runtime.RawExtension `json:"valuesPatch,omitempty" protobuf:"bytes,7,opt,name=valuesPatch"`
	Resources   *ResourceDefinitions  `json:"resources,omitempty" protobuf:"bytes,8,opt,name=resources"`
	WaitFors    []WaitFlags           `json:"waitFors,omitempty" protobuf:"bytes,9,rep,name=waitFors"`
}

// ResourceID identifies a resource
type ResourceID struct {
	Group   string `json:"group" protobuf:"bytes,1,opt,name=group"`
	Version string `json:"version" protobuf:"bytes,2,opt,name=version"`
	// Name is the plural name of the resource to serve.  It must match the name of the CustomResourceDefinition-registration
	// too: plural.group and it must be all lowercase.
	Resource string `json:"resource" protobuf:"bytes,3,opt,name=resource"`
}

type Feature struct {
	Trait string `json:"trait" protobuf:"bytes,1,opt,name=trait"`
	Value string `json:"value" protobuf:"bytes,2,opt,name=value"`
}

type ResourceDefinitions struct {
	Owned    []ResourceID `json:"owned" protobuf:"bytes,1,rep,name=owned"`
	Required []ResourceID `json:"required" protobuf:"bytes,2,rep,name=required"`
}

// wait ([-f FILENAME] | resource.group/resource.name | resource.group [(-l label | --all)]) [--for=delete|--for condition=available]

type WaitFlags struct {
	Resource     GroupResource         `json:"resource" protobuf:"bytes,1,opt,name=resource"`
	Labels       *metav1.LabelSelector `json:"labels" protobuf:"bytes,2,opt,name=labels"`
	All          bool                  `json:"all" protobuf:"varint,3,opt,name=all"`
	Timeout      metav1.Duration       `json:"timeout" protobuf:"bytes,4,opt,name=timeout"`
	ForCondition string                `json:"for" protobuf:"bytes,5,opt,name=for"`
}

type GroupVersionResource struct {
	Group    string `json:"group" protobuf:"bytes,1,opt,name=group"`
	Version  string `json:"version" protobuf:"bytes,2,opt,name=version"`
	Resource string `json:"resource" protobuf:"bytes,3,opt,name=resource"`
}

type GroupResource struct {
	Group string `json:"group" protobuf:"bytes,1,opt,name=group"`
	Name  string `json:"name" protobuf:"bytes,2,opt,name=name"`
}
