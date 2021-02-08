module github.com/layer5io/meshery-linkerd

go 1.14

replace github.com/kudobuilder/kuttl => github.com/layer5io/kuttl v0.4.1-0.20200723152044-916f10574334

require (
	github.com/golang/protobuf v1.4.3 // indirect
	github.com/layer5io/meshery-adapter-library v0.1.12-0.20210127214045-50f4c3bbd783
	github.com/layer5io/meshkit v0.2.1-0.20210127211805-88e99ca45457
	github.com/layer5io/service-mesh-performance v0.3.3
	google.golang.org/grpc v1.33.1 // indirect
	k8s.io/apimachinery v0.18.12
)
