package conf

import (
	"github.com/scylladb/gocqlx/v2"
	tc "go.temporal.io/sdk/client"
)

type github struct {
	AppID         int64  `env:"GITHUB_APP_ID"`
	ClinetID      string `env:"GITHUB_CLIENT_ID"`
	WebhookSecret string `env:"GITHUB_WEBHOOK_SECRET"`
	PrivateKey    string `env:"GITHUB_PRIVATE_KEY"`
	// PrivateKey Base64EncodedValue `env:"GITHUB_PRIVATE_KEY"`
}

type kratos struct {
	ServerUrl string `env:"KRATOS_SERVER_URL"`
}

type temporal struct {
	ServerHost string `env:"TEMPORAL_HOST"`
	ServerPort string `env:"TEMPORAL_PORT" env-default:"7233"`
	Client     tc.Client
	Queues     struct {
		Webhooks string `env-default:"webhooks"`
	}
}

func (t *temporal) GetConnectionString() string {
	return t.ServerHost + ":" + t.ServerPort
}

type service struct {
	Name    string
	Debug   bool   `env:"DEBUG" env-default:"false"`
	Version string `env:"VERSION" env-default:"0.0.0-dev"`
}

type cassandra struct {
	Hosts    []string `env:"CASSANDRA_HOSTS" env-default:"cassandra"`
	KeySpace string   `env:"CASSANDRA_KEYSPACE" env-default:"ctrlplane"`
	Session  gocqlx.Session
}
