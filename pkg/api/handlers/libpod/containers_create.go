package libpod

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/containers/podman/v4/libpod"
	"github.com/containers/podman/v4/pkg/api/handlers/utils"
	api "github.com/containers/podman/v4/pkg/api/types"
	"github.com/containers/podman/v4/pkg/domain/entities"
	"github.com/containers/podman/v4/pkg/specgen"
	"github.com/containers/podman/v4/pkg/specgen/generate"
	"github.com/pkg/errors"
)

// CreateContainer takes a specgenerator and makes a container. It returns
// the new container ID on success along with any warnings.
func CreateContainer(w http.ResponseWriter, r *http.Request) {
	runtime := r.Context().Value(api.RuntimeKey).(*libpod.Runtime)
	var sg specgen.SpecGenerator

	if err := json.NewDecoder(r.Body).Decode(&sg); err != nil {
		utils.Error(w, http.StatusInternalServerError, errors.Wrap(err, "Decode()"))
		return
	}
	if sg.Passwd == nil {
		t := true
		sg.Passwd = &t
	}
	warn, err := generate.CompleteSpec(r.Context(), runtime, &sg)
	if err != nil {
		utils.InternalServerError(w, err)
		return
	}

	for i, d := range sg.Devices {
		if len(d.Path) > 0 && d.Major == 0 && d.Minor == 0 {
			// source:destination:permissions
			tokens := strings.SplitN(d.Path, ":", 3)
			spec, err := generate.DeviceFromPath(tokens[0])
			if err != nil {
				utils.InternalServerError(w, errors.Wrapf(err, "failed to query linux device %q for container %q", d.Path, sg.Name))
			}
			sg.Devices[i].FileMode = spec.FileMode
			sg.Devices[i].GID = spec.GID
			sg.Devices[i].Major = spec.Major
			sg.Devices[i].Minor = spec.Minor
			sg.Devices[i].Type = spec.Type
			sg.Devices[i].UID = spec.UID

			switch len(tokens) {
			case 3:
				mode, err := generate.ParseFileMode(tokens[2])
				if err != nil {
					utils.InternalServerError(w, fmt.Errorf("invalid device permission specification %q for container %q", tokens[2], sg.Name))
				}
				sg.Devices[i].FileMode = &mode
				fallthrough
			case 2:
				sg.Devices[i].Path = tokens[1]
			case 1:
				sg.Devices[i].Path = tokens[0]
			default:
				utils.InternalServerError(w, fmt.Errorf("invalid device specification %q for container %q", d.Path, sg.Name))
				return
			}
		}
	}

	rtSpec, spec, opts, err := generate.MakeContainer(context.Background(), runtime, &sg, false, nil)
	if err != nil {
		utils.InternalServerError(w, err)
		return
	}
	ctr, err := generate.ExecuteCreate(context.Background(), runtime, rtSpec, spec, false, opts...)
	if err != nil {
		utils.InternalServerError(w, err)
		return
	}
	response := entities.ContainerCreateResponse{ID: ctr.ID(), Warnings: warn}
	utils.WriteJSON(w, http.StatusCreated, response)
}
