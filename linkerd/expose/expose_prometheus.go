package expose

import (
	"github.com/layer5io/meshkit/logger"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// Prometheus exposes the "linkerd-prometheus" deployment
func Prometheus(clientSet *kubernetes.Clientset, restConfig rest.Config, logger logger.Handler, del bool) error {
	if !del {
		return Expose(clientSet, restConfig, Config{
			Name:      "linkerd-prometheus-meshery",
			Namespace: "linkerd",
			Type:      "NodePort",
			Logger:    logger,
		}, []string{"linkerd deployment linkerd-prometheus"})
	}

	return Remove("linkerd-prometheus-meshery", "linkerd", clientSet)
}
