package connect

import (
	"crypto/tls"
	"database/sql"
	"fmt"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/DIMO-Network/clickhouse-infra/pkg/connect/config"
)

const (
	defaultDialTimeout = time.Second * 30
	defaultReadTimeout = time.Minute * 5
)

// GetClickhouseDB returns a sql.DB connection to clickhouse.
func GetClickhouseDB(settings *config.Settings) *sql.DB {
	dialTimeout := defaultDialTimeout
	if settings.DialTimeout != "" {
		sDialTimeout, err := time.ParseDuration(settings.DialTimeout)
		if err == nil {
			dialTimeout = sDialTimeout
		}
	}
	readTimeout := defaultReadTimeout
	if settings.ReadTimeout != "" {
		sReadTimeout, err := time.ParseDuration(settings.ReadTimeout)
		if err == nil {
			readTimeout = sReadTimeout
		}
	}
	addr := fmt.Sprintf("%s:%d", settings.Host, settings.Port)
	conn := clickhouse.OpenDB(&clickhouse.Options{
		Addr: []string{addr},
		Auth: clickhouse.Auth{
			Username: settings.User,
			Password: settings.Password,
			Database: settings.Database,
		},
		DialTimeout: dialTimeout,
		ReadTimeout: readTimeout,
		TLS: &tls.Config{
			MinVersion: tls.VersionTLS12,
			RootCAs:    settings.RootCAs,
		},
	})
	return conn
}

// GetClickhouseConn returns a clickhouse.Conn connection to clickhouse.
func GetClickhouseConn(settings *config.Settings) (clickhouse.Conn, error) {
	dialTimeout := defaultDialTimeout
	if settings.DialTimeout != "" {
		var err error
		dialTimeout, err = time.ParseDuration(settings.DialTimeout)
		if err != nil {
			return nil, fmt.Errorf("failed to parse dial timeout: %w", err)
		}
	}
	readTimeout := defaultReadTimeout
	if settings.ReadTimeout != "" {
		var err error
		readTimeout, err = time.ParseDuration(settings.ReadTimeout)
		if err != nil {
			return nil, fmt.Errorf("failed to parse read timeout: %w", err)
		}
	}
	addr := fmt.Sprintf("%s:%d", settings.Host, settings.Port)
	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{addr},
		Auth: clickhouse.Auth{
			Username: settings.User,
			Password: settings.Password,
			Database: settings.Database,
		},
		DialTimeout: dialTimeout,
		ReadTimeout: readTimeout,
		TLS: &tls.Config{
			MinVersion: tls.VersionTLS12,
			RootCAs:    settings.RootCAs,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to open clickhouse connection: %w", err)
	}
	return conn, nil
}
