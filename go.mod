module github.com/layer5io/meshery-linkerd

go 1.14

require (
	github.com/Azure/go-autorest/autorest/adal v0.6.0 // indirect
	github.com/alecthomas/template v0.0.0-20190718012654-fb15b899a751
	github.com/fatih/color v1.9.0
	github.com/ghodss/yaml v1.0.0
	github.com/golang/protobuf v1.4.2
	github.com/googleapis/gnostic v0.3.1 // indirect
	github.com/gophercloud/gophercloud v0.4.0 // indirect
	github.com/hashicorp/go-getter v1.4.1
	github.com/linkerd/linkerd2 v0.0.0-00010101000000-000000000000
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.6.0
	github.com/stretchr/testify v1.5.1
	golang.org/x/net v0.0.0-20200202094626-16171245cfb2
	google.golang.org/grpc v1.29.1
	k8s.io/apiextensions-apiserver v0.17.4
	k8s.io/apimachinery v0.17.4
	k8s.io/client-go v0.17.4
	k8s.io/kube-aggregator v0.17.4
)

replace github.com/linkerd/linkerd2 => github.com/Aisuko/linkerd2 v0.0.0-20200813135257-77538bd6c81d

//replace github.com/linkerd/linkerd2 => /Users/aisuko/Documents/linkerd2
