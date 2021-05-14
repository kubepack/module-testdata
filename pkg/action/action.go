/*
Copyright The Helm Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package action

import (
	"fmt"
	"os"

	kubex "github.com/tamalsaha/hell-flow/pkg/kube"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/storage"
	"helm.sh/helm/v3/pkg/storage/driver"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/discovery/cached/memory"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	driver2 "kubepack.dev/lib-app/pkg/storage/driver"
	appcs "sigs.k8s.io/application/client/clientset/versioned"
)

// Configuration injects the dependencies that all actions share.
type Configuration struct {
	*action.Configuration
}

// Init initializes the action configuration
func (c *Configuration) Init(getter genericclioptions.RESTClientGetter, namespace, helmDriver string, log action.DebugLog) error {
	//kc2 := kube.New(getter)
	//kc2.Log = log
	kc := kubex.New(getter, log)

	lazyClient := &lazyClient{
		namespace: namespace,
		clientFn:  kc.Factory.KubernetesClientSet,
		appClientFn: func() (*appcs.Clientset, error) {
			config, err := kc.Factory.ToRawKubeConfigLoader().ClientConfig()
			if err != nil {
				return nil, err
			}
			return appcs.NewForConfig(config)
		},
	}

	var store *storage.Storage
	switch helmDriver {
	case "app", "apps", "application", "applications":
		config, err := kc.Factory.ToRawKubeConfigLoader().ClientConfig()
		if err != nil {
			return err
		}
		d := driver2.NewApplications(
			newApplicationClient(lazyClient),
			dynamic.NewForConfigOrDie(config),
			memory.NewMemCacheClient(kubernetes.NewForConfigOrDie(config).Discovery()),
		)
		d.Log = log
		store = storage.Init(d)
	case "secret", "secrets", "":
		d := driver.NewSecrets(newSecretClient(lazyClient))
		d.Log = log
		store = storage.Init(d)
	case "configmap", "configmaps":
		d := driver.NewConfigMaps(newConfigMapClient(lazyClient))
		d.Log = log
		store = storage.Init(d)
	case "memory":
		var d *driver.Memory
		if c.Releases != nil {
			if mem, ok := c.Releases.Driver.(*driver.Memory); ok {
				// This function can be called more than once (e.g., helm list --all-namespaces).
				// If a memory driver was already initialized, re-use it but set the possibly new namespace.
				// We re-use it in case some releases where already created in the existing memory driver.
				d = mem
			}
		}
		if d == nil {
			d = driver.NewMemory()
		}
		d.SetNamespace(namespace)
		store = storage.Init(d)
	case "sql":
		d, err := driver.NewSQL(
			os.Getenv("HELM_DRIVER_SQL_CONNECTION_STRING"),
			log,
			namespace,
		)
		if err != nil {
			panic(fmt.Sprintf("Unable to instantiate SQL driver: %v", err))
		}
		store = storage.Init(d)
	default:
		// Not sure what to do here.
		panic("Unknown driver in HELM_DRIVER: " + helmDriver)
	}

	c.RESTClientGetter = getter
	c.KubeClient = kc
	c.Releases = store
	c.Log = log

	return nil
}
