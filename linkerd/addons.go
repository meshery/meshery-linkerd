package linkerd

import (
	"context"
	"net/url"
	"sync"

	"github.com/layer5io/meshery-adapter-library/status"
	"github.com/layer5io/meshery-linkerd/internal/config"
	"github.com/layer5io/meshkit/utils"
	"github.com/layer5io/meshkit/utils/kubernetes"
	mesherykube "github.com/layer5io/meshkit/utils/kubernetes"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

// installAddon installs/uninstalls an addon in the given namespace
func (linkerd *Linkerd) installAddon(namespace string, del bool, service string, patches []string, helmChartURL string, addon string, kubeconfigs []string) (string, error) {
	act := kubernetes.INSTALL
	st := status.Installing

	if del {
		st = status.Removing
		act = kubernetes.UNINSTALL
	}
	var errs []error
	var wg sync.WaitGroup
	for _, k8sconfig := range kubeconfigs {
		wg.Add(1)
		go func(k8sconfig string) {
			defer wg.Done()
			kClient, err := mesherykube.New([]byte(k8sconfig))
			if err != nil {
				errs = append(errs, err)
				return
			}
			switch addon {
			case config.JaegerAddon:
				err = kClient.ApplyHelmChart(kubernetes.ApplyHelmChartConfig{
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
				err = kClient.ApplyHelmChart(kubernetes.ApplyHelmChartConfig{
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
				err = kClient.ApplyHelmChart(kubernetes.ApplyHelmChartConfig{
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
				err = kClient.ApplyHelmChart(kubernetes.ApplyHelmChartConfig{
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
				errs = append(errs, err)
				return
			}

			for _, patch := range patches {
				if !del {
					_, err := url.ParseRequestURI(patch)
					if err != nil {
						errs = append(errs, err)
						return
					}

					content, err := utils.ReadFileSource(patch)
					if err != nil {
						errs = append(errs, err)
						return
					}

					_, err = kClient.KubeClient.CoreV1().Services(namespace).Patch(context.TODO(), service, types.MergePatchType, []byte(content), metav1.PatchOptions{})
					if err != nil {
						errs = append(errs, err)
						return
					}
				}
			}
		}(k8sconfig)
	}
	wg.Wait()
	if len(errs) != 0 {
		return st, ErrAddonFromHelm(mergeErrors(errs))
	}
	return st, nil
}
