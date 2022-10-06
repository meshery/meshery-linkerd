package linkerd

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/layer5io/meshery-adapter-library/common"
	"github.com/layer5io/meshery-adapter-library/meshes"
	"github.com/layer5io/meshery-linkerd/internal/config"
	"github.com/layer5io/meshkit/models/oam/core/v1alpha1"
	"gopkg.in/yaml.v3"
)

// CompHandler is the type for functions which can handle OAM components
type CompHandler func(*Linkerd, v1alpha1.Component, bool, []string) (string, error)

// HandleComponents handles the processing of OAM components
func (linkerd *Linkerd) HandleComponents(comps []v1alpha1.Component, isDel bool, kubeconfigs []string) (string, error) {
	var errs []error
	var msgs []string
	stat1 := "deploying"
	stat2 := "deployed"
	if isDel {
		stat1 = "removing"
		stat2 = "removed"
	}
	compFuncMap := map[string]CompHandler{
		"LinkerdMesh":              handleComponentLinkerdMesh,
		"JaegerLinkerdAddon":       handleComponentLinkerdAddon,
		"VizLinkerdAddon":          handleComponentLinkerdAddon,
		"MultiClusterLinkerdAddon": handleComponentLinkerdAddon,
		"SMILinkerdAddon":          handleComponentLinkerdAddon,
	}

	for _, comp := range comps {
		ee := &meshes.EventsResponse{
			OperationId:   uuid.New().String(),
			Component:     config.ServerConfig["type"],
			ComponentName: config.ServerConfig["name"],
		}
		fnc, ok := compFuncMap[comp.Spec.Type]
		if !ok {
			msg, err := handleLinkerdCoreComponent(linkerd, comp, isDel, "", "", kubeconfigs)
			if err != nil {
				ee.Summary = fmt.Sprintf("Error while %s %s", stat1, comp.Spec.Type)
				linkerd.streamErr(ee.Summary, ee, err)
				errs = append(errs, err)
				continue
			}
			ee.Summary = fmt.Sprintf("%s %s successfully", comp.Spec.Type, stat2)
			ee.Details = fmt.Sprintf("The %s is now %s.", comp.Spec.Type, stat2)
			linkerd.StreamInfo(ee)
			msgs = append(msgs, msg)
			continue
		}

		msg, err := fnc(linkerd, comp, isDel, kubeconfigs)
		if err != nil {
			ee.Summary = fmt.Sprintf("Error while %s %s", stat1, comp.Spec.Type)
			linkerd.streamErr(ee.Summary, ee, err)
			errs = append(errs, err)
			continue
		}
		ee.Summary = fmt.Sprintf("%s %s %s successfully", comp.Name, comp.Spec.Type, stat2)
		ee.Details = fmt.Sprintf("The %s %s is now %s.", comp.Name, comp.Spec.Type, stat2)
		linkerd.StreamInfo(ee)
		msgs = append(msgs, msg)
	}

	if err := mergeErrors(errs); err != nil {
		return mergeMsgs(msgs), err
	}

	return mergeMsgs(msgs), nil
}

// HandleApplicationConfiguration handles the processing of OAM application configuration
func (linkerd *Linkerd) HandleApplicationConfiguration(config v1alpha1.Configuration, isDel bool, kubeconfigs []string) (string, error) {
	var errs []error
	var msgs []string
	for _, comp := range config.Spec.Components {
		for _, trait := range comp.Traits {
			if trait.Name == "automaticSidecarInjection.Linkerd" {
				namespaces := castSliceInterfaceToSliceString(trait.Properties["namespaces"].([]interface{}))
				if err := handleNamespaceLabel(linkerd, namespaces, isDel, kubeconfigs); err != nil {
					errs = append(errs, err)
				}
			}

			msgs = append(msgs, fmt.Sprintf("applied trait \"%s\" on service \"%s\"", trait.Name, comp.ComponentName))
		}
	}

	if err := mergeErrors(errs); err != nil {
		return mergeMsgs(msgs), err
	}

	return mergeMsgs(msgs), nil
}

func handleNamespaceLabel(linkerd *Linkerd, namespaces []string, isDel bool, kubeconfigs []string) error {
	var errs []error
	for _, ns := range namespaces {
		if err := linkerd.AnnotateNamespace(ns, isDel, map[string]string{
			"linkerd.io/inject": "enabled",
		}, kubeconfigs); err != nil {
			errs = append(errs, err)
		}
	}

	return mergeErrors(errs)
}

func handleComponentLinkerdMesh(linkerd *Linkerd, comp v1alpha1.Component, isDel bool, kubeconfigs []string) (string, error) {
	version := comp.Spec.Version
	return linkerd.installLinkerd(isDel, version, comp.Namespace, kubeconfigs)
}

