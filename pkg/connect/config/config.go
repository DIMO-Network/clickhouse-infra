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
	// ReadTimeout is the timeout for reading from the clickhouse database.
	// defaults to 5ms
	ReadTimeout string `env:"CLICKHOUSE_READ_TIMEOUT" yaml:"CLICKHOUSE_READ_TIMEOUT"`
	// DialTimeout is the timeout for dialing the clickhouse database.
	// defaults to 30s
	DialTimeout string `env:"CLICKHOUSE_DIAL_TIMEOUT" yaml:"CLICKHOUSE_DIAL_TIMEOUT"`

	RootCAs *x509.CertPool `env:"-" yaml:"-"`
}

// DSN returns the Data Source Name (DSN) for connecting to ClickHouse.
func (s Settings) DSN() string {
	dialTimeout := s.DialTimeout
	if dialTimeout == "" {
		dialTimeout = "200ms"
	}
	readTimeout := s.ReadTimeout
	if readTimeout == "" {
		readTimeout = "5m"
	}
	return fmt.Sprintf("clickhouse://%s:%d/%s?username=%s&password=%s&secure=true&dial_timeout=%s&read_timeout=%s&max_execution_time=60", s.Host, s.Port, s.Database, url.QueryEscape(s.User), url.QueryEscape(s.Password), dialTimeout, readTimeout)
}
