package linkerd

import (
	"context"
	"fmt"

	"github.com/layer5io/meshery-adapter-library/adapter"
	"github.com/layer5io/meshery-adapter-library/common"
	adapterconfig "github.com/layer5io/meshery-adapter-library/config"
	"github.com/layer5io/meshery-adapter-library/status"
	internalconfig "github.com/layer5io/meshery-linkerd/internal/config"
	"github.com/layer5io/meshkit/logger"
)

// Linkerd is the handler for the adapter
type Linkerd struct {
	adapter.Adapter // Type Embedded
}

// New initializes linkerd handler.
func New(c adapterconfig.Handler, l logger.Handler, kc adapterconfig.Handler) adapter.Handler {
	return &Linkerd{
		Adapter: adapter.Adapter{
			Config:            c,
			Log:               l,
			KubeconfigHandler: kc,
		},
	}
}

// ApplyOperation applies the operation on linkerd
func (linkerd *Linkerd) ApplyOperation(ctx context.Context, opReq adapter.OperationRequest) error {
	operations := make(adapter.Operations)
	err := linkerd.Config.GetObject(adapter.OperationsKey, &operations)
	if err != nil {
		return err
	}

	e := &adapter.Event{
		Operationid: opReq.OperationID,
		Summary:     status.Deploying,
		Details:     "Operation is not supported",
	}

	switch opReq.OperationName {
	case internalconfig.LinkerdOperation:
		go func(hh *Linkerd, ee *adapter.Event) {
			version := string(operations[opReq.OperationName].Versions[0])
			stat, err := hh.installLinkerd(opReq.IsDeleteOperation, version, opReq.Namespace)
			if err != nil {
				e.Summary = fmt.Sprintf("Error while %s Linkerd service mesh", stat)
				e.Details = err.Error()
				hh.StreamErr(e, err)
				return
			}
			ee.Summary = fmt.Sprintf("Linkerd service mesh %s successfully", stat)
			ee.Details = fmt.Sprintf("The Linkerd service mesh is now %s.", stat)
			hh.StreamInfo(e)
		}(linkerd, e)
	case common.BookInfoOperation, common.HTTPBinOperation, common.ImageHubOperation, common.EmojiVotoOperation:
		go func(hh *Linkerd, ee *adapter.Event) {
			appName := operations[opReq.OperationName].AdditionalProperties[common.ServiceName]
			stat, err := hh.installSampleApp(opReq.Namespace, opReq.IsDeleteOperation, operations[opReq.OperationName].Templates)
			if err != nil {
				e.Summary = fmt.Sprintf("Error while %s %s application", stat, appName)
				e.Details = err.Error()
				hh.StreamErr(e, err)
				return
			}
			ee.Summary = fmt.Sprintf("%s application %s successfully", appName, stat)
			ee.Details = fmt.Sprintf("The %s application is now %s.", appName, stat)
			hh.StreamInfo(e)
		}(linkerd, e)
	case common.SmiConformanceOperation:
		go func(hh *Linkerd, ee *adapter.Event) {
			name := operations[opReq.OperationName].Description
			_, err := hh.RunSMITest(adapter.SMITestOptions{
				Ctx:         context.TODO(),
				OperationID: ee.Operationid,
				Namespace:   "meshery",
				Manifest:    string(operations[opReq.OperationName].Templates[0]),
				Labels:      make(map[string]string),
				Annotations: map[string]string{
					"linkerd.io/inject": "enabled",
				},
			})
			if err != nil {
				e.Summary = fmt.Sprintf("Error while %s %s test", status.Running, name)
				e.Details = err.Error()
				hh.StreamErr(e, err)
				return
			}
			ee.Summary = fmt.Sprintf("%s test %s successfully", name, status.Completed)
			ee.Details = ""
			hh.StreamInfo(e)
		}(linkerd, e)
	case common.CustomOperation:
		go func(hh *Linkerd, ee *adapter.Event) {
			stat, err := hh.applyCustomOperation(opReq.Namespace, opReq.CustomBody, opReq.IsDeleteOperation)
			if err != nil {
				e.Summary = fmt.Sprintf("Error while %s custom operation", stat)
				e.Details = err.Error()
				hh.StreamErr(e, err)
				return
			}
			ee.Summary = fmt.Sprintf("Manifest %s successfully", status.Deployed)
			ee.Details = ""
			hh.StreamInfo(e)
		}(linkerd, e)
	case internalconfig.AnnotateNamespace:
		go func(hh *Linkerd, ee *adapter.Event) {
			err := hh.LoadNamespaceToMesh(opReq.Namespace, opReq.IsDeleteOperation)
			if err != nil {
				e.Summary = fmt.Sprintf("Error while annotating %s", opReq.Namespace)
				e.Details = err.Error()
				hh.StreamErr(e, err)
				return
			}
			ee.Summary = "Annotation successful"
			ee.Details = ""
			hh.StreamInfo(e)
		}(linkerd, e)
	default:
		e.Summary = "Invalid Request"
		linkerd.StreamErr(e, ErrOpInvalid)
	}

	return nil
}
