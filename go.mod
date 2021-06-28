module github.com/layer5io/meshery-linkerd

go 1.14

replace github.com/kudobuilder/kuttl => github.com/layer5io/kuttl v0.4.1-0.20200723152044-916f10574334

require (
	github.com/golang/protobuf v1.4.3 // indirect
	github.com/layer5io/meshery-adapter-library v0.1.20
	github.com/layer5io/meshkit v0.2.14
	github.com/layer5io/service-mesh-performance v0.3.3
	k8s.io/apimachinery v0.18.12
)
