package genericapiserver

import (
	"kis/apiserver/endpoints"
	"strings"
	"kis/apiserver/registry/rest"
	"net/http"
	"log"
	"io"
	"k8s.io/kubernetes/src/github.com/emicklei/go-restful"
)

type APIGroupInfo struct {
	VersionedResourcesStorageMap map[string]map[string] rest.RestStorage
}


// GenericAPIServer contains state for a Kubernetes cluster api server.
type GenericAPIServer struct {

	// "Outputs"
	// Handler holds the handlers being used by this API server
	Handler *APIServerHandler

	// SecureServingInfo holds configuration of the TLS server.
	SecureServingInfo *SecureServingInfo
}

// preparedGenericAPIServer is a private wrapper that enforces a call of PrepareRun() before Run can be invoked.
type preparedGenericAPIServer struct {
	*GenericAPIServer
}

func (s *GenericAPIServer) PrepareRun() *preparedGenericAPIServer{
	return &preparedGenericAPIServer{s}
}

func (s * preparedGenericAPIServer) Run(stopCh <-chan struct{}) error {
	return nil
}

func (s *GenericAPIServer) InstallLegacyAPIGroup(apiPrefix string, groupInfo *APIGroupInfo){
	s.installAPIResources(apiPrefix, groupInfo)
}

func(s *GenericAPIServer) installAPIResources(apiPrefix string, apiGroupInfo *APIGroupInfo){
	apiGroupVersion := s.getAPIGroupVersion(apiGroupInfo, "v1", apiPrefix)
	apiGroupVersion.InstallREST(s.Handler.GoRestfulContainer)
}

func (s *GenericAPIServer) getAPIGroupVersion(apiGroupInfo *APIGroupInfo, groupVersion string,apiPrefix string) *endpoints.APIGroupVersion{
	storages := make(map[string]rest.RestStorage)
	for k, v := range apiGroupInfo.VersionedResourcesStorageMap[groupVersion]{
		storages[strings.ToLower(k)] = v
	}
	version := s.newAPIGroupVersion(apiGroupInfo, groupVersion)
	version.Root = apiPrefix
	version.Storage = storages
	return version
}

func (s *GenericAPIServer) newAPIGroupVersion(apiGroupInfo *APIGroupInfo, groupVersion string) *endpoints.APIGroupVersion {
	return &endpoints.APIGroupVersion{
		GroupVersion: groupVersion,
	}
}

// serveSecurely runs the secure http server. It fails only if certificates cannot
// be loaded or the initial listen call fails. The actual server loop (stoppable by closing
// stopCh) runs in a go routine, i.e. serveSecurely does not block.
func (s *SecureServingInfo) Serve(handler http.Handler) error {
	secureServer := &http.Server{
		Addr: "127.0.0.1:8080",
		Handler:        handler,
	}
	log.Fatal(secureServer.ListenAndServe())
	return nil
}

func hello2(req *restful.Request, resp *restful.Response) {
	io.WriteString(resp, "default world")
}