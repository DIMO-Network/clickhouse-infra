// Package codegen provides a function for generating migration files.
package codegen

import (
	"bytes"
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"golang.org/x/tools/imports"
)

var (
	delemReplacer = strings.NewReplacer("_", " ", "-", " ", ".", " ")
	titleCaser    = cases.Title(language.AmericanEnglish, cases.NoLower)
	lowerCaser    = cases.Lower(language.AmericanEnglish)
)

const (
	migrationFileFormat = "%s_%s_migration.go"
	generatedFileName   = "generated.go"
)

//go:embed migration.tmpl
var migrationFileTemplate string

//go:embed generated.tmpl
var generatedFileTemplate string

// Config is the configuration for the migration generator.
type Config struct {
	// fileName is the name of the migration file.
	FileName string
}

type tmplData struct {
	PackageName string
	FuncName    string
}

// Generate creates a new ClickHouse table file.
func Generate(outputDir, fileName, packageName string) error {
	version, err := getVersion(outputDir)
	if err != nil {
		return fmt.Errorf("error getting version: %w", err)
	}

	migrationTmpl, err := template.New("migrationTemplate").Parse(migrationFileTemplate)
	if err != nil {
		return fmt.Errorf("error parsing migration template: %w", err)
	}

	generatedTmpl, err := template.New("generatedTemplate").Parse(generatedFileTemplate)
	if err != nil {
		return fmt.Errorf("error parsing generated template: %w", err)
	}

	var outBuf bytes.Buffer
	data := tmplData{
		PackageName: packageName,
		FuncName:    strings.ReplaceAll(titleCaser.String(fileName), " ", ""),
	}
	err = migrationTmpl.Execute(&outBuf, &data)
	if err != nil {
		return fmt.Errorf("error executing ClickHouse table template: %w", err)
	}

	outFile := delemReplacer.Replace(fileName)
	outFile = getFilePath(outputDir, outFile, version)
	err = formatAndWriteToFile(outBuf.Bytes(), outFile)
	if err != nil {
		return fmt.Errorf("error writing file: %w", err)
	}

	outBuf.Reset()
	err = generatedTmpl.Execute(&outBuf, &data)
	if err != nil {
		return fmt.Errorf("error executing generated template: %w", err)
	}
	outFile = filepath.Clean(filepath.Join(outputDir, generatedFileName))
	err = formatAndWriteToFile(outBuf.Bytes(), outFile)
	if err != nil {
		return fmt.Errorf("error writing file: %w", err)
	}

	return nil
}

func getFilePath(outputDir, fileName, version string) string {
	noSpaceName := lowerCaser.String(strings.ReplaceAll(fileName, " ", "_"))
	migrationFileName := fmt.Sprintf(migrationFileFormat, version, noSpaceName)
	return filepath.Clean(filepath.Join(outputDir, migrationFileName))
}

func getVersion(outputDir string) (string, error) {
	files, err := os.ReadDir(outputDir)
	if err != nil {
		return "", fmt.Errorf("error reading directory: %w", err)
	}
	var maxVersion int
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		name := file.Name()
		if strings.HasSuffix(name, ".go") {
			versionStr := strings.Split(name, "_")[0]
			version, err := strconv.Atoi(versionStr)
			if err != nil {
				continue
			}
			if version > maxVersion {
				maxVersion = version
			}
		}
	}
	return fmt.Sprintf("%05d", maxVersion+1), nil
}

// formatAndWriteToFile formats the go source with goimports and writes it to the output file.
func formatAndWriteToFile(goData []byte, outputFilePath string) (err error) {
	cleanPath := filepath.Clean(outputFilePath)
	formatted, fmtErr := imports.Process(cleanPath, goData, &imports.Options{
		AllErrors: true,
		Comments:  true,
	})
	if fmtErr != nil {
		// do not return early, we still want to write the file
		fmtErr = fmt.Errorf("error formatting go source: %w", fmtErr)
		formatted = goData
	}
	goOutputFile, err := os.Create(cleanPath)
	if err != nil {
		return fmt.Errorf("error creating output file: %w", err)
	}
	defer func() {
		if cerr := goOutputFile.Close(); err == nil && cerr != nil {
			err = cerr
		}
	}()
	_, err = goOutputFile.Write(formatted)
	if err != nil {
		return fmt.Errorf("error writing to file: %w", err)
	}

	// return the formatting error if there is one
	return fmtErr
}
