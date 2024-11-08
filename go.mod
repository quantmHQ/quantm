module go.breu.io/quantm

go 1.23.1

require (
	connectrpc.com/connect v1.17.0
	github.com/avast/retry-go/v4 v4.6.0
	github.com/bradleyfalzon/ghinstallation/v2 v2.11.0
	github.com/go-jose/go-jose/v4 v4.0.4
	github.com/go-playground/validator/v10 v10.22.1
	github.com/gocql/gocql v1.7.0
	github.com/google/go-github/v62 v62.0.0
	github.com/gosimple/slug v1.14.0
	github.com/jackc/pgx/v5 v5.7.1
	github.com/knadh/koanf/providers/env v1.0.0
	github.com/knadh/koanf/providers/structs v0.1.0
	github.com/knadh/koanf/v2 v2.1.1
	github.com/labstack/echo/v4 v4.12.0
	github.com/sethvargo/go-password v0.3.1
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/testify v1.9.0
	go.breu.io/durex v0.6.2
	go.breu.io/graceful v0.1.0
	go.opentelemetry.io/otel/trace v1.31.0
	go.step.sm/crypto v0.54.0
	go.temporal.io/sdk v1.30.0
	golang.org/x/crypto v0.28.0
	golang.org/x/net v0.30.0
	google.golang.org/genproto/googleapis/rpc v0.0.0-20241021214115-324edc3d5d38
)

require github.com/go-test/deep v1.0.8 // indirect

require (
	cloud.google.com/go/compute/metadata v0.5.2
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/facebookgo/clock v0.0.0-20150410010913-600d898af40a // indirect
	github.com/fatih/structs v1.1.0 // indirect
	github.com/gabriel-vasile/mimetype v1.4.3 // indirect
	github.com/go-jose/go-jose/v3 v3.0.3 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-viper/mapstructure/v2 v2.0.0-alpha.1 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang-jwt/jwt/v4 v4.5.0 // indirect
	github.com/golang-migrate/migrate/v4 v4.18.1
	github.com/golang/mock v1.6.0 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/google/go-querystring v1.1.0 // indirect
	github.com/google/uuid v1.6.0
	github.com/gorilla/websocket v1.5.3 // indirect
	github.com/gosimple/unidecode v1.0.1 // indirect
	github.com/grpc-ecosystem/go-grpc-middleware v1.4.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.22.0 // indirect
	github.com/hailocab/go-hostpool v0.0.0-20160125115350-e80d13ce29ed // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/knadh/koanf/maps v0.1.1 // indirect
	github.com/labstack/gommon v0.4.2 // indirect
	github.com/leodido/go-urn v1.4.0 // indirect
	github.com/lib/pq v1.10.9 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mitchellh/copystructure v1.2.0 // indirect
	github.com/mitchellh/reflectwalk v1.0.2 // indirect
	github.com/nexus-rpc/sdk-go v0.0.11 // indirect
	github.com/pborman/uuid v1.2.1 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/robfig/cron v1.2.0 // indirect
	github.com/slack-go/slack v0.15.0
	github.com/stretchr/objx v0.5.2 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasttemplate v1.2.2 // indirect
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.55.0 // indirect
	go.opentelemetry.io/otel v1.31.0 // indirect
	go.temporal.io/api v1.41.0 // indirect
	go.uber.org/atomic v1.11.0 // indirect
	golang.org/x/exp v0.0.0-20240409090435-93d18d7e34b8 // indirect
	golang.org/x/sync v0.8.0 // indirect
	golang.org/x/sys v0.26.0 // indirect
	golang.org/x/text v0.19.0 // indirect
	golang.org/x/time v0.7.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20241007155032-5fefd90f89a9 // indirect
	google.golang.org/grpc v1.67.1
	google.golang.org/protobuf v1.35.1
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/Guilospanck/igocqlx v1.0.0 => github.com/debuggerpk/igocqlx v1.2.0

replace github.com/deepmap/oapi-codegen/v2 v2.1.0 => github.com/breuHQ/oapi-codegen/v2 v2.1.1-breu

// replace github.com/deepmap/oapi-codegen/v2 v2.1.0 => /Users/jay/Work/opensource/oapi-codegen

// replace go.breu.io/durex v0.4.0 => /Users/jay/Work/breu/durex
