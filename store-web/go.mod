module github.com/stackus/ftgogo/store-web

go 1.16

replace github.com/stackus/ftgogo/serviceapis => ./../serviceapis

replace shared-go => ../shared-go

// Development replacements
//replace github.com/stackus/edat => ../../edat
//replace github.com/stackus/edat-msgpack => ../../edat-msgpack
//replace github.com/stackus/edat-pgx => ../../edat-pgx

require (
	github.com/deepmap/oapi-codegen v1.8.2
	github.com/getkin/kin-openapi v0.68.0
	github.com/go-chi/chi/v5 v5.0.3
	github.com/go-chi/render v1.0.1
	github.com/stackus/ftgogo/serviceapis v0.0.0-20210116185538-3dd9fbb69179
	shared-go v0.0.0-00010101000000-000000000000
)
