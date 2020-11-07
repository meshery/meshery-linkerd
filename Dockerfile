FROM golang:1.13 as bd
WORKDIR /github.com/layer5io/meshery-linkerd
ADD . .
RUN GOPROXY=direct GOSUMDB=off go build -ldflags="-w -s" -a -o /meshery-linkerd .
RUN find . -name "*.go" -type f -delete; mv linkerd /

FROM alpine
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
RUN apk --update add ca-certificates curl && \
    mkdir /lib64 && \
    ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2
# Install kubectl
RUN curl -LO "https://storage.googleapis.com/kubernetes-release/release/v1.18.0/bin/linux/amd64/kubectl" && \
	chmod +x ./kubectl && \
	mv ./kubectl /usr/local/bin/kubectl

USER appuser
RUN mkdir -p /home/appuser/.kube
WORKDIR /home/appuser
COPY --from=bd /meshery-linkerd /home/appuser
COPY --from=bd /linkerd /home/appuser/linkerd
CMD ./meshery-linkerd
