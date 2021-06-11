module github.com/stackus/ftgogo/order

go 1.16

replace github.com/stackus/ftgogo/serviceapis => ./../serviceapis

replace shared-go => ../shared-go

require (
	github.com/mattn/go-colorable v0.1.8 // indirect
	github.com/prometheus/client_golang v1.11.0
	github.com/prometheus/common v0.29.0 // indirect
	github.com/rs/zerolog v1.22.0
	github.com/stackus/edat v0.0.3
	github.com/stackus/edat-pgx v0.0.2
	github.com/stackus/errors v0.0.3
	github.com/stackus/ftgogo/serviceapis v0.0.0-20210116185538-3dd9fbb69179
	go.uber.org/atomic v1.8.0 // indirect
	golang.org/x/net v0.0.0-20210610132358-84b48f89b13b // indirect
	google.golang.org/genproto v0.0.0-20210610141715-e7a9b787a5a4 // indirect
	google.golang.org/grpc v1.38.0
	shared-go v0.0.0-00010101000000-000000000000
)
