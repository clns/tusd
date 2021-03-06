package main

import (
	"flag"
	"github.com/tus/tusd"
	"github.com/tus/tusd/filestore"
	"log"
	"net/http"
	"os"
)

var httpHost string
var httpPort string
var maxSize int64
var dir string

var stdout = log.New(os.Stdout, "[tusd] ", 0)
var stderr = log.New(os.Stderr, "[tusd] ", 0)

func init() {
	flag.StringVar(&httpHost, "host", "0.0.0.0", "Host to bind HTTP server to")
	flag.StringVar(&httpPort, "port", "1080", "Port to bind HTTP server to")
	flag.Int64Var(&maxSize, "max-size", 0, "Maximum size of uploads in bytes")
	flag.StringVar(&dir, "dir", "./data", "Directory to store uploads in")

	flag.Parse()
}

func main() {

	stdout.Printf("Using '%s' as directory storage.\n", dir)
	if err := os.MkdirAll(dir, os.FileMode(0775)); err != nil {
		stderr.Fatalf("Unable to ensure directory exists: %s", err)
	}

	store := filestore.FileStore{
		Path: dir,
	}

	stdout.Printf("Using %.2fMB as maximum size.\n", float64(maxSize)/1024/1024)

	handler, err := tusd.NewHandler(tusd.Config{
		MaxSize:   maxSize,
		BasePath:  "files/",
		DataStore: store,
	})
	if err != nil {
		stderr.Fatalf("Unable to create handler: %s", err)
	}

	address := httpHost + ":" + httpPort
	stdout.Printf("Using %s as address to listen.\n", address)

	http.Handle("/files/", http.StripPrefix("/files/", handler))
	err = http.ListenAndServe(address, nil)
	if err != nil {
		stderr.Fatalf("Unable to listen: %s", err)
	}
}
