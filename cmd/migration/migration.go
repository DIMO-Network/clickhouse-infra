package main

import (
	"flag"
	"log"

	"github.com/DIMO-Network/clickhouse-infra/internal/codegen"
)

func main() {
	// Migration flags
	migrationFileName := flag.String("filename", "migration", "Name of the migration file. Default is the migration")
	packageName := flag.String("package", "migration", "Name of the package for the migration file. Default is migration")
	outputDir := flag.String("output", ".", "Output directory for the generated Go file")
	flag.Parse()
	err := codegen.Generate(*outputDir, *migrationFileName, *packageName)
	if err != nil {
		log.Fatalf("failed to generate migration file: %v", err)
	}
}
