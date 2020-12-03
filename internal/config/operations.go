package config

import (
	"github.com/layer5io/meshery-adapter-library/adapter"
	"github.com/layer5io/meshery-adapter-library/meshes"
)

var (
	ServiceName = "service_name"
)

func getOperations(dev adapter.Operations) adapter.Operations {
	versions, _ := getLatestReleaseNames(3)

	dev[LinkerdOperation] = &adapter.Operation{
		Type:        int32(meshes.OpCategory_INSTALL),
		Description: "Linkerd Service Mesh",
		Versions:    versions,
		Templates: []adapter.Template{
			"templates/linkerd.yaml",
		},
		AdditionalProperties: map[string]string{},
	}

	dev[AnnotateNamespace] = &adapter.Operation{
		Type:        int32(meshes.OpCategory_CONFIGURE),
		Description: "Annotate Namespace",
	}

	return dev
}
