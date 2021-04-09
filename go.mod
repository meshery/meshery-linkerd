module github.com/layer5io/meshery-linkerd

go 1.14

replace (
	github.com/docker/docker => github.com/moby/moby v17.12.0-ce-rc1.0.20200618181300-9dc6525e6118+incompatible
	github.com/kudobuilder/kuttl => github.com/layer5io/kuttl v0.4.1-0.20200723152044-916f10574334
	github.com/layer5io/meshery-adapter-library v0.1.12 => /Users/abishekk/Documents/layer5/meshery-adapter-library
	golang.org/x/sys => golang.org/x/sys v0.0.0-20200826173525-f9321e4c35a6
)

require (
	github.com/golang/protobuf v1.4.3 // indirect
	github.com/layer5io/meshery-adapter-library v0.1.14
	github.com/layer5io/meshkit v0.2.7
	github.com/layer5io/service-mesh-performance v0.3.3
	k8s.io/apimachinery v0.18.12
)
