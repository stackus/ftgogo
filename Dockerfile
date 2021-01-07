ARG service
ARG cmd=service
FROM golang:alpine AS builder
ARG service
ARG cmd
# Install dependencies in a cachable way
RUN mkdir $GOPATH/src/shared-go
COPY ./shared-go/go.* $GOPATH/src/shared-go
WORKDIR $GOPATH/src/shared-go
RUN go mod download

RUN mkdir $GOPATH/src/serviceapis
COPY serviceapis/go.* $GOPATH/src/serviceapis
WORKDIR $GOPATH/src/serviceapis
RUN go mod download

RUN mkdir $GOPATH/src/${service}
COPY ./${service}/go.* $GOPATH/src/${service}
WORKDIR $GOPATH/src/${service}/cmd/${cmd}
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