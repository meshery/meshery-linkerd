package linkerd

import (
	"context"
	"fmt"
	"sync"

	"github.com/layer5io/meshery-adapter-library/adapter"
	"github.com/layer5io/meshery-adapter-library/status"
	mesherykube "github.com/layer5io/meshkit/utils/kubernetes"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (linkerd *Linkerd) installSampleApp(namespace string, del bool, templates []adapter.Template, kubeconfigs []string) (string, error) {
	linkerd.Log.Info(fmt.Sprintf("Requested action is delete: %v", del))
	st := status.Installing

	if del {
		st = status.Removing
	}

	for _, template := range templates {
		err := linkerd.applyManifest([]byte(template.String()), del, namespace, kubeconfigs)
		if err != nil {
			return st, ErrSampleApp(err)
		}
	}

	return status.Installed, nil
}

// LoadToMesh adds annotation to service
func (linkerd *Linkerd) LoadToMesh(namespace string, service string, remove bool, kubeconfigs []string) error {
	var errs []error
	var wg sync.WaitGroup
	var errMx sync.Mutex
	for _, k8sconfig := range kubeconfigs {
		wg.Add(1)
		go func(k8sconfig string) {
			defer wg.Done()
			kClient, err := mesherykube.New([]byte(k8sconfig))
			if err != nil {
				errMx.Lock()
				errs = append(errs, err)
				errMx.Unlock()
				return
			}
			deploy, err := kClient.KubeClient.AppsV1().Deployments(namespace).Get(context.TODO(), service, metav1.GetOptions{})
			if err != nil {
				errMx.Lock()
				errs = append(errs, err)
				errMx.Unlock()
				return
			}

			if deploy.ObjectMeta.Annotations == nil {
				deploy.ObjectMeta.Annotations = map[string]string{}
			}
			deploy.ObjectMeta.Annotations["linkerd.io/inject"] = "enabled"

			if remove {
				delete(deploy.ObjectMeta.Annotations, "linkerd.io/inject")
			}

			_, err = kClient.KubeClient.AppsV1().Deployments(namespace).Update(context.TODO(), deploy, metav1.UpdateOptions{})
			if err != nil {
				errMx.Lock()
				errs = append(errs, err)
				errMx.Unlock()
				return
			}
		}(k8sconfig)
	}
	wg.Wait()
	if len(errs) != 0 {
		return mergeErrors(errs)
	}
	return nil
}
