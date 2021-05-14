package main

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
	"kmodules.xyz/client-go/discovery"
)

func IsResourceReady(mapper discovery.ResourceMapper, gvr schema.GroupVersionResource) (bool, error) {
	return false, nil
}