func handleLinkerdCoreComponent(
	linkerd *Linkerd,
	comp v1alpha1.Component,
	isDel bool,
	apiVersion,
	kind string,
	kubeconfigs []string) (string, error) {
	if apiVersion == "" {
		apiVersion = getAPIVersionFromComponent(comp)
		if apiVersion == "" {
			return "", ErrLinkerdCoreComponentFail(fmt.Errorf("failed to get API Version for: %s", comp.Name))
		}
	}

	if kind == "" {
		kind = getKindFromComponent(comp)
		if kind == "" {
			return "", ErrLinkerdCoreComponentFail(fmt.Errorf("failed to get kind for: %s", comp.Name))
		}
	}

	component := map[string]interface{}{
		"apiVersion": apiVersion,
		"kind":       kind,
		"metadata": map[string]interface{}{
			"name":        comp.Name,
			"annotations": comp.Annotations,
			"labels":      comp.Labels,
		},
		"spec": comp.Spec.Settings,
	}

	// Convert to yaml
	yamlByt, err := yaml.Marshal(component)
	if err != nil {
		err = ErrParseLinkerdCoreComponent(err)
		linkerd.Log.Error(err)
		return "", err
	}

	msg := fmt.Sprintf("created %s \"%s\" in namespace \"%s\"", kind, comp.Name, comp.Namespace)
	if isDel {
		msg = fmt.Sprintf("deleted %s config \"%s\" in namespace \"%s\"", kind, comp.Name, comp.Namespace)
	}

	return msg, linkerd.applyManifest(yamlByt, isDel, comp.Namespace, kubeconfigs)
}

func handleComponentLinkerdAddon(istio *Linkerd, comp v1alpha1.Component, isDel bool, kubeconfigs []string) (string, error) {
	var addonName string
	var helmURL string
	version := removePrefixFromVersionIfPresent(comp.Spec.Version)
	switch comp.Spec.Type {
	case "JaegerLinkerdAddon":
		addonName = config.JaegerAddon
		helmURL = "https://helm.linkerd.io/stable/linkerd-jaeger-" + version + ".tgz"
	case "VizLinkerdAddon":
		addonName = config.VizAddon
		helmURL = "https://helm.linkerd.io/stable/linkerd-viz-" + version + ".tgz"
	case "MultiClusterLinkerdAddon":
		addonName = config.MultiClusterAddon
		helmURL = "https://helm.linkerd.io/stable/linkerd-multicluster-" + version + ".tgz"
	case "SMIClusterLinkerdAddon":
		addonName = config.SMIAddon
		helmURL = "https://github.com/linkerd/linkerd-smi/releases/download/v0.1.0/linkerd-smi-0.1.0.tgz"
	default:
		return "", nil
	}

	// Get the service
	svc := config.Operations[addonName].AdditionalProperties[common.ServiceName]

	// Get the patches
	patches := make([]string, 0)
	patches = append(patches, config.Operations[addonName].AdditionalProperties[config.ServicePatchFile])

	_, err := istio.installAddon(comp.Namespace, isDel, svc, patches, helmURL, addonName, kubeconfigs)
	msg := fmt.Sprintf("created service of type \"%s\"", comp.Spec.Type)
	if isDel {
		msg = fmt.Sprintf("deleted service of type \"%s\"", comp.Spec.Type)
	}

	return msg, err
}
func getAPIVersionFromComponent(comp v1alpha1.Component) string {
	return comp.Annotations["pattern.meshery.io.mesh.workload.k8sAPIVersion"]
}

func getKindFromComponent(comp v1alpha1.Component) string {
	return comp.Annotations["pattern.meshery.io.mesh.workload.k8sKind"]
}

func castSliceInterfaceToSliceString(in []interface{}) []string {
	var out []string

	for _, v := range in {
		cast, ok := v.(string)
		if ok {
			out = append(out, cast)
		}
	}

	return out
}

func mergeErrors(errs []error) error {
	if len(errs) == 0 {
		return nil
	}

	var errMsgs []string

	for _, err := range errs {
		errMsgs = append(errMsgs, err.Error())
	}

	return fmt.Errorf(strings.Join(errMsgs, "\n"))
}

func mergeMsgs(strs []string) string {
	return strings.Join(strs, "\n")
}
func removePrefixFromVersionIfPresent(version string) string {
	if version == "" {
		return "2.10.1" //default, to avoid any errors
	}
	if strings.HasPrefix(version, "stable-") {
		return strings.TrimPrefix(version, "stable-")
	}
	if strings.HasPrefix(version, "edge-") {
		return strings.TrimPrefix(version, "edge-")
	}
	return version
}
