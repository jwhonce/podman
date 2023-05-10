package endpoints

import (
	"net/http"

	"github.com/emicklei/go-restful/v3"
)

type Service struct{}

func (s Service) RegisterTo(container *restful.Container) {
	systemService := new(restful.WebService)
	systemService.Path("/ping")

	systemService.Route(systemService.GET("/ping").To(Ping).
		Operation("getPing").
		Doc("Returns Service details").
		Writes("OK").
		Returns(http.StatusOK, "OK", nil))

	systemService.Route(systemService.HEAD("/ping").To(Ping).
		Operation("headPing").
		Doc("Returns Service details").
		Returns(http.StatusOK, "", nil),
	)
	container.Add(systemService)
}
