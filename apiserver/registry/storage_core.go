package registry

import (
	"kis/apiserver/storage"
	"kis/apiserver/genericapiserver"
	podstore "kis/apiserver/registry/core/pod/storage"
	"kis/apiserver/registry/generic"
	"kis/apiserver/registry/rest"
)
type LegacyRESTStorageProvider struct {
	StorageFactory storage.StorageFactory
}

func(c LegacyRESTStorageProvider) NewLegacyRESTStorage(restOptionsGetter generic.RESTOptionsGetter) genericapiserver.APIGroupInfo {
	apiGroupInfo := genericapiserver.APIGroupInfo{
		VersionedResourcesStorageMap: map[string]map[string]rest.RestStorage{},
	}
	podStorage := podstore.NewStorage(restOptionsGetter)

	restStorageMap := map[string]rest.RestStorage{
		"pods":             podStorage.Pod,
	}
	apiGroupInfo.VersionedResourcesStorageMap["v1"] = restStorageMap
	return apiGroupInfo
}