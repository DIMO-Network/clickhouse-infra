// Package migrate provides the functionality to run goose migrations on a clickhouse database.
package migrate

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"sync"

	"github.com/pressly/goose/v3"
)

// migrationLock allows to run goose commands concurrently. Since goose leverages global variables to store the migrations.
var migrationLock sync.Mutex

// setMigrations sets the migrations for the goose tool.
// this will reset the global migrations and FS to avoid any unwanted migrations registers.
func setMigrations(registerFuncs []func()) {
	emptyFs := embed.FS{}
	goose.SetBaseFS(emptyFs)
	goose.ResetGlobalMigrations()
	for _, regFunc := range registerFuncs {
		regFunc()
	}
}

// RunGoose runs the goose command with the provided arguments.
// args should be the command and the arguments to pass to goose.
// eg RunGoose(ctx, []string{"up", "-v"}, db).
// registerFuncs should be a list of functions that register the migrations.
// This function is safe to run concurrently.
func RunGoose(ctx context.Context, gooseArgs []string, registerFuncs []func(), db *sql.DB) error {
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
	setMigrations(registerFuncs)
	if err := goose.SetDialect("clickhouse"); err != nil {
		return fmt.Errorf("failed to set dialect: %w", err)
	}
	err := goose.RunContext(ctx, cmd, db, ".", args...)
	if err != nil {
		return fmt.Errorf("failed to run goose command: %w", err)
	}
	return nil
}
