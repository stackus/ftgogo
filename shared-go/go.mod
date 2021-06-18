module shared-go

go 1.16

//replace github.com/stackus/edat => ../edat

require (
	github.com/go-chi/chi/v5 v5.0.3
	github.com/go-chi/cors v1.2.0
	github.com/go-chi/render v1.0.1
	github.com/google/uuid v1.2.0 // indirect
	github.com/jackc/pgproto3/v2 v2.1.0 // indirect
	github.com/jackc/pgx/v4 v4.11.0
	github.com/joho/godotenv v1.3.0
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/klauspost/compress v1.13.1 // indirect
	github.com/nats-io/nats-streaming-server v0.19.0 // indirect
	github.com/nats-io/stan.go v0.9.0
	github.com/pierrec/lz4 v2.6.1+incompatible // indirect
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.11.0
	github.com/prometheus/common v0.29.0 // indirect
	github.com/rs/zerolog v1.23.0
	github.com/shamaton/msgpack v1.2.1 // indirect
	github.com/spf13/cobra v1.1.3
	github.com/stackus/edat v0.0.6
	github.com/stackus/edat-kafka-go v0.0.1
	github.com/stackus/edat-msgpack v0.0.2
	github.com/stackus/edat-pgx v0.0.2
	github.com/stackus/edat-stan v0.0.1
	github.com/stackus/errors v0.0.3
	github.com/xdg/scram v1.0.3 // indirect
	github.com/xdg/stringprep v1.0.3 // indirect
	go.uber.org/atomic v1.8.0 // indirect
	go.uber.org/multierr v1.7.0 // indirect
	go.uber.org/zap v1.17.0
	golang.org/x/crypto v0.0.0-20210616213533-5ff15b29337e // indirect
	golang.org/x/net v0.0.0-20210614182718-04defd469f4e // indirect
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c
	golang.org/x/sys v0.0.0-20210616094352-59db8d763f22 // indirect
	google.golang.org/genproto v0.0.0-20210617175327-b9e0b3197ced // indirect
	google.golang.org/grpc v1.38.0
)
