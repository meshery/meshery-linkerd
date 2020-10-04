FROM golang:1.13 as bd
RUN adduser --disabled-login appuser
WORKDIR /github.com/layer5io/meshery-linkerd
ADD . .
RUN GOPROXY=direct GOSUMDB=off go build -ldflags="-w -s" -a -o /meshery-linkerd .
RUN find . -name "*.go" -type f -delete; mv linkerd /

FROM alpine
RUN apk --update add ca-certificates curl
RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2
COPY --from=bd /meshery-linkerd /app/
COPY --from=bd /linkerd /app/linkerd
COPY --from=bd /etc/passwd /etc/passwd
# Install kubectl
RUN curl -LO "https://storage.googleapis.com/kubernetes-release/release/v1.18.0/bin/linux/amd64/kubectl" && \
	chmod +x ./kubectl && \
	mv ./kubectl /usr/local/bin/kubectl
USER appuser
WORKDIR /app
CMD ./meshery-linkerd
