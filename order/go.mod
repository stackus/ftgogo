module github.com/stackus/ftgogo/order

go 1.16

replace github.com/stackus/ftgogo/serviceapis => ./../serviceapis

replace shared-go => ../shared-go

// Development replacements
//replace github.com/stackus/edat => ../../edat
//replace github.com/stackus/edat-msgpack => ../../edat-msgpack
//replace github.com/stackus/edat-pgx => ../../edat-pgx

require (
	github.com/cucumber/godog v0.11.0
	github.com/getkin/kin-openapi v0.68.0 // indirect
	github.com/google/uuid v1.3.0
	github.com/hashicorp/go-immutable-radix v1.3.1 // indirect
	github.com/hashicorp/go-memdb v1.3.2 // indirect
	github.com/jackc/pgx/v4 v4.13.0 // indirect
	github.com/mattn/go-colorable v0.1.8 // indirect
	github.com/prometheus/client_golang v1.11.0
	github.com/prometheus/common v0.30.0 // indirect
	github.com/rdumont/assistdog v0.0.0-20201106100018-168b06230d14
	github.com/spf13/pflag v1.0.5
	github.com/stackus/edat v0.0.6
	github.com/stackus/edat-msgpack v0.0.2
	github.com/stackus/edat-pgx v0.0.2
	github.com/stackus/errors v0.0.3
	github.com/stackus/ftgogo/serviceapis v0.0.0-20210116185538-3dd9fbb69179
	golang.org/x/net v0.0.0-20210726213435-c6fcb2dbf985 // indirect
	google.golang.org/genproto v0.0.0-20210729151513-df9385d47c1b // indirect
	google.golang.org/grpc v1.39.0
	shared-go v0.0.0-00010101000000-000000000000
)
