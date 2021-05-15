package main

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"kmodules.xyz/resource-metadata/apis/meta/v1alpha1"
)

type Flow struct {
	Actions  []Action       `json:"actions"`
	EdgeList []DirectedEdge `json:"edge_list"`
}

// Check array, map, etc
// can this be always string like in --set keys?
// Keep is such that we can always generate helm equivalent command
type KV struct {
	Key string
	// type is an OpenAPI type definition for this column.
	// See https://github.com/OAI/OpenAPI-Specification/blob/master/versions/2.0.md#data-types for more.
	Type string `json:"type"`
	// format is an optional OpenAPI type definition for this column. The 'name' format is applied
	// to the primary identifier column to assist in clients identifying column is the resource name.
	// See https://github.com/OAI/OpenAPI-Specification/blob/master/versions/2.0.md#data-types for more.
	// +optional
	Format string `json:"format,omitempty"`
	// PathTemplate is a Go text template that will be evaluated to determine cell value.
	// Users can use JSONPath expression to extract nested fields and apply template functions from Masterminds/sprig library.
	// The template function for JSON path is called `jp`.
	// Example: {{ jp "{.a.b}" . }} or {{ jp "{.a.b}" true }}, if json output is desired from JSONPath parser
	// +optional
	PathTemplate string `json:"pathTemplate,omitempty"`
	//
	//
	// Directly use path from object
	Path string `json:"path"`
}

type LoadValue struct {
	From   ObjectLocator `json:"from"`
	Values []KV          `json:"values"`
}

type ObjectLocator struct {
	// Use the values from that release == action to render templates
	UseRelease string    `json:"use_release"`
	Src        ObjectRef `json:"src"`
	Paths      []string  `json:"paths"` // sequence of DirectedEdge names
}

type DirectedEdge struct {
	Name       string
	Src        metav1.TypeMeta
	Dst        metav1.TypeMeta
	Connection v1alpha1.ResourceConnectionSpec
}

type ObjectRef struct {
	Target       metav1.TypeMeta       `json:"target"`
	Selector     *metav1.LabelSelector `json:"selector,omitempty"`
	Name         string                `json:"name,omitempty"`
	NameTemplate string                `json:"nameTemplate,omitempty"`
	// Namespace always same as Workflow
}

type Action struct {
	// Also the action name
	ReleaseName string `json:"releaseName" protobuf:"bytes,3,opt,name=releaseName"`

	ChartRef `json:",inline" protobuf:"bytes,1,opt,name=chartRef"`
	Version  string `json:"version" protobuf:"bytes,2,opt,name=version"`

	// Namespace   string `json:"namespace" protobuf:"bytes,4,opt,name=namespace"`

	ValuesFile string `json:"valuesFile,omitempty" protobuf:"bytes,6,opt,name=valuesFile"`
	// RFC 6902 compatible json patch. ref: http://jsonpatch.com
	// +optional
	// +kubebuilder:pruning:PreserveUnknownFields
	ValuesPatch *runtime.RawExtension `json:"valuesPatch,omitempty" protobuf:"bytes,7,opt,name=valuesPatch"`

	ValueOverrides []LoadValue `json:"overrideValues"`

	// https://github.com/tamalsaha/kstatus-demo
	ReadinessCriteria ReadinessCriteria `json:"readiness_criteria"`

	Prerequisites Prerequisites `json:"prerequisites"`
}

type Prerequisites struct {
	RequiredResources []ResourceID `json:"required_resources"`
}

type ReadinessCriteria struct {
	Timeout metav1.Duration `json:"timeout"`

	// List objects for which to wait to reconcile using kstatus == Current
	// Same as helm --wait
	WaitForReconciled bool `json:"wait_for_reconciled"`

	ResourcesExist []ResourceID `json:"required_resources"`
	WaitFors       []WaitFlags  `json:"waitFors,omitempty" protobuf:"bytes,9,rep,name=waitFors"`
}

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
