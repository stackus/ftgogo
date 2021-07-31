module github.com/stackus/ftgogo/kitchen

go 1.16

replace github.com/stackus/ftgogo/serviceapis => ./../serviceapis

replace shared-go => ../shared-go

// Development replacements
//replace github.com/stackus/edat => ../../edat
//replace github.com/stackus/edat-msgpack => ../../edat-msgpack
//replace github.com/stackus/edat-pgx => ../../edat-pgx

require (
	github.com/cucumber/godog v0.11.0
	github.com/hashicorp/go-immutable-radix v1.3.1 // indirect
	github.com/hashicorp/go-memdb v1.3.2 // indirect
	github.com/mattn/go-colorable v0.1.8 // indirect
	github.com/rdumont/assistdog v0.0.0-20201106100018-168b06230d14
	github.com/spf13/pflag v1.0.5
	github.com/stackus/edat v0.0.6
	github.com/stackus/edat-msgpack v0.0.2
	github.com/stackus/edat-pgx v0.0.2
	github.com/stackus/errors v0.0.3
	github.com/stackus/ftgogo/serviceapis v0.0.0-20210116185538-3dd9fbb69179
	google.golang.org/grpc v1.39.0
	shared-go v0.0.0-00010101000000-000000000000
)
