ARG service
ARG cmd=service
FROM golang:1.16-alpine AS builder
ARG service
ARG cmd

WORKDIR $GOPATH/src/${service}/cmd/${cmd}
COPY . $GOPATH/src
RUN go install .

FROM alpine:latest AS runtime
ARG cmd
ENV cmd=$cmd
COPY --from=builder /go/bin/${cmd} /usr/local/bin/${cmd}

ENTRYPOINT /usr/local/bin/${cmd}
