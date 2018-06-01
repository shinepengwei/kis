package genericapiserver

import (
	"net"
	"kis/apiserver/registry/generic"
)

type SecureServingInfo struct {
	// Listener is the secure server network listener.
	Listener net.Listener
}

// Config is a structure used to configure a GenericAPIServer.
// Its members are sorted roughly in order of importance for composers.
type Config struct {
	// SecureServing is required to serve https
	SecureServing *SecureServingInfo

	// RESTOptionsGetter is used to construct RESTStorage types via the generic registry.
	RESTOptionsGetter generic.RESTOptionsGetter
}


func NewGenericAPIServerConfig() *Config{
	ln,_ := net.Listen("tcp", "127.0.0.1:8081")
	return &Config{
		SecureServing:&SecureServingInfo{
			Listener:ln,
		},
	}
}

func (s *Config) NewGenericAPIServer(name string)(*GenericAPIServer){
	apiServerHandler := NewAPIServerHandler("test")
	return &GenericAPIServer{
		Handler:apiServerHandler,
		SecureServingInfo:s.SecureServing,
	}
}
