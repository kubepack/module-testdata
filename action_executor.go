package main

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"kmodules.xyz/client-go/discovery"
)

type ActionRunner struct {
	dc dynamic.Interface
	mapper discovery.ResourceMapper

	action Action
	err error
}

/*
	runner := new(ActionRunner)
	if runner.MeetsPrerequisites() {
		runner.Apply().WaitUntilReady()
		if runner.Err() != nil {
			// if ae, ok := runner.Err().(AlreadyErrored); ok {
				// action was already errored out, Rest before reusing
			}
		}

	} else {
		if runner.Err() != nil {
			// if ae, ok := runner.Err().(AlreadyErrored); ok {
				// action was already errored out, Rest before reusing
			}
		}
	}

*/

// Do we need this?
func (e *ActionRunner) ResetError() {
	e.err = nil
}

func (e *ActionRunner) Err() error {
	return e.err
}

func (e *ActionRunner) MeetsPrerequisites() bool {
	if e.err != nil {
		e.err = NewAlreadyErrored(e.err)
		return false
	}

	for _, r := range e.action.Prerequisites.RequiredResources {
		exists, err := IsResourceExistsAndReady(e.dc, e.mapper, schema.GroupVersionResource{
			Group:    r.Group,
			Version:  r.Version,
			Resource: r.Resource,
		})
		if err != nil {
			e.err = err
			return false
		}
		if !exists {
			return false
		}
	}
	return true
}

func (e *ActionRunner) Apply() *ActionRunner {
	if e.err != nil {
		e.err = NewAlreadyErrored(e.err)
		return e
	}

	return e
}

func (e *ActionRunner) WaitUntilReady() {
	if e.err != nil {
		e.err = NewAlreadyErrored(e.err)
		return
	}




}

func (e *ActionRunner) IsReady() bool {
	if e.err != nil {
		e.err = NewAlreadyErrored(e.err)
		return false
	}

	return false
}

type AlreadyErrored struct {
	underlying error
}

func NewAlreadyErrored(underlying error) error {
	if _, ok := underlying.(AlreadyErrored); ok {
		return underlying
	}
	return AlreadyErrored{underlying: underlying}
}

func (e AlreadyErrored) Error() string {
	return e.underlying.Error()
}

func (e AlreadyErrored) Underlying() error {
	return e.underlying
}
