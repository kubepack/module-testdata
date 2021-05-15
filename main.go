package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"

	"github.com/tamalsaha/hell-flow/pkg/lib/action"

	"gomodules.xyz/x/crypto/rand"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery/cached/memory"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/restmapper"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	"k8s.io/client-go/util/homedir"
	"k8s.io/klog/v2"
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

func main__hcart_install() {
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
	i, err := action.NewInstaller(getter, namespace, "secret")
	if err != nil {
		klog.Fatal(err)
	}
	i.WithRegistry(lib.DefaultRegistry).
		WithOptions(action.InstallOptions{
			ChartURL:     url,
			ChartName:    name,
			Version:      version,
			ValuesFile:   "",
			ValuesPatch:  nil,
			DryRun:       false,
			DisableHooks: false,
			Replace:      false,
			Wait:         false,
			Devel:        false,
			Timeout:      0,
			Namespace:    namespace,
			ReleaseName:  rand.WithUniqSuffix(name),
			Atomic:       false,
			SkipCRDs:     false,
		})
	rel, err := i.Run()
	if err != nil {
		klog.Fatal(err)
	}
	klog.Infoln(rel)
}

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

	vt, err := InstallOrUpgrade(getter, "default", ChartLocator{
		URL:     url,
		Name:    name,
		Version: version,
	}, name)
	if err != nil {
		klog.Fatal(err)
	}
	klog.Infof("Chart %s", vt)
}
