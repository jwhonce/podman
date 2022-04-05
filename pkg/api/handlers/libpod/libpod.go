package libpod

import (
	_ "github.com/containers/podman/v4/pkg/errorhandling" // nolint:golint  imported into pkg for side effect
	jsoniter "github.com/json-iterator/go"
)

// pull frozen jsoniter into package
var json = jsoniter.ConfigCompatibleWithStandardLibrary
