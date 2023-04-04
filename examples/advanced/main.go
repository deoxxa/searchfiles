package main

import (
	"context"
	"flag"
	"fmt"
	"strings"

	"fknsrs.biz/p/searchfiles"
	_ "fknsrs.biz/p/searchfiles/driver/ag"
	_ "fknsrs.biz/p/searchfiles/driver/grep"
	_ "fknsrs.biz/p/searchfiles/driver/native"
	_ "fknsrs.biz/p/searchfiles/driver/pt"
	_ "fknsrs.biz/p/searchfiles/driver/rg"
)

var (
	flagDirectory string
	flagQuery     string
	flagRegexp    bool
	flagDriver    string
)

func init() {
	flag.StringVar(&flagDirectory, "directory", ".", "Directory to search in.")
	flag.StringVar(&flagQuery, "query", "", "Query to search for.")
	flag.BoolVar(&flagRegexp, "regexp", false, "Search for a regular expression rather than a static string.")
	flag.StringVar(&flagDriver, "driver", "native", "Choose a driver to use (ag, grep, native, pt, rg).")
}

func main() {
	flag.Parse()

	var files []string

	if flagRegexp {
		a, err := searchfiles.SearchRegexpUsing(context.Background(), flagDriver, flagDirectory, flagQuery)
		if err != nil {
			panic(err)
		}
		files = a
	} else {
		a, err := searchfiles.SearchLiteralUsing(context.Background(), flagDriver, flagDirectory, flagQuery)
		if err != nil {
			panic(err)
		}
		files = a
	}

	fmt.Println(strings.Join(files, "\n"))
}
