package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"

	"github.com/tamalsaha/hell-flow/pkg/lib/action"
	"github.com/tamalsaha/hell-flow/pkg/values"

	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery/cached/memory"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/restmapper"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	"k8s.io/client-go/util/homedir"
	"k8s.io/klog/v2"
	"kmodules.xyz/client-go/discovery"
	clientcmdutil "kmodules.xyz/client-go/tools/clientcmd"
	"kubepack.dev/kubepack/pkg/lib"
	"sigs.k8s.io/yaml"
)

func print_yaml() {
	data, err := yaml.Marshal(myflow)
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile("myflow.yaml", data, 0644)
	if err != nil {
		panic(err)
	}
}

func main__() {
	print_yaml()

	masterURL := ""
	kubeconfigPath := filepath.Join(homedir.HomeDir(), ".kube", "config")

	config, err := clientcmd.BuildConfigFromFlags(masterURL, kubeconfigPath)
	if err != nil {
		log.Fatalf("Could not get Kubernetes config: %s", err)
	}

	client := kubernetes.NewForConfigOrDie(config)

	var mapper meta.RESTMapper
	mapper = restmapper.NewDeferredDiscoveryRESTMapper(memory.NewMemCacheClient(client.Discovery()))

	gvr := schema.GroupVersionResource{
		Group:    "admissionregistration.k8s.io",
		Version:  "",
		Resource: "mutatingwebhookconfigurations",
	}
	gvrs, err := mapper.ResourceFor(gvr)
	if err != nil {
		panic(err)
	}
	fmt.Println(gvrs)

	gvk := schema.GroupVersionKind{
		Group:   "admissionregistration.k8s.io",
		Version: "",
		Kind:    "MutatingWebhookConfiguration",
	}
	mappings, err := mapper.RESTMappings(gvk.GroupKind(), "v1alpha1")
	if err != nil {
		if meta.IsNoMatchError(err) {
			fmt.Println(err.(*meta.NoKindMatchError).Error())
			return
		}
		panic(err)
	}
	for _, m2 := range mappings {
		fmt.Println(m2.GroupVersionKind)
	}
}

func main_config_overriding() {
	bc := &BaseConfig{}
	bc.Init("it")

	cc := &ChildConfig{
		bc,
	}
	cc.Init("xyz")
	cc.BaseConfig.Init("xyz")
}

var (
	masterURL      = ""
	kubeconfigPath = filepath.Join(homedir.HomeDir(), ".kube", "config")

	//url     = "https://charts.appscode.com/stable/"
	//name    = "kubedb"
	//version = "v0.13.0-rc.0"

	url     = "https://raw.githubusercontent.com/tamalsaha/hell-flow/master/stable/"
	name    = "first"
	version = "0.1.0"
)

func main() {
	print_yaml()

	cc := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeconfigPath},
		&clientcmd.ConfigOverrides{ClusterInfo: clientcmdapi.Cluster{Server: masterURL}})
	kubeconfig, err := cc.RawConfig()
	if err != nil {
		klog.Fatal(err)
	}
	getter := clientcmdutil.NewClientGetter(&kubeconfig)

	config, err := cc.ClientConfig() // clientcmd.BuildConfigFromFlags(masterURL, kubeconfigPath)
	if err != nil {
		log.Fatalf("Could not get Kubernetes config: %s", err)
	}

	dc := dynamic.NewForConfigOrDie(config)
	gvrNode := schema.GroupVersionResource{
		Group:    "",
		Version:  "v1",
		Resource: "nodes",
	}
	_, err = dc.Resource(gvrNode).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err)
	}

	namespace := "default"
	i, err := action.NewInstaller(getter, namespace, "app")
	if err != nil {
		klog.Fatal(err)
	}
	i.WithRegistry(lib.DefaultRegistry).
		WithOptions(action.InstallOptions{
			ChartURL:  url,
			ChartName: name,
			Version:   version,
			Values: values.Options{
				ValuesFile:  "",
				ValuesPatch: nil,
			},
			DryRun:       false,
			DisableHooks: false,
			Replace:      false,
			Wait:         false,
			Devel:        false,
			Timeout:      0,
			Namespace:    namespace,
			ReleaseName:  name,
			Atomic:       false,
			SkipCRDs:     false,
		})
	rel, err := i.Run()
	if err != nil {
		klog.Fatal(err)
	}
	klog.Infoln(rel)
}

