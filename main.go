package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func init() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage of genorm:")
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "  genorm -type=<type>")
		fmt.Fprintln(os.Stderr, "")
		flag.PrintDefaults()
	}

	flag.CommandLine.Init("", flag.ExitOnError)
}

func main() {
	driver := flag.String("driver", "", "The driver to be used")
	typeName := flag.String("type", "", "Type that hosts io.WriterTo interface implementation")
	packageName := flag.String("package", "", "Package name")
	url := flag.String("url", "", "Database url")
	tableName := flag.String("tableName", "", "Existing table name")

	flag.Parse()

	pkgDir, err := packageDir(*packageName)
	if err != nil {
		panic(err)
	}

	outputFile := formatFileName(*typeName)
	//fmt.Printf("Generating file: %s\n", outputFile)

	writer, err := os.Create(filepath.Join(pkgDir, outputFile))
	if err != nil {
		os.Remove(filepath.Join(pkgDir, outputFile))
		writer, err = os.Create(filepath.Join(pkgDir, outputFile))
		if err != nil {
			panic(err)
		}
	}
	defer writer.Close()

	generator := &Generator{}

	m := metadata(*typeName, *driver, *url, *tableName, pkgDir)

	//fmt.Printf("Metadata to apply %+v\n", m)

	if err := generator.Generate(writer, m); err != nil {
		panic(err)
	}

	fmt.Printf("Generated %s\n", outputFile)
}

func formatFileName(typeName string) string {
	return fmt.Sprintf("%s-repository.go", strings.ToLower(typeName))
}

func packageDir(packageName string) (string, error) {
	if packageName == "" {
		return os.Getwd()
	}

	path := os.Getenv("GOPATH")
	if path == "" {
		return "", errors.New("GOPATH is not set")
	}

	workDir := filepath.Join(path, "src", packageName)
	if _, err := os.Stat(workDir); err != nil {
		return "", err
	}

	return workDir, nil
}

func metadata(typeName string, driver string, url string, tableName string, packageDir string) (m Metadata) {
	m.Type = typeName
	m.PackageName = filepath.Base(packageDir)
	m.Driver = driver
	m.URL = url
	m.TableName = tableName
	return m
}
