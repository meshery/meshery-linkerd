package build

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/layer5io/meshery-adapter-library/adapter"

	"github.com/layer5io/meshkit/utils"
	"github.com/layer5io/meshkit/utils/manifests"
	"github.com/layer5io/meshkit/utils/walker"
	smp "github.com/layer5io/service-mesh-performance/spec"
)

var DefaultGenerationMethod string
var LatestVersion string
var MeshModelPath string
var AllVersions []string

const Component = "Linkerd"

var MeshModelConfig = adapter.MeshModelConfig{ //Move to build/config.go
	Category: "Cloud Native Network",
	Metadata: map[string]interface{}{},
}

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

var VersionToURL = make(map[string][]string)

func init() {
	// Initialize Metadata including logo svgs
	f, _ := os.Open("./build/meshmodel_metadata.json")
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Printf("Error closing file: %s\n", err)
		}
	}()
	byt, _ := io.ReadAll(f)

	_ = json.Unmarshal(byt, &MeshModelConfig.Metadata)
	wd, _ := os.Getwd()
	MeshModelPath = filepath.Join(wd, "templates", "meshmodel", "components")
	AllVersions, _ = utils.GetLatestReleaseTagsSorted("linkerd", "linkerd2")
	if len(AllVersions) == 0 {
		return
	}
	for _, v := range AllVersions {
		walker.NewGithub().Owner("linkerd").Repo("linkerd2").Branch(v).Root("charts/linkerd-crds/templates").RegisterFileInterceptor(func(gca walker.GithubContentAPI) error {
			VersionToURL[v] = append(VersionToURL[v], DefaultGenerationURL(v, gca.Name))
			return nil
		}).Walk()
	}
	LatestVersion = AllVersions[len(AllVersions)-1]
	DefaultGenerationMethod = adapter.Manifests

}

func DefaultGenerationURL(version string, crd string) string {
	return fmt.Sprintf("https://raw.githubusercontent.com/linkerd/linkerd2/%s/charts/linkerd-crds/templates/policy/%s", version, crd)
}
