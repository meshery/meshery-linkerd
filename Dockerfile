FROM golang:1.15 as build-env
WORKDIR /github.com/layer5io/meshery-linkerd
COPY go.mod go.sum ./
RUN go mod download

COPY main.go main.go
COPY internal/ internal/
COPY linkerd/ linkerd/

RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -ldflags="-w -s" -a -o meshery-linkerd main.go

FROM gcr.io/distroless/base:nonroot-amd64
ENV DISTRO="debian"
ENV GOARCH="amd64"
WORKDIR /$HOME/.meshery
COPY --from=build-env /github.com/layer5io/meshery-linkerd/meshery-linkerd .
ENTRYPOINT ["./meshery-linkerd"]
