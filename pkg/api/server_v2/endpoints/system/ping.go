package endpoints

import (
	"net/http"

	"github.com/containers/buildah"
	"github.com/emicklei/go-restful/v3"
)

func Ping(request *restful.Request, response *restful.Response) {
	// TODO: Add filter to set API version headers
	// Note: API-Version and Libpod-API-Version are set in handler_api.go
	response.AddHeader("BuildKit-Version", "")
	response.AddHeader("Builder-Version", "")
	response.AddHeader("Docker-Experimental", "true")
	response.AddHeader("Cache-Control", "no-cache")
	response.AddHeader("Pragma", "no-cache")
	response.AddHeader("Libpod-Buildah-Version", buildah.Version)

	if request.Request.Method == http.MethodGet {
		response.WriteEntity("OK")
		return
	}
	response.WriteHeader(http.StatusOK)
}