func main_install_or_upgrdae() {
	print_yaml()

	cc := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeconfigPath},
		&clientcmd.ConfigOverrides{ClusterInfo: clientcmdapi.Cluster{Server: masterURL}})
	kubeconfig, err := cc.RawConfig()
	if err != nil {
		klog.Fatal(err)
	}
	getter := clientcmdutil.NewClientGetter(&kubeconfig)

	vt, err := InstallOrUpgrade(getter, "default", ChartLocator{
		URL:     url,
		Name:    name,
		Version: version,
	}, name, "", values.Options{})
	if err != nil {
		klog.Fatal(err)
	}
	klog.Infof("Chart %s", vt)
}

func main____() {
	print_yaml()

	cc := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeconfigPath},
		&clientcmd.ConfigOverrides{ClusterInfo: clientcmdapi.Cluster{Server: masterURL}})
	kubeconfig, err := cc.RawConfig()
	if err != nil {
		klog.Fatal(err)
	}
	getter := clientcmdutil.NewClientGetter(&kubeconfig)

	config, err := cc.ClientConfig() // clientcmd.BuildConfigFromFlags(masterURL, kubeconfigPath)
	if err != nil {
		log.Fatalf("Could not get Kubernetes config: %s", err)
	}

	dc := dynamic.NewForConfigOrDie(config)
	mapper, err := getter.ToRESTMapper()
	if err != nil {
		klog.Fatal(err)
	}

	flowstore := map[string]*FlowState{}

	for _, action := range myflow.Actions {
		runner := ActionRunner{
			dc:        dc,
			mapper:    discovery.NewResourceMapper(mapper),
			flowstore: flowstore,
			FlowName:  myflow.Name,
			action:    action,
			// err:    nil,
		}
		err := runner.Execute()
		if err != nil {
			klog.Fatalln(err)
		}
	}

	gvrNode := schema.GroupVersionResource{
		Group:    "",
		Version:  "v1",
		Resource: "nodes",
	}
	_, err = dc.Resource(gvrNode).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err)
	}

	namespace := "default"
	i, err := action.NewInstaller(getter, namespace, "secret")
	if err != nil {
		klog.Fatal(err)
	}
	i.WithRegistry(lib.DefaultRegistry).
		WithOptions(action.InstallOptions{
			ChartURL:  url,
			ChartName: name,
			Version:   version,
			Values: values.Options{
				ValuesFile:  "",
				ValuesPatch: nil,
			},
			DryRun:       false,
			DisableHooks: false,
			Replace:      false,
			Wait:         false,
			Devel:        false,
			Timeout:      0,
			Namespace:    namespace,
			ReleaseName:  name,
			Atomic:       false,
			SkipCRDs:     false,
		})
	rel, err := i.Run()
	if err != nil {
		klog.Fatal(err)
	}
	klog.Infoln(rel)
}

func main_tpllist() {
	tpls := TemplateList{}
	tpls.Add("xyz")
	tpls.Add("abc")
	fmt.Println(tpls)
}

func main__copy_label_selector() {
	src := metav1.LabelSelector{
		MatchLabels: map[string]string{
			"abc": "xyz",
		},
		MatchExpressions: []metav1.LabelSelectorRequirement{
			{
				Key:      "nodename",
				Operator: metav1.LabelSelectorOpIn,
				Values: []string{
					"node-1",
					"node-2",
				},
			},
			{
				Key:      "hostname",
				Operator: metav1.LabelSelectorOpIn,
				Values: []string{
					"host-1",
					"host-2",
				},
			},
		},
	}

	var sel metav1.LabelSelector
	if src.MatchLabels != nil {
		sel.MatchLabels = make(map[string]string)
	}
	for k, v := range src.MatchLabels {
		sel.MatchLabels[k] = v + "-copy"
	}
	if len(src.MatchExpressions) > 0 {
		sel.MatchExpressions = make([]metav1.LabelSelectorRequirement, 0, len(src.MatchExpressions))
	}
	for _, expr := range src.MatchExpressions {
		ne := expr // src.MatchExpressions[i]
		ne.Values = make([]string, 0, len(expr.Values))
		for _, v := range expr.Values {
			ne.Values = append(ne.Values, v+"-copy")
		}
		sel.MatchExpressions = append(sel.MatchExpressions, ne)
	}
	fmt.Println(sel)
}

// must use int64
func main_set_values() {
	u := unstructured.Unstructured{
		Object: map[string]interface{}{
			"a": map[string]interface{}{
				"b": int64(2),
			},
		},
	}
	override := map[string]interface{}{
		"b": int64(20),
		"c": "c3",
	}
	err := unstructured.SetNestedField(u.Object, override, "a")
	if err != nil {
		panic(err)
	}
}

func main_print_yaml() {
	print_yaml()
}
