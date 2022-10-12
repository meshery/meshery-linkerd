package linkerd

import (
	"context"
	"fmt"
	"sync"

	"github.com/layer5io/meshery-adapter-library/adapter"
	"github.com/layer5io/meshery-adapter-library/common"
	adapterconfig "github.com/layer5io/meshery-adapter-library/config"
	"github.com/layer5io/meshery-adapter-library/meshes"
	"github.com/layer5io/meshery-adapter-library/status"
	internalconfig "github.com/layer5io/meshery-linkerd/internal/config"
	"github.com/layer5io/meshery-linkerd/linkerd/oam"
	"github.com/layer5io/meshkit/errors"
	"github.com/layer5io/meshkit/logger"
	"github.com/layer5io/meshkit/models"
	"github.com/layer5io/meshkit/models/oam/core/v1alpha1"
	"github.com/layer5io/meshkit/utils/events"
	mesherykube "github.com/layer5io/meshkit/utils/kubernetes"
	"gopkg.in/yaml.v3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Linkerd is the handler for the adapter
type Linkerd struct {
	adapter.Adapter // Type Embedded
}

// New initializes linkerd handler.
func New(c adapterconfig.Handler, l logger.Handler, kc adapterconfig.Handler, ev *events.EventStreamer) adapter.Handler {
	return &Linkerd{
		Adapter: adapter.Adapter{
			Config:            c,
			Log:               l,
			KubeconfigHandler: kc,
			EventStreamer:     ev,
		},
	}
}

// CreateKubeconfigs creates and writes passed kubeconfig onto the filesystem
func (linkerd *Linkerd) CreateKubeconfigs(kubeconfigs []string) error {
	var errs = make([]error, 0)
	for _, kubeconfig := range kubeconfigs {
		kconfig := models.Kubeconfig{}
		err := yaml.Unmarshal([]byte(kubeconfig), &kconfig)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		// To have control over what exactly to take in on kubeconfig
		linkerd.KubeconfigHandler.SetKey("kind", kconfig.Kind)
		linkerd.KubeconfigHandler.SetKey("apiVersion", kconfig.APIVersion)
		linkerd.KubeconfigHandler.SetKey("current-context", kconfig.CurrentContext)
		err = linkerd.KubeconfigHandler.SetObject("preferences", kconfig.Preferences)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		err = linkerd.KubeconfigHandler.SetObject("clusters", kconfig.Clusters)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		err = linkerd.KubeconfigHandler.SetObject("users", kconfig.Users)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		err = linkerd.KubeconfigHandler.SetObject("contexts", kconfig.Contexts)
		if err != nil {
			errs = append(errs, err)
			continue
		}
	}
	if len(errs) == 0 {
		return nil
	}
	return mergeErrors(errs)
}

// ApplyOperation applies the operation on linkerd
func (linkerd *Linkerd) ApplyOperation(ctx context.Context, opReq adapter.OperationRequest) error {
	err := linkerd.CreateKubeconfigs(opReq.K8sConfigs)
	if err != nil {
		return err
	}
	operations := make(adapter.Operations)
	kubeConfigs := opReq.K8sConfigs
	err = linkerd.Config.GetObject(adapter.OperationsKey, &operations)
	if err != nil {
		return err
	}

	e := &meshes.EventsResponse{
		OperationId:   opReq.OperationID,
		Summary:       status.Deploying,
		Details:       "Operation is not supported",
		Component:     internalconfig.ServerConfig["type"],
		ComponentName: internalconfig.ServerConfig["name"],
	}

	switch opReq.OperationName {
	case internalconfig.LinkerdOperation:
		go func(hh *Linkerd, ee *meshes.EventsResponse) {
			var err error
			var stat, version string
			if len(operations[opReq.OperationName].Versions) == 0 {
				err = ErrFetchLinkerdVersions
			} else {
				version = string(operations[opReq.OperationName].Versions[len(operations[opReq.OperationName].Versions)-1])
				stat, err = hh.installLinkerd(opReq.IsDeleteOperation, version, opReq.Namespace, kubeConfigs)
			}
			if err != nil {
				summary := fmt.Sprintf("Error while %s Linkerd service mesh", stat)
				hh.streamErr(summary, ee, err)
				return
			}
			ee.Summary = fmt.Sprintf("Linkerd service mesh %s successfully", stat)
			ee.Details = fmt.Sprintf("The Linkerd service mesh is now %s.", stat)
			hh.StreamInfo(ee)
		}(linkerd, e)
	case common.BookInfoOperation, common.HTTPBinOperation, common.ImageHubOperation, common.EmojiVotoOperation:
		go func(hh *Linkerd, ee *meshes.EventsResponse) {
			appName := operations[opReq.OperationName].AdditionalProperties[common.ServiceName]
			stat, err := hh.installSampleApp(opReq.Namespace, opReq.IsDeleteOperation, operations[opReq.OperationName].Templates, kubeConfigs)
			if err != nil {
				summary := fmt.Sprintf("Error while %s %s application", stat, appName)
				hh.streamErr(summary, ee, err)
				return
			}
			ee.Summary = fmt.Sprintf("%s application %s successfully", appName, stat)
			ee.Details = fmt.Sprintf("The %s application is now %s.", appName, stat)
			hh.StreamInfo(ee)
		}(linkerd, e)
	case common.SmiConformanceOperation:
		go func(hh *Linkerd, ee *meshes.EventsResponse) {
			name := operations[opReq.OperationName].Description
			_, err := hh.RunSMITest(adapter.SMITestOptions{
				Ctx:         context.TODO(),
				OperationID: ee.OperationId,
				Namespace:   "meshery",
				Manifest:    string(operations[opReq.OperationName].Templates[0]),
				Labels:      make(map[string]string),
				Annotations: map[string]string{
					"linkerd.io/inject": "enabled",
				},
			})
			if err != nil {
				summary := fmt.Sprintf("Error while %s %s test", status.Running, name)
				hh.streamErr(summary, ee, err)
				return
			}
			ee.Summary = fmt.Sprintf("%s test %s successfully", name, status.Completed)
			ee.Details = ""
			hh.StreamInfo(ee)
		}(linkerd, e)
	case common.CustomOperation:
		go func(hh *Linkerd, ee *meshes.EventsResponse) {
			stat, err := hh.applyCustomOperation(opReq.Namespace, opReq.CustomBody, opReq.IsDeleteOperation, kubeConfigs)
			if err != nil {
				summary := fmt.Sprintf("Error while %s custom operation", stat)
				hh.streamErr(summary, ee, err)
				return
			}
			ee.Summary = fmt.Sprintf("Manifest %s successfully", status.Deployed)
			ee.Details = ""
			hh.StreamInfo(ee)
		}(linkerd, e)
	case internalconfig.JaegerAddon, internalconfig.VizAddon, internalconfig.MultiClusterAddon, internalconfig.SMIAddon:
		go func(hh *Linkerd, ee *meshes.EventsResponse) {
			svcname := operations[opReq.OperationName].AdditionalProperties[common.ServiceName]
			patches := make([]string, 0)
			patches = append(patches, operations[opReq.OperationName].AdditionalProperties[internalconfig.ServicePatchFile])
			helmChartURL := operations[opReq.OperationName].AdditionalProperties[internalconfig.HelmChartURL]
			_, err := hh.installAddon(opReq.Namespace, opReq.IsDeleteOperation, svcname, patches, helmChartURL, opReq.OperationName, kubeConfigs)
			operation := "install"
			if opReq.IsDeleteOperation {
				operation = "uninstall"
			}

			if err != nil {
				summary := fmt.Sprintf("Error while %sing %s", operation, opReq.OperationName)
				hh.streamErr(summary, ee, err)
				return
			}
			ee.Summary = fmt.Sprintf("Successfully %sed %s", operation, opReq.OperationName)
			ee.Details = fmt.Sprintf("Successfully %sed %s from the %s namespace", operation, opReq.OperationName, opReq.Namespace)
			hh.StreamInfo(ee)
		}(linkerd, e)
	case internalconfig.AnnotateNamespace:
		go func(hh *Linkerd, ee *meshes.EventsResponse) {
			err := hh.AnnotateNamespace(opReq.Namespace, opReq.IsDeleteOperation, map[string]string{
				"linkerd.io/inject": "enabled",
			}, kubeConfigs)
			if err != nil {
				summary := fmt.Sprintf("Error while annotating %s", opReq.Namespace)
				hh.streamErr(summary, ee, err)
				return
			}
			ee.Summary = "Annotation successful"
			ee.Details = ""
			hh.StreamInfo(ee)
		}(linkerd, e)
	default:
		summary := "Invalid Request"
		linkerd.streamErr(summary, e, ErrOpInvalid)
	}

	return nil
}

// ProcessOAM will handles the grpc invocation for handling OAM objects
func (linkerd *Linkerd) ProcessOAM(ctx context.Context, oamReq adapter.OAMRequest) (string, error) {
	err := linkerd.CreateKubeconfigs(oamReq.K8sConfigs)
	if err != nil {
		return "", err
	}
	kubeconfigs := oamReq.K8sConfigs
	var comps []v1alpha1.Component
	for _, acomp := range oamReq.OamComps {
		comp, err := oam.ParseApplicationComponent(acomp)
		if err != nil {
			linkerd.Log.Error(ErrParseOAMComponent)
			continue
		}

		comps = append(comps, comp)
	}

	config, err := oam.ParseApplicationConfiguration(oamReq.OamConfig)
	if err != nil {
		linkerd.Log.Error(ErrParseOAMConfig)
	}

	// If operation is delete then first HandleConfiguration and then handle the deployment
	if oamReq.DeleteOp {
		// Process configuration
		msg2, err := linkerd.HandleApplicationConfiguration(config, oamReq.DeleteOp, kubeconfigs)
		if err != nil {
			return msg2, ErrProcessOAM(err)
		}

		// Process components
		msg1, err := linkerd.HandleComponents(comps, oamReq.DeleteOp, kubeconfigs)
		if err != nil {
			return msg1 + "\n" + msg2, ErrProcessOAM(err)
		}

		return msg1 + "\n" + msg2, nil
	}

	// Process components
	msg1, err := linkerd.HandleComponents(comps, oamReq.DeleteOp, kubeconfigs)
	if err != nil {
		return msg1, ErrProcessOAM(err)
	}

	// Process configuration
	msg2, err := linkerd.HandleApplicationConfiguration(config, oamReq.DeleteOp, kubeconfigs)
	if err != nil {
		return msg1 + "\n" + msg2, ErrProcessOAM(err)
	}

	return msg1 + "\n" + msg2, nil
}

// AnnotateNamespace is used to label namespaces ,for cases like automatic sidecar injection (or not). If the namespace is not present, it will create one, instead of throwing error.
func (linkerd *Linkerd) AnnotateNamespace(namespace string, remove bool, labels map[string]string, kubeconfigs []string) error {
	var errs []error
	var errMx sync.Mutex
	var wg sync.WaitGroup
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
			ns, err := kClient.KubeClient.CoreV1().Namespaces().Get(context.TODO(), namespace, metav1.GetOptions{})
			if err != nil {
				linkerd.Log.Info("Namespace \"", namespace, "\" not present. Creating namespace")
				var er error
				ns, er = createNS(kClient, namespace)
				if er != nil {
					errMx.Lock()
					errs = append(errs, err)
					errMx.Unlock()
					return
				}
			}

			if ns.ObjectMeta.Annotations == nil {
				ns.ObjectMeta.Annotations = map[string]string{}
			}
			for key, val := range labels {
				ns.ObjectMeta.Annotations[key] = val
			}

			if remove {
				for key := range labels {
					delete(ns.ObjectMeta.Annotations, key)
				}
			}

			_, err = kClient.KubeClient.CoreV1().Namespaces().Update(context.TODO(), ns, metav1.UpdateOptions{})
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
		return ErrAnnotatingNamespace(mergeErrors(errs))
	}
	return nil
}

func (linkerd *Linkerd) streamErr(summary string, e *meshes.EventsResponse, err error) {
	e.Summary = summary
	e.Details = err.Error()
	e.ErrorCode = errors.GetCode(err)
	e.ProbableCause = errors.GetCause(err)
	e.SuggestedRemediation = errors.GetRemedy(err)
	linkerd.StreamErr(e, err)
}
