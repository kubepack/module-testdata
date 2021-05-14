package kube

import (
	"io"
	"time"

	"helm.sh/helm/v3/pkg/kube"
	v1 "k8s.io/api/core/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

type Client struct {
	*kube.Client
}

var _ kube.Interface = &Client{}

func New(getter genericclioptions.RESTClientGetter, log func(string, ...interface{})) *Client {
	kc := kube.New(getter)
	kc.Log = log
	return &Client{
		Client: kc,
	}
}

func (c *Client) Create(resources kube.ResourceList) (*kube.Result, error) {
	return c.Client.Create(resources)
}

func (c *Client) Wait(resources kube.ResourceList, timeout time.Duration) error {
	return c.Client.Wait(resources, timeout)
}

func (c *Client) WaitWithJobs(resources kube.ResourceList, timeout time.Duration) error {
	return c.Client.WaitWithJobs(resources, timeout)
}

func (c *Client) Delete(resources kube.ResourceList) (*kube.Result, []error) {
	return c.Client.Delete(resources)
}

func (c *Client) WatchUntilReady(resources kube.ResourceList, timeout time.Duration) error {
	return c.Client.WatchUntilReady(resources, timeout)
}

func (c *Client) Update(original, target kube.ResourceList, force bool) (*kube.Result, error) {
	return c.Client.Update(original, target, force)
}

func (c *Client) Build(reader io.Reader, validate bool) (kube.ResourceList, error) {
	return c.Client.Build(reader, validate)
}

func (c *Client) WaitAndGetCompletedPodPhase(name string, timeout time.Duration) (v1.PodPhase, error) {
	return c.Client.WaitAndGetCompletedPodPhase(name, timeout)
}

func (c *Client) IsReachable() error {
	return c.Client.IsReachable()
}
