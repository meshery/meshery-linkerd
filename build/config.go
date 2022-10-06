package build

import (
	"fmt"
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
var CRDnamesURL map[string]string

const Component = "Linkerd"

// NewConfig creates the configuration for creating components
func NewConfig(version string) manifests.Config {
	return manifests.Config{
		Name:        smp.ServiceMesh_Type_name[int32(smp.ServiceMesh_LINKERD)],
		Type:        Component,
		MeshVersion: version,
		CrdFilter: manifests.NewCueCrdFilter(manifests.ExtractorPaths{
			NamePath:    "spec.names.kind",
			IdPath:      "spec.names.kind",
			VersionPath: "spec.versions[0].name",
			GroupPath:   "spec.group",
			SpecPath:    "spec.versions[0].schema.openAPIV3Schema.properties.spec"}, false),
		ExtractCrds: func(manifest string) []string {
			manifests.RemoveHelmTemplatingFromCRD(&manifest)
			crds := strings.Split(manifest, "---")
			return crds
		},
	}
}
func init() {
	wd, _ := os.Getwd()
	WorkloadPath = filepath.Join(wd, "templates", "oam", "workloads")
	vs, err := config.GetLatestReleaseNames(30)
	if len(vs) == 0 {
		fmt.Println("dynamic component generation failure: ", err.Error())
		return
	}
	for _, v := range vs {
		AllVersions = append(AllVersions, string(v))
	}
	LatestVersion = AllVersions[0]
	DefaultGenerationMethod = adapter.Manifests
	names, err := config.GetFileNames("linkerd", "linkerd2", "charts/linkerd-crds/templates/**")
	if err != nil {
		fmt.Println("dynamic component generation failure: ", err.Error())
		return
	}
	for n := range names {
		if !strings.HasSuffix(n, ".yaml") {
			delete(names, n)
		}
	}
	CRDnamesURL = names
}
