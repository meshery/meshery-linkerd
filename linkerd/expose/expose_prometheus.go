package expose

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// Prometheus exposes the "linkerd-prometheus" deployment
func Prometheus(clientSet *kubernetes.Clientset, restConfig rest.Config, logger Logger) error {
	return Expose(clientSet, restConfig, Config{
		Name:      "linkerd-prometheus-meshery",
		Namespace: "linkerd",
		Type:      "NodePort",
		Logger:    logger,
	}, []string{"linkerd deployment linkerd-prometheus"})
}
