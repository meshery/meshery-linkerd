module github.com/layer5io/meshery-linkerd

go 1.15

replace github.com/kudobuilder/kuttl => github.com/layer5io/kuttl v0.4.1-0.20200723152044-916f10574334

require (
	github.com/layer5io/meshery-adapter-library v0.1.24
	github.com/layer5io/meshkit v0.2.29
	github.com/layer5io/service-mesh-performance v0.3.3
	gopkg.in/yaml.v2 v2.4.0
	helm.sh/helm/v3 v3.6.3 // indirect
	k8s.io/apimachinery v0.21.0
)
