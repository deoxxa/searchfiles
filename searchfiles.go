package searchfiles

import (
	"context"
	"fmt"
)

var (
	ErrUnimplemented = fmt.Errorf("unimplemented")
	ErrUnknownDriver = fmt.Errorf("no driver found with this name")
	ErrNoDrivers     = fmt.Errorf("no drivers registered; try using fknsrs.biz/p/searchfiles/detect or fknsrs.biz/p/searchfiles/driver/native")
)

type Driver interface {
	SelfTest(ctx context.Context) error
	SearchLiteral(ctx context.Context, directory, query string) ([]string, error)
	SearchRegexp(ctx context.Context, directory, query string) ([]string, error)
}

var (
	drivers         = map[string]Driver{}
	preferredDriver = "native"
)

func Register(driverName string, driver Driver) {
	drivers[driverName] = driver
}

func DriverNames() []string {
	var a []string

	for k := range drivers {
		a = append(a, k)
	}

	return a
}

func SetPreferredDriver(driverName string) {
	preferredDriver = driverName
}

func getDriver(driverName string) (Driver, error) {
	if len(drivers) == 0 {
		return nil, fmt.Errorf("searchfiles.getDriver: %w", ErrNoDrivers)
	}

	if driverName == "" {
		driverName = preferredDriver
	}

	driver, ok := drivers[driverName]
	if !ok {
		return nil, fmt.Errorf("searchfiles.getDriver: %w", ErrUnknownDriver)
	}

	return driver, nil
}

func TestDriver(ctx context.Context, driverName string) error {
	driver, ok := drivers[driverName]
	if !ok {
		return fmt.Errorf("searchfiles.TestDriver: %w", ErrUnknownDriver)
	}

	if err := driver.SelfTest(ctx); err != nil {
		return fmt.Errorf("searchfiles.TestDriver: %w", err)
	}

	return nil
}

func SearchLiteral(ctx context.Context, directory, query string) ([]string, error) {
	res, err := SearchLiteralUsing(ctx, "", directory, query)
	if err != nil {
		return nil, fmt.Errorf("searchfiles.SearchLiteral: %w", err)
	}

	return res, nil
}

func SearchRegexp(ctx context.Context, directory, query string) ([]string, error) {
	res, err := SearchRegexpUsing(ctx, "", directory, query)
	if err != nil {
		return nil, fmt.Errorf("searchfiles.SearchRegexp: %w", err)
	}

	return res, nil
}

func SearchLiteralUsing(ctx context.Context, driverName string, directory, query string) ([]string, error) {
	driver, err := getDriver(driverName)
	if err != nil {
		return nil, fmt.Errorf("searchfiles.SearchLiteralUsing: %w", err)
	}

	a, err := driver.SearchLiteral(ctx, directory, query)
	if err != nil {
		return nil, fmt.Errorf("searchfiles.SearchLiteralUsing: %w", err)
	}

	return a, nil
}

func SearchRegexpUsing(ctx context.Context, driverName string, directory, query string) ([]string, error) {
	driver, err := getDriver(driverName)
	if err != nil {
		return nil, fmt.Errorf("searchfiles.SearchRegexpUsing: %w", err)
	}

	a, err := driver.SearchRegexp(ctx, directory, query)
	if err != nil {
		return nil, fmt.Errorf("searchfiles.SearchRegexpUsing: %w", err)
	}

	return a, nil
}
