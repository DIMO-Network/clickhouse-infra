# clickhouse-infra

![GitHub license](https://img.shields.io/badge/license-Apache%202.0-blue.svg)
[![GoDoc](https://godoc.org/github.com/DIMO-Network/clickhouse-infra?status.svg)](https://godoc.org/github.com/DIMO-Network/clickhouse-infra)
[![Go Report Card](https://goreportcard.com/badge/github.com/DIMO-Network/clickhouse-infra)](https://goreportcard.com/report/github.com/DIMO-Network/clickhouse-infra)

A Go library providing common ClickHouse database utilities for DIMO Network applications.

## Features

- **Database Migrations**: Thread-safe migration utilities using Goose for ClickHouse schema management
- **Connection Management**: Simplified connection helpers for both native ClickHouse protocol and SQL database connections
- **Test Containers**: Docker-based ClickHouse test containers with TLS support for integration testing
- **Configuration**: Structured configuration management with environment variable support

## Development

This project uses a standard Makefile. Try `make help` to see the available commands.

## License

[Apache 2.0](LICENSE)
