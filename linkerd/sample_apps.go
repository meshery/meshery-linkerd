package linkerd

import (
	"github.com/layer5io/meshery-adapter-library/adapter"
	"github.com/layer5io/meshery-adapter-library/status"
	"github.com/layer5io/meshkit/utils"
)

func (linkerd *Linkerd) installSampleApp(del bool, templates []adapter.Template) (string, error) {
	st := status.Installing

	if del {
		st = status.Removing
	}

	for _, template := range templates {
		contents, err := utils.ReadRemoteFile(string(template))
		if err != nil {
			return st, ErrSampleApp(err)
		}

		err = linkerd.applyManifest([]byte(contents))
		if err != nil {
			return st, ErrSampleApp(err)
		}
	}

	return status.Installed, nil
}
