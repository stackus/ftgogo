module github.com/stackus/ftgogo/order-history

go 1.16

replace github.com/stackus/ftgogo/serviceapis => ./../serviceapis

replace shared-go => ../shared-go

// Development replacements
//replace github.com/stackus/edat => ../../edat
//replace github.com/stackus/edat-msgpack => ../../edat-msgpack
//replace github.com/stackus/edat-pgx => ../../edat-pgx

require (
	github.com/mattn/go-colorable v0.1.8 // indirect
	github.com/stackus/edat v0.0.6
	github.com/stackus/edat-pgx v0.0.2
	github.com/stackus/ftgogo/serviceapis v0.0.0-20210116185538-3dd9fbb69179
	google.golang.org/grpc v1.39.0
	google.golang.org/protobuf v1.27.1
	shared-go v0.0.0-00010101000000-000000000000
)
