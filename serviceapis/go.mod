module github.com/stackus/ftgogo/serviceapis

go 1.16

// Development replacements
//replace github.com/stackus/edat => ../../edat
//replace github.com/stackus/edat-msgpack => ../../edat-msgpack
//replace github.com/stackus/edat-pgx => ../../edat-pgx

require (
	github.com/getkin/kin-openapi v0.68.0
	github.com/go-openapi/swag v0.19.15 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/stackus/edat v0.0.6
	golang.org/x/net v0.0.0-20210726213435-c6fcb2dbf985 // indirect
	golang.org/x/sys v0.0.0-20210630005230-0f9fa26af87c // indirect
	google.golang.org/genproto v0.0.0-20210729151513-df9385d47c1b // indirect
	google.golang.org/grpc v1.39.0
	google.golang.org/protobuf v1.27.1
)
