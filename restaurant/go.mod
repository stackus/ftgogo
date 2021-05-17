module github.com/stackus/ftgogo/restaurant

go 1.16

replace serviceapis => ./../serviceapis

replace shared-go => ../shared-go

require (
	github.com/deepmap/oapi-codegen v1.6.1
	github.com/getkin/kin-openapi v0.53.0
	github.com/go-chi/chi/v5 v5.0.3
	github.com/go-chi/render v1.0.1
	github.com/go-openapi/swag v0.19.14 // indirect
	github.com/google/uuid v1.2.0
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/stackus/edat v0.0.3
	github.com/stackus/edat-pgx v0.0.1
	serviceapis v0.0.0-00010101000000-000000000000
	shared-go v0.0.0-00010101000000-000000000000
)
