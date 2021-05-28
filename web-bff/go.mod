module github.com/stackus/ftgogo/web-bff

go 1.16

replace github.com/stackus/ftgogo/serviceapis => ./../serviceapis

replace shared-go => ../shared-go

require (
	github.com/deepmap/oapi-codegen v1.7.0
	github.com/getkin/kin-openapi v0.62.0
	github.com/go-chi/chi/v5 v5.0.3
	github.com/go-chi/jwtauth/v5 v5.0.1
	github.com/go-chi/render v1.0.1
	github.com/goccy/go-json v0.5.1 // indirect
	github.com/lestrrat-go/backoff/v2 v2.0.8 // indirect
	github.com/lestrrat-go/jwx v1.2.0
	github.com/stackus/errors v0.0.2
	github.com/stackus/ftgogo/serviceapis v0.0.0-20210116185538-3dd9fbb69179
	google.golang.org/grpc v1.38.0
	google.golang.org/protobuf v1.26.0
	shared-go v0.0.0-00010101000000-000000000000
)
