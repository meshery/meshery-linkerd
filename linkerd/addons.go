package linkerd

import (
	"context"
	"net/url"

	"github.com/layer5io/meshery-adapter-library/status"
	"github.com/layer5io/meshery-linkerd/internal/config"
	"github.com/layer5io/meshkit/utils"
	"github.com/layer5io/meshkit/utils/kubernetes"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

// installAddon installs/uninstalls an addon in the given namespace
func (linkerd *Linkerd) installAddon(namespace string, del bool, service string, patches []string, helmChartURL string, addon string) (string, error) {
	act := kubernetes.INSTALL
	st := status.Installing

	if del {
		st = status.Removing
		act = kubernetes.UNINSTALL
	}
	var err error
	switch addon {
	case config.JaegerAddon:
		err = linkerd.MesheryKubeclient.ApplyHelmChart(kubernetes.ApplyHelmChartConfig{
			URL:             helmChartURL,
			Namespace:       namespace,
			CreateNamespace: true,
			Action:          act,
			OverrideValues: map[string]interface{}{
				"installNamespace": false, //Set to false when installing in a custom namespace.
				"namespace":        namespace,
			},
		})
	case config.VizAddon:
		err = linkerd.MesheryKubeclient.ApplyHelmChart(kubernetes.ApplyHelmChartConfig{
			URL:             helmChartURL,
			Namespace:       namespace,
			CreateNamespace: true,
			Action:          act,
			OverrideValues: map[string]interface{}{
				"installNamespace": false, //Set to false when installing in a custom namespace.
				"linkerdNamespace": linkerdNamespace,
				"namespace":        namespace,
			},
		})
	case config.MultiClusterAddon:
		err = linkerd.MesheryKubeclient.ApplyHelmChart(kubernetes.ApplyHelmChartConfig{
			URL:             helmChartURL,
			Namespace:       namespace,
			CreateNamespace: true,
			Action:          act,
			OverrideValues: map[string]interface{}{
				"installNamespace": false, //Set to false when installing in a custom namespace.
				"linkerdNamespace": linkerdNamespace,
				"namespace":        namespace,
			},
		})
	case config.SMIAddon:
		err = linkerd.MesheryKubeclient.ApplyHelmChart(kubernetes.ApplyHelmChartConfig{
			URL:             helmChartURL,
			Namespace:       namespace,
			Action:          act,
			CreateNamespace: true,
			OverrideValues: map[string]interface{}{
				"installNamespace": false, //Set to false when installing in a custom namespace.
				"namespace":        namespace,
			},
		})
	}

	if err != nil {
		return st, ErrAddonFromHelm(err)
	}

	for _, patch := range patches {
		if !del {
			_, err := url.ParseRequestURI(patch)
			if err != nil {
				return st, ErrAddonFromHelm(err)
			}

			content, err := utils.ReadFileSource(patch)
			if err != nil {
				return st, ErrAddonFromHelm(err)
			}

			_, err = linkerd.KubeClient.CoreV1().Services(namespace).Patch(context.TODO(), service, types.MergePatchType, []byte(content), metav1.PatchOptions{})
			if err != nil {
				return st, ErrAddonFromHelm(err)
			}
		}
	}
	return st, nil
}
