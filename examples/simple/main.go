package main

import (
	"context"
	"flag"
	"fmt"
	"strings"

	"fknsrs.biz/p/searchfiles"
	_ "fknsrs.biz/p/searchfiles/driver/native"
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

	files, err := searchfiles.SearchLiteral(context.Background(), flagDirectory, flagQuery)
	if err != nil {
		panic(err)
	}

	fmt.Println(strings.Join(files, "\n"))
}
