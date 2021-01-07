module github.com/stackus/ftgogo/delivery

go 1.15

replace serviceapis => ./../serviceapis

replace shared-go => ../shared-go

require (
	github.com/deepmap/oapi-codegen v1.4.1
	github.com/getkin/kin-openapi v0.35.0
	github.com/go-chi/chi v4.1.2+incompatible
	github.com/go-chi/render v1.0.1
	github.com/google/uuid v1.1.4
	github.com/stackus/edat v0.0.1
	github.com/stackus/edat-pgx v0.0.1
	serviceapis v0.0.0-00010101000000-000000000000
	shared-go v0.0.0-00010101000000-000000000000
)
