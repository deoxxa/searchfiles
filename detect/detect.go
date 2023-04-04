package detect

import (
	"context"
	"fmt"

	"fknsrs.biz/p/searchfiles"
	_ "fknsrs.biz/p/searchfiles/driver/ag"
	_ "fknsrs.biz/p/searchfiles/driver/grep"
	_ "fknsrs.biz/p/searchfiles/driver/native"
	_ "fknsrs.biz/p/searchfiles/driver/pt"
	_ "fknsrs.biz/p/searchfiles/driver/rg"
)

var (
	ErrNoWorkingDriver = fmt.Errorf("no working driver found")
)

var DefaultSearchOrder = []string{"ag", "rg", "grep", "pt", "native"}

func Detect(ctx context.Context, searchOrder []string) (string, error) {
	if searchOrder == nil {
		searchOrder = DefaultSearchOrder
	}

	for _, driverName := range searchOrder {
		if err := searchfiles.TestDriver(ctx, driverName); err == nil {
			return driverName, nil
		}
	}

	return "", fmt.Errorf("detect.Detect: %w", ErrNoWorkingDriver)
}

func DetectAndSetPreferred(ctx context.Context, searchOrder []string) (string, error) {
	driverName, err := Detect(ctx, searchOrder)
	if err != nil {
		return "", fmt.Errorf("detect.DetectAndSetPreferred: %w", err)
	}

	searchfiles.SetPreferredDriver(driverName)

	return driverName, nil
}
