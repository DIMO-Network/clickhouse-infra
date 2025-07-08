package config

import (
	"crypto/x509"
	"fmt"
	"net/url"
)

// Clickhouse is the configuration for connecting to a clickhouse database.
type Settings struct {
	Host     string `env:"CLICKHOUSE_HOST"     yaml:"CLICKHOUSE_HOST"`
	Port     int    `env:"CLICKHOUSE_TCP_PORT" yaml:"CLICKHOUSE_TCP_PORT"`
	User     string `env:"CLICKHOUSE_USER"     yaml:"CLICKHOUSE_USER"`
	Password string `env:"CLICKHOUSE_PASSWORD" yaml:"CLICKHOUSE_PASSWORD"`
	Database string `env:"CLICKHOUSE_DATABASE" yaml:"CLICKHOUSE_DATABASE"`

	RootCAs *x509.CertPool `env:"-" yaml:"-"`
}

// DSN returns the Data Source Name (DSN) for connecting to ClickHouse.
func (s Settings) DSN() string {
	return fmt.Sprintf("clickhouse://%s:%d/%s?username=%s&password=%s&secure=true&dial_timeout=200ms&max_execution_time=60", s.Host, s.Port, s.Database, url.QueryEscape(s.User), url.QueryEscape(s.Password))
}
