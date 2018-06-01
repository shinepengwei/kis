package storage

import (
	genericregistry "kis/apiserver/registry/generic"
	"kis/apiserver/runtime"
	apis "kis/apis/core"
)

type PodStorage struct {
	Pod         *REST
}

type REST struct{
	*genericregistry.Store
}

// NewStorage returns a RESTStorage object that will work against pods.
func NewStorage(optsGetter genericregistry.RESTOptionsGetter) PodStorage {
	store := &genericregistry.Store{
		NewFunc: 					func() runtime.Object { return &apis.Pod{} },
		NewListFunc:              	func() runtime.Object { return &apis.PodList{} },
		//KeyFunc:					func() string{return }
	}
	//初始化etcd连接，store.Storage
	store.CompleteWithOptions(optsGetter)

	return PodStorage{
		Pod:         &REST{store},
	}
}

