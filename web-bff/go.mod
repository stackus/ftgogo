module github.com/stackus/ftgogo/web-bff

go 1.16

replace github.com/stackus/ftgogo/serviceapis => ./../serviceapis

replace shared-go => ../shared-go

require (
	github.com/deepmap/oapi-codegen v1.7.0
	github.com/getkin/kin-openapi v0.62.0 // indirect
	github.com/go-chi/chi/v5 v5.0.3
	github.com/go-chi/jwtauth/v5 v5.0.1
	github.com/go-chi/render v1.0.1
	github.com/go-openapi/swag v0.19.15 // indirect
	github.com/labstack/echo/v4 v4.3.0 // indirect
	github.com/lestrrat-go/jwx v1.2.0
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/stackus/errors v0.0.2
	github.com/stackus/ftgogo/serviceapis v0.0.0-00010101000000-000000000000
	golang.org/x/net v0.0.0-20210525063256-abc453219eb5 // indirect
	golang.org/x/sys v0.0.0-20210525143221-35b2ab0089ea // indirect
	golang.org/x/tools v0.1.2 // indirect
	google.golang.org/grpc v1.38.0
	google.golang.org/protobuf v1.26.0
	shared-go v0.0.0-00010101000000-000000000000
)
