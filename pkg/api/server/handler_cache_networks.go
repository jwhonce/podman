package server

import (
	"net/http"
	"os"
	"sync"

	"github.com/containers/podman/v3/libpod"
	api "github.com/containers/podman/v3/pkg/api/types"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

var (
	previousNetworkUpdate os.FileInfo
	networkUpdateSync     sync.Once
)

// cacheNetworksHandler refreshes the in-memory network cache if the
// network configuration directory has been changed
func cacheNetworksHandler() mux.MiddlewareFunc {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			runtime := r.Context().Value(api.RuntimeKey).(*libpod.Runtime)
			cfg, err := runtime.GetConfig()
			if err != nil {
				logrus.Infof("Failed to retrieve podman configuration: %v", err)
			} else {
				netCfg := cfg.Network.NetworkConfigDir

				// First pass of this handler is for initialization.
				// 1+ passes will update memory network cache/configuration if the modification
				// time of the network configuration directory changes.
				networkUpdateSync.Do(func() {
					var err error
					previousNetworkUpdate, err = os.Stat(netCfg)
					logrus.Infof("Failed to access %q: %v", netCfg, err)
				})

				if now, err := os.Stat(netCfg); err != nil {
					logrus.Infof("Failed to access %q: %v", netCfg, err)
				} else {
					if now.ModTime().After(previousNetworkUpdate.ModTime()) {
						//
						// Reload network cache/configuration
						//
						previousNetworkUpdate = now
					}
				}
			}

			h.ServeHTTP(w, r)
		})
	}
}
