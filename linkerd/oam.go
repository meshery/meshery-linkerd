package linkerd

import (
	"fmt"
	"strings"

	"github.com/layer5io/meshery-adapter-library/common"
	"github.com/layer5io/meshery-linkerd/internal/config"
	"github.com/layer5io/meshkit/models/oam/core/v1alpha1"
	"gopkg.in/yaml.v2"
)

// CompHandler is the type for functions which can handle OAM components
type CompHandler func(*Linkerd, v1alpha1.Component, bool) (string, error)

// HandleComponents handles the processing of OAM components
func (linkerd *Linkerd) HandleComponents(comps []v1alpha1.Component, isDel bool) (string, error) {
	var errs []error
	var msgs []string

	compFuncMap := map[string]CompHandler{
		"LinkerdMesh":              handleComponentLinkerdMesh,
		"JaegerLinkerdAddon":       handleComponentLinkerdAddon,
		"VizLinkerdAddon":          handleComponentLinkerdAddon,
		"MultiClusterLinkerdAddon": handleComponentLinkerdAddon,
		"SMILinkerdAddon":          handleComponentLinkerdAddon,
	}

	for _, comp := range comps {
		fnc, ok := compFuncMap[comp.Spec.Type]
		if !ok {
			msg, err := handleLinkerdCoreComponent(linkerd, comp, isDel, "", "")
			if err != nil {
				errs = append(errs, err)
				continue
			}

			msgs = append(msgs, msg)
			continue
		}

		msg, err := fnc(linkerd, comp, isDel)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		msgs = append(msgs, msg)
	}

	if err := mergeErrors(errs); err != nil {
		return mergeMsgs(msgs), err
	}

	return mergeMsgs(msgs), nil
}

// HandleApplicationConfiguration handles the processing of OAM application configuration
func (linkerd *Linkerd) HandleApplicationConfiguration(config v1alpha1.Configuration, isDel bool) (string, error) {
	var errs []error
	var msgs []string
	for _, comp := range config.Spec.Components {
		for _, trait := range comp.Traits {
			if trait.Name == "automaticSidecarInjection.Linkerd" {
				namespaces := castSliceInterfaceToSliceString(trait.Properties["namespaces"].([]interface{}))
				if err := handleNamespaceLabel(linkerd, namespaces, isDel); err != nil {
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

func handleNamespaceLabel(linkerd *Linkerd, namespaces []string, isDel bool) error {
	var errs []error
	for _, ns := range namespaces {
		if err := linkerd.AnnotateNamespace(ns, isDel, map[string]string{
			"linkerd.io/inject": "enabled",
		}); err != nil {
			errs = append(errs, err)
		}
	}

	return mergeErrors(errs)
}

func handleComponentLinkerdMesh(linkerd *Linkerd, comp v1alpha1.Component, isDel bool) (string, error) {
	// Get the linkerd version from the settings
	// we are sure that the version of linkerd would be present
	// because the configuration is already validated against the schema
	version := comp.Spec.Settings["version"].(string)

	return linkerd.installLinkerd(isDel, version, comp.Namespace)
}

func handleLinkerdCoreComponent(
	linkerd *Linkerd,
	comp v1alpha1.Component,
	isDel bool,
	apiVersion,
	kind string) (string, error) {
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

	return msg, linkerd.applyManifest(yamlByt, isDel, comp.Namespace)
}

func handleComponentLinkerdAddon(istio *Linkerd, comp v1alpha1.Component, isDel bool) (string, error) {
	var addonName string
	var helmURL string
	version := removePrefixFromVersionIfPresent(comp.Spec.Settings["version"].(string))
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

	_, err := istio.installAddon(comp.Namespace, isDel, svc, patches, helmURL, addonName)
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
