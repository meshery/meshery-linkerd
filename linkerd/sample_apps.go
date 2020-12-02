package linkerd

import (
	"context"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/layer5io/meshery-adapter-library/adapter"
	"github.com/layer5io/meshery-adapter-library/status"
	"github.com/layer5io/meshkit/utils"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (linkerd *Linkerd) installSampleApp(namespace string, del bool, templates []adapter.Template) (string, error) {
	linkerd.Log.Info(fmt.Sprintf("Requested action is delete: %v", del))
	st := status.Installing

	if del {
		st = status.Removing
	}

	for _, template := range templates {
		contents, err := readFileSource(string(template))
		if err != nil {
			return st, ErrSampleApp(err)
		}

		err = linkerd.applyManifest([]byte(contents), del, namespace)
		if err != nil {
			return st, ErrSampleApp(err)
		}
	}

	return status.Installed, nil
}

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

// AddAnnotation is used to mark namespaces for automatic sidecar injection (or not)
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

// readFileSource supports "http", "https" and "file" protocols.
// it takes in the location as a uri and returns the contents of
// file as a string.
//
// TODO: May move this function to meshkit
func readFileSource(uri string) (string, error) {
	if strings.HasPrefix(uri, "http") {
		return utils.ReadRemoteFile(uri)
	}
	if strings.HasPrefix(uri, "file") {
		return readLocalFile(uri)
	}

	return "", fmt.Errorf("invalid protocol: only http, https and file are valid protocols")
}

// readLocalFile takes in the location of a local file
// in the format `file://location/of/file` and returns
// the content of the file if the path is valid and no
// error occurs
func readLocalFile(location string) (string, error) {
	// remove the protocol prefix
	location = strings.TrimPrefix(location, "file://")

	// Need to support variable file locations hence
	// #nosec
	data, err := ioutil.ReadFile(location)
	if err != nil {
		return "", err
	}

	return string(data), nil
}
