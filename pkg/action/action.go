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
	driver2 "github.com/tamalsaha/hell-flow/pkg/storage/driver"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/kube"
	"helm.sh/helm/v3/pkg/storage"
	"helm.sh/helm/v3/pkg/storage/driver"
	crd_cs "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/dynamic"
	"kmodules.xyz/client-go/apiextensions"
	disco_util "kmodules.xyz/client-go/discovery"
	"kubepack.dev/kubepack/apis/kubepack/v1alpha1"
	appcs "sigs.k8s.io/application/client/clientset/versioned"
)

// Configuration injects the dependencies that all actions share.
type Configuration struct {
	action.Configuration
}

// Init initializes the action configuration
func (c *Configuration) Init(getter genericclioptions.RESTClientGetter, namespace, helmDriver string, log action.DebugLog) error {
	var kc kube.Interface
	var factory kube.Factory

	switch helmDriver {
	case "secret", "secrets", "", "configmap", "configmaps", "memory", "sql":
		client := kube.New(getter)
		client.Log = log
		kc = client
		factory = client.Factory
	default:
		// register Application CRD
		crds := []*apiextensions.CustomResourceDefinition{
			v1alpha1.ApplicationCustomResourceDefinition(),
		}
		restcfg, err := getter.ToRESTConfig()
		if err != nil {
			return fmt.Errorf("failed to get rest config, reason %v", err)
		}
		crdClient, err := crd_cs.NewForConfig(restcfg)
		if err != nil {
			return fmt.Errorf("failed to create crd client, reason %v", err)
		}
		err = apiextensions.RegisterCRDs(crdClient, crds)
		if err != nil {
			return fmt.Errorf("failed to register application crd, reason %v", err)
		}

		client, err := kubex.New(getter, log)
		if err != nil {
			return err
		}
		kc = client
		factory = client.Factory
	}

	lazyClient := &lazyClient{
		namespace: namespace,
		clientFn:  factory.KubernetesClientSet,
		appClientFn: func() (*appcs.Clientset, error) {
			config, err := factory.ToRawKubeConfigLoader().ClientConfig()
			if err != nil {
				return nil, err
			}
			return appcs.NewForConfig(config)
		},
	}

	var store *storage.Storage
	switch helmDriver {
	case "app", "apps", "application", "applications", "editor":
		config, err := factory.ToRawKubeConfigLoader().ClientConfig()
		if err != nil {
			return err
		}
		mapper, err := getter.ToRESTMapper()
		if err != nil {
			return err
		}
		d := driver2.NewApplications(
			newApplicationClient(lazyClient),
			dynamic.NewForConfigOrDie(config),
			disco_util.NewResourceMapper(mapper),
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
