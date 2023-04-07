package common

import (
	"fmt"
	"strconv"
	"time"

	e "github.com/Walker088/gorealestate/error"
)

const (
	RocEraFormattingError = "PS00002"

	currentPackage = "github.com/Walker088/gorealestate/common/time"
)

// Input: ROC Era, e.g., 1011019 = 2012/10/19
func RocEraToCommonEra(dateStr string) (*time.Time, *e.ErrorData) {
	if len(dateStr) < 6 {
		return nil, e.NewErrorData(
			RocEraFormattingError,
			fmt.Sprintf("Incorrect ROC era format %s, expect length > 5, e.g, 1011019, 891011", dateStr),
			fmt.Sprintf("%s.RocEraToCommonEra", currentPackage),
			nil,
			nil,
		)
	}

	rocToCommonEra := func(era int) int {
		return era + 1911
	}

	rocYearStr := dateStr[:len(dateStr)-4]
	rocYear, err := strconv.Atoi(rocYearStr)
	if err != nil {
		return nil, e.NewErrorData(
			RocEraFormattingError,
			fmt.Sprintf("Incorrect ROC era format: %s, year extracting error %s", dateStr, err.Error()),
			fmt.Sprintf("%s.RocEraToCommonEra", currentPackage),
			nil,
			nil,
		)
	}

	monthAndDay := dateStr[len(rocYearStr):] //SliceDiff(dateStr, rocYearStr)

	month, err := strconv.Atoi(monthAndDay[:len(monthAndDay)-2])
	if err != nil {
		return nil, e.NewErrorData(
			RocEraFormattingError,
			fmt.Sprintf("Incorrect ROC era format: %s, month extracting error %s", dateStr, err.Error()),
			fmt.Sprintf("%s.RocEraToCommonEra", currentPackage),
			nil,
			nil,
		)
	}

	day, err := strconv.Atoi(monthAndDay[2:])
	if err != nil {
		return nil, e.NewErrorData(
			RocEraFormattingError,
			fmt.Sprintf("Incorrect ROC era format: %s, day extracting error %s", dateStr, err.Error()),
			fmt.Sprintf("%s.RocEraToCommonEra", currentPackage),
			nil,
			nil,
		)
	}

	converted := time.Date(rocToCommonEra(int(rocYear)), time.Month(month), day, 0, 0, 0, 0, time.Local)
	return &converted, nil
}
