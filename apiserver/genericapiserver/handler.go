package genericapiserver

import (
	"github.com/emicklei/go-restful"
)

// APIServerHandlers holds the different http.Handlers used by the API server.
// This includes the full handler chain, the director (which chooses between gorestful and nonGoRestful,
// the gorestful handler (used for the API) which falls through to the nonGoRestful handler on unregistered paths,
// and the nonGoRestful handler (which can contain a fallthrough of its own)
// FullHandlerChain -> Director -> {GoRestfulContainer,NonGoRestfulMux} based on inspection of registered web services
type APIServerHandler struct {
	// The registered APIs.  InstallAPIs uses this.  Other servers probably shouldn't access this directly.
	GoRestfulContainer *restful.Container
}

func NewAPIServerHandler(name string) *APIServerHandler {
	gorestfulContainer := restful.NewContainer()
	//gorestfulContainer.ServeMux = http.NewServeMux()
	//gorestfulContainer.Router(restful.CurlyRouter{})
	return &APIServerHandler{
		GoRestfulContainer: gorestfulContainer,
	}
}
