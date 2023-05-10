package server

import (
	"net"
	"time"

	"github.com/containers/podman/v4/libpod"
	images "github.com/containers/podman/v4/pkg/api/server_v2/endpoints/images"
	system "github.com/containers/podman/v4/pkg/api/server_v2/endpoints/system"
	"github.com/containers/podman/v4/pkg/domain/entities"
	"github.com/emicklei/go-restful/v3"
)

type APIServer struct {
	*restful.Container
}

const (
	DefaultCorsHeaders     = ""                // default Cross-Origin Resource Sharing (CORS) headers
	DefaultServiceDuration = 300 * time.Second // Number of seconds to wait for next request, if exceeded shutdown server
)

func NewServer(runtime *libpod.Runtime) (*APIServer, error) {
	return newServer(runtime, nil, entities.ServiceOptions{
		CorsHeaders: DefaultCorsHeaders,
		Timeout:     DefaultServiceDuration,
	})
}

func NewServerWithSettings(runtime *libpod.Runtime, listener net.Listener, opts entities.ServiceOptions) (*APIServer, error) {
	return newServer(runtime, listener, opts)
}

func newServer(runtime *libpod.Runtime, listener net.Listener, opts entities.ServiceOptions) (*APIServer, error) {
	server := &APIServer{
		Container: restful.DefaultContainer,
	}
	// TODO: restful.CrossOriginResourceSharing setup...

	system.Service{}.RegisterTo(server.Container)
	images.Service{}.RegisterTo(server.Container)
	return server, nil
}
