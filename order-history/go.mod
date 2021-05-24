module github.com/stackus/ftgogo/order-history

go 1.16

replace github.com/stackus/ftgogo/serviceapis => ./../serviceapis

replace shared-go => ../shared-go

require (
	github.com/mattn/go-colorable v0.1.8 // indirect
	github.com/stackus/edat v0.0.3
	github.com/stackus/edat-pgx v0.0.1
	github.com/stackus/ftgogo/serviceapis v0.0.0-00010101000000-000000000000
	google.golang.org/protobuf v1.26.0
	shared-go v0.0.0-00010101000000-000000000000
)
