module shared-go

go 1.16

replace github.com/stackus/edat-pgx => ./../../edat-pgx

require (
	github.com/Shopify/sarama v1.29.0
	github.com/go-chi/chi/v5 v5.0.3
	github.com/go-chi/cors v1.2.0
	github.com/go-chi/render v1.0.1
	github.com/google/uuid v1.2.0 // indirect
	github.com/jackc/pgx/v4 v4.11.0
	github.com/joho/godotenv v1.3.0
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/nats-io/nats-streaming-server v0.19.0 // indirect
	github.com/nats-io/stan.go v0.9.0
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.10.0
	github.com/prometheus/common v0.25.0 // indirect
	github.com/rs/zerolog v1.22.0
	github.com/shamaton/msgpack v1.2.1 // indirect
	github.com/spf13/cobra v1.1.3
	github.com/stackus/edat v0.0.3
	github.com/stackus/edat-msgpack v0.0.2
	github.com/stackus/edat-pgx v0.0.1
	github.com/stackus/edat-stan v0.0.1
	github.com/stackus/errors v0.0.1
	go.uber.org/multierr v1.7.0 // indirect
	go.uber.org/zap v1.16.0
	golang.org/x/crypto v0.0.0-20210513164829-c07d793c2f9a // indirect
	golang.org/x/net v0.0.0-20210521195947-fe42d452be8f // indirect
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c
	golang.org/x/sys v0.0.0-20210521203332-0cec03c779c1 // indirect
	google.golang.org/genproto v0.0.0-20210521181308-5ccab8a35a9a // indirect
	google.golang.org/grpc v1.38.0
)
