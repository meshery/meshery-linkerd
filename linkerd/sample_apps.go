package linkerd

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/layer5io/meshery-adapter-library/adapter"
	"github.com/layer5io/meshery-adapter-library/status"
	"github.com/layer5io/meshkit/utils"
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

	data, err := ioutil.ReadFile(location)
	if err != nil {
		return "", err
	}

	return string(data), nil
}
