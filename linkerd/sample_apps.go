package linkerd

import (
	"context"
	"fmt"

	"github.com/layer5io/meshery-adapter-library/adapter"
	"github.com/layer5io/meshery-adapter-library/status"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (linkerd *Linkerd) installSampleApp(namespace string, del bool, templates []adapter.Template) (string, error) {
	linkerd.Log.Info(fmt.Sprintf("Requested action is delete: %v", del))
	st := status.Installing

	if del {
		st = status.Removing
	}

	for _, template := range templates {
		err := linkerd.applyManifest([]byte(template.String()), del, namespace)
		if err != nil {
			return st, ErrSampleApp(err)
		}
	}

	return status.Installed, nil
}

// LoadToMesh adds annotation to service
func (linkerd *Linkerd) LoadToMesh(namespace string, service string, remove bool) error {
	deploy, err := linkerd.KubeClient.AppsV1().Deployments(namespace).Get(context.TODO(), service, metav1.GetOptions{})
	if err != nil {
		return err
	}

	if deploy.ObjectMeta.Annotations == nil {
		deploy.ObjectMeta.Annotations = map[string]string{}
	}
	deploy.ObjectMeta.Annotations["linkerd.io/inject"] = "enabled"

	if remove {
		delete(deploy.ObjectMeta.Annotations, "linkerd.io/inject")
	}

	_, err = linkerd.KubeClient.AppsV1().Deployments(namespace).Update(context.TODO(), deploy, metav1.UpdateOptions{})
	if err != nil {
		return err
	}

	return nil
}

// LoadNamespaceToMesh is used to mark namespaces for automatic sidecar injection (or not)
func (linkerd *Linkerd) LoadNamespaceToMesh(namespace string, remove bool) error {
	ns, err := linkerd.KubeClient.CoreV1().Namespaces().Get(context.TODO(), namespace, metav1.GetOptions{})
	if err != nil {
		return err
	}

	if ns.ObjectMeta.Annotations == nil {
		ns.ObjectMeta.Annotations = map[string]string{}
	}
	ns.ObjectMeta.Annotations["linkerd.io/inject"] = "enabled"

	if remove {
		delete(ns.ObjectMeta.Annotations, "linkerd.io/inject")
	}

	_, err = linkerd.KubeClient.CoreV1().Namespaces().Update(context.TODO(), ns, metav1.UpdateOptions{})
	if err != nil {
		return err
	}
	return nil
}
