protoc-setup:
	cd meshes
	wget https://raw.githubusercontent.com/layer5io/meshery/master/meshes/meshops.proto

proto:	
	protoc -I meshes/ meshes/meshops.proto --go_out=plugins=grpc:./meshes/

docker:
	docker build -t layer5/meshery-linkerd .

docker-run:
	(docker rm -f meshery-linkerd) || true
	docker run --name meshery-linkerd -d \
	-p 10001:10001 \
	-e DEBUG=true \
	layer5/meshery-linkerd

run:
	DEBUG=true go run main.go