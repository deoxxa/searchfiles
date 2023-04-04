# searchfiles

[![PkgGoDev](https://pkg.go.dev/badge/fknsrs.biz/p/searchfiles)](https://pkg.go.dev/fknsrs.biz/p/searchfiles)

## [Example: Simple](./examples/simple/main.go)

This uses a simple internal search implementation. It doesn't call out to any
external processes, so it's fast to start up, but if you have to search a lot
of files it might be quite slow.

```go
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
```

## [Example: Detect](./examples/detect/main.go)

This tries to detect the "best" driver available. The `Detect` strategy
assumes you have a lot of files to search.

```go
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
```

## [Example: Advanced](./examples/advanced/main.go)

This shows most of the options of the library.

```go
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
```
