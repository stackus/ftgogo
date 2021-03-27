ARG service
ARG cmd=service
FROM golang:1.16-alpine AS builder
ARG service
ARG cmd
# Install dependencies in a cachable way
WORKDIR $GOPATH/src/shared-go
COPY ./shared-go/go.* $GOPATH/src/shared-go
RUN go mod download

WORKDIR $GOPATH/src/serviceapis
COPY serviceapis/go.* $GOPATH/src/serviceapis
RUN go mod download

WORKDIR $GOPATH/src/${service}/cmd/${cmd}
COPY ./${service}/go.* $GOPATH/src/${service}
RUN go mod download

COPY serviceapis $GOPATH/src/serviceapis
COPY ./shared-go $GOPATH/src/shared-go
COPY ./${service} $GOPATH/src/${service}

RUN go install .

FROM alpine:latest AS runtime
ARG cmd
ENV cmd=$cmd
COPY --from=builder /go/bin/${cmd} /usr/local/bin/${cmd}

ENTRYPOINT /usr/local/bin/${cmd}