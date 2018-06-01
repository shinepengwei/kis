package master

import (
	"kis/apiserver/storage"
	"kis/apiserver/registry"
	"kis/apiserver/genericapiserver"
	"kis/apiserver/registry/generic"
)
type ExtraConfig struct{
	StorageFactory	storage.StorageFactory
}

type Config struct {
	GenericConfig *genericapiserver.Config
	ExtraConfig   ExtraConfig
}

func(c Config) New()(*Master){
	s := c.GenericConfig.NewGenericAPIServer("kis-apiserver")
	m := &Master{
		GenericAPIServer: s,
	}
	legacyRESTStorageProvider := registry.LegacyRESTStorageProvider{
		StorageFactory:             c.ExtraConfig.StorageFactory,
	}
	m.InstallLegacyAPI(&c, c.GenericConfig.RESTOptionsGetter, legacyRESTStorageProvider)

	return m
}

func (m *Master) InstallLegacyAPI(c *Config, restOptionGetter generic.RESTOptionsGetter,legacyRESTStorageProvider registry.LegacyRESTStorageProvider){
	apiGroupInfo := legacyRESTStorageProvider.NewLegacyRESTStorage(restOptionGetter)
	m.GenericAPIServer.InstallLegacyAPIGroup("/",&apiGroupInfo)
}


type Master struct {
	GenericAPIServer *genericapiserver.GenericAPIServer
}

func (m *Master) Run(stopCh <-chan struct{}){
	m.GenericAPIServer.SecureServingInfo.Serve(m.GenericAPIServer.Handler.GoRestfulContainer)
	<-stopCh
}