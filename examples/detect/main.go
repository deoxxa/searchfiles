package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

	"fknsrs.biz/p/searchfiles"
	"fknsrs.biz/p/searchfiles/detect"
)

var (
	flagDirectory string
	flagQuery     string
)

func init() {
	flag.StringVar(&flagDirectory, "directory", ".", "Directory to search in.")
	flag.StringVar(&flagQuery, "query", "", "Query to search for.")
}

func main() {
	flag.Parse()

	driverName, err := detect.Detect(context.Background(), nil)
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(os.Stderr, "> best detected driver: %s\n", driverName)

	files, err := searchfiles.SearchLiteral(context.Background(), flagDirectory, flagQuery)
	if err != nil {
		panic(err)
	}

	fmt.Println(strings.Join(files, "\n"))
}
