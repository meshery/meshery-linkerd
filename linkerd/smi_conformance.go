package linkerd

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/layer5io/gokit/smi"
	"github.com/layer5io/meshery-linkerd/meshes"
)

func (iClient *Client) runConformanceTest(id string, name string, version string) error {

	annotations := map[string]string{
		"linkerd.io/inject": "enabled",
	}

	test, err := smi.New(context.TODO(), id, version, name, iClient.k8sClientset)
	if err != nil {
		iClient.eventChan <- &meshes.EventsResponse{
			OperationId: id,
			EventType:   meshes.EventType_ERROR,
			Summary:     "Error while creating smi-conformance tool",
			Details:     err.Error(),
		}
		return err
	}

	result, err := test.Run(nil, annotations)
	if err != nil {
		iClient.eventChan <- &meshes.EventsResponse{
			OperationId: id,
			EventType:   meshes.EventType_ERROR,
			Summary:     fmt.Sprintf("Error while %s running smi-conformance test", result.Status),
			Details:     err.Error(),
		}
		return err
	}

	jsondata, _ := json.Marshal(result)
	iClient.eventChan <- &meshes.EventsResponse{
		OperationId: id,
		EventType:   meshes.EventType_INFO,
		Summary:     fmt.Sprintf("Smi conformance test %s successfully", result.Status),
		Details:     string(jsondata),
	}

	return nil
}
