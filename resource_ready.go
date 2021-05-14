package main

import (
	"context"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"kmodules.xyz/client-go/discovery"
	"sigs.k8s.io/cli-utils/pkg/kstatus/status"
)

var grCRD = schema.GroupResource{
	Group:    "apiextensions.k8s.io",
	Resource: "customresourcedefinitions",
}

func IsResourceExistsAndReady(dc dynamic.Interface, mapper discovery.ResourceMapper, gvr schema.GroupVersionResource) (bool, error) {
	exists, err := mapper.ExistsGVR(gvr)
	if gvr.GroupResource() != grCRD {
		return exists, err
	}

	crdGVR, err := mapper.Preferred(gvr)
	if err != nil {
		return false, err
	}
	obj, err := dc.Resource(crdGVR).Get(context.TODO(), gvr.GroupResource().String(), metav1.GetOptions{})
	if err != nil {
		return false, err
	}
	s, err := status.Compute(obj)
	if err != nil {
		return false, err
	}
	if s.Status == status.CurrentStatus {
		return true, nil
	}
	return false, fmt.Errorf("%+v", s)
}
