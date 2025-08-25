// Package migrate helps with concurrency safe goose migrations for clickhouse.
package migrate

import (
	"context"
	"database/sql"
	"fmt"
	"io/fs"
	"os"
	"sync"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/pressly/goose/v3"
)

// migrationLock allows to run goose commands concurrently. Since goose leverages global variables to store the migrations.
var migrationLock sync.Mutex

// setMigrations sets the migrations for the goose tool.
// this will reset the global migrations and FS to avoid any unwanted migrations registers.
func setMigrations(baseFS fs.FS) {
	goose.SetBaseFS(baseFS)
	goose.ResetGlobalMigrations()
}

// RunGoose runs the goose command with the provided arguments.
// args should be the command and the arguments to pass to goose.
// eg RunGoose(ctx, []string{"up", "-v"}, db).
// baseFS should be the filesystem that holds versioned migration files.
// This function is safe to run concurrently.
func RunGoose(ctx context.Context, gooseArgs []string, baseFS fs.FS, db *sql.DB) error {
	migrationLock.Lock()
	defer migrationLock.Unlock()
	if len(gooseArgs) == 0 {
		return fmt.Errorf("command not provided")
	}
	cmd := gooseArgs[0]
	var args []string
	if len(gooseArgs) > 1 {
		args = gooseArgs[1:]
	}
	setMigrations(baseFS)
	if err := goose.SetDialect("clickhouse"); err != nil {
		return fmt.Errorf("failed to set dialect: %w", err)
	}
	err := goose.RunContext(ctx, cmd, db, ".", args...)
	if err != nil {
		return fmt.Errorf("failed to run goose command: %w", err)
	}
	return nil
}

// RunGooseCmd parses cmdline arguments and runs the goose command using the provided baseFS.
func RunGooseCmd(ctx context.Context, baseFS fs.FS) error {
	args := os.Args

	if len(args) < 2 {
		return fmt.Errorf("usage: %s <dbstring> <command> [args]", args[0])
	}

	dbstring := args[1]
	dbOptions, err := clickhouse.ParseDSN(dbstring)
	if err != nil {
		return fmt.Errorf("failed to parse DSN: %w", err)
	}
	sqlDB := clickhouse.OpenDB(dbOptions)

	err = RunGoose(ctx, args[2:], baseFS, sqlDB)
	if err != nil {
		_ = sqlDB.Close()
		return fmt.Errorf("failed to run goose command: %w", err)
	}
	if err := sqlDB.Close(); err != nil {
		return fmt.Errorf("failed to close db: %w", err)
	}
	return nil
}
