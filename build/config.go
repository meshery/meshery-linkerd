package build

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/layer5io/meshery-adapter-library/adapter"
	"github.com/layer5io/meshery-linkerd/internal/config"

	"github.com/layer5io/meshkit/utils/manifests"
	smp "github.com/layer5io/service-mesh-performance/spec"
)

var DefaultGenerationMethod string
var LatestVersion string
var WorkloadPath string
var AllVersions []string
var CRDnames []string

const Component = "Linkerd"

//NewConfig creates the configuration for creating components
func NewConfig(version string) manifests.Config {
	return manifests.Config{
		Name:        smp.ServiceMesh_Type_name[int32(smp.ServiceMesh_LINKERD)],
		Type:        Component,
		MeshVersion: version,
		Filter: manifests.CrdFilter{
			RootFilter:    []string{"$[?(@.kind==\"CustomResourceDefinition\")]"},
			NameFilter:    []string{"$..[\"spec\"][\"names\"][\"kind\"]"},
			VersionFilter: []string{"$[0]..spec.versions[0]"},
			GroupFilter:   []string{"$[0]..spec"},
			SpecFilter:    []string{"$[0]..openAPIV3Schema.properties.spec"},
			ItrFilter:     []string{"$[?(@.spec.names.kind"},
			ItrSpecFilter: []string{"$[?(@.spec.names.kind"},
			VField:        "name",
			GField:        "group",
		},
	}
}
func init() {
	wd, _ := os.Getwd()
	WorkloadPath = filepath.Join(wd, "templates", "oam", "workloads")
	vs, _ := config.GetLatestReleaseNames(30)
	if len(vs) == 0 {
		return
	}
	for _, v := range vs {
		AllVersions = append(AllVersions, string(v))
	}
	LatestVersion = AllVersions[0]
	DefaultGenerationMethod = adapter.Manifests
	names, err := config.GetFileNames("linkerd", "linkerd2", "charts/linkerd-crds/templates/**")
	if err != nil {
		return
	}
	for _, n := range names {
		if strings.HasSuffix(n, ".yaml") {
			CRDnames = append(CRDnames, n)
		}
	}
}

func GenerationURL(crd string) string {
	return "https://raw.githubusercontent.com/linkerd/linkerd2/main/charts/linkerd-crds/templates/" + crd
}
