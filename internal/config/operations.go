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
		Type:                 int32(meshes.OpCategory_INSTALL),
		Description:          "Linkerd Service Mesh",
		Versions:             versions,
		Templates:            []adapter.Template{},
		AdditionalProperties: map[string]string{},
	}

	dev[AnnotateNamespace] = &adapter.Operation{
		Type:        int32(meshes.OpCategory_CONFIGURE),
		Description: "Annotate Namespace",
	}
	dev[JaegerAddon] = &adapter.Operation{
		Type:        int32(meshes.OpCategory_CONFIGURE),
		Description: "Add-on: Jaeger",
		AdditionalProperties: map[string]string{
			ServiceName:      "jaeger",
			ServicePatchFile: "file://templates/oam/patches/service-loadbalancer.json",
			HelmChartURL:     "https://helm.linkerd.io/stable/linkerd-jaeger-2.10.2.tgz",
		},
	}
	dev[VizAddon] = &adapter.Operation{
		Type:        int32(meshes.OpCategory_CONFIGURE),
		Description: "Add-on: Viz",
		AdditionalProperties: map[string]string{
			ServiceName:      "web",
			ServicePatchFile: "file://templates/oam/patches/service-loadbalancer.json",
			HelmChartURL:     "https://helm.linkerd.io/stable/linkerd-viz-2.10.2.tgz",
		},
	}
	dev[MultiClusterAddon] = &adapter.Operation{
		Type:        int32(meshes.OpCategory_CONFIGURE),
		Description: "Add-on: Multi-cluster",
		AdditionalProperties: map[string]string{
			ServiceName:      "linkerd-gateway",
			ServicePatchFile: "file://templates/oam/patches/service-loadbalancer.json",
			HelmChartURL:     "https://helm.linkerd.io/stable/linkerd-multicluster-2.10.2.tgz",
		},
	}
	dev[SMIAddon] = &adapter.Operation{
		Type:        int32(meshes.OpCategory_CONFIGURE),
		Description: "Add-on: SMI Addon",
		AdditionalProperties: map[string]string{
			// ServiceName:      "linkerd-gateway",
			// ServicePatchFile: "file://templates/oam/patches/service-loadbalancer.json",
			HelmChartURL: "https://github.com/linkerd/linkerd-smi/releases/download/v0.1.0/linkerd-smi-0.1.0.tgz",
		},
	}
	return dev
}
