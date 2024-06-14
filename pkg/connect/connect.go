package connect

import (
	"crypto/tls"
	"database/sql"
	"fmt"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/DIMO-Network/clickhouse-infra/pkg/connect/config"
)

// GetClickhouseDB returns a sql.DB connection to clickhouse.
func GetClickhouseDB(settings *config.Settings) *sql.DB {
	addr := fmt.Sprintf("%s:%d", settings.Host, settings.Port)
	conn := clickhouse.OpenDB(&clickhouse.Options{
		Addr: []string{addr},
		Auth: clickhouse.Auth{
			Username: settings.User,
			Password: settings.Password,
			Database: settings.Database,
		},
		DialTimeout: time.Minute * 30,
		TLS: &tls.Config{
			MinVersion: tls.VersionTLS12,
			RootCAs:    settings.RootCAs,
		},
	})
	return conn
}

// GetClickhouseConn returns a clickhouse.Conn connection to clickhouse.
func GetClickhouseConn(settings *config.Settings) (clickhouse.Conn, error) {
	addr := fmt.Sprintf("%s:%d", settings.Host, settings.Port)
	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{addr},
		Auth: clickhouse.Auth{
			Username: settings.User,
			Password: settings.Password,
			Database: settings.Database,
		},
		DialTimeout: time.Minute * 30,
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
