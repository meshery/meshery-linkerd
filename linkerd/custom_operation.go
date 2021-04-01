// Package linkerd provides custom operation ability for deploying Linkerd2
package linkerd

import (
	"github.com/layer5io/meshery-adapter-library/status"
)

// applyCustomOperation return deployed status or the error message if it exist in deployment process
func (linkerd *Linkerd) applyCustomOperation(namespace string, manifest string, isDel bool) (string, error) {
	st := status.Starting

	err := linkerd.applyManifest([]byte(manifest), isDel, namespace)
	if err != nil {
		return st, ErrCustomOperation(err)
	}

	return status.Completed, nil
}
