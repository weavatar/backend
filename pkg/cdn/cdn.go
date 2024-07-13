package cdn

import (
	"slices"
	"strings"

	"github.com/goravel/framework/facades"
	"github.com/goravel/framework/support/carbon"
)

var drivers []Driver

func register(driver Driver) {
	drivers = append(drivers, driver)
}

func RefreshUrl(urls []string) error {
	for _, driver := range drivers {
		err := driver.RefreshUrl(urls)
		if err != nil {
			return err
		}
	}

	return nil
}

func RefreshPath(paths []string) error {
	for _, driver := range drivers {
		err := driver.RefreshPath(paths)
		if err != nil {
			return err
		}
	}

	return nil
}

func GetUsage(domain string, startTime, endTime carbon.Carbon) (uint, error) {
	var total uint
	for _, driver := range drivers {
		usage, err := driver.GetUsage(domain, startTime, endTime)
		if err != nil {
			return 0, err
		}
		total += usage
	}

	return total, nil
}

func driverInUse(driver string) bool {
	config := facades.Config().GetString("cdn.driver")
	configs := strings.Split(config, ",")
	if slices.Contains(configs, driver) {
		return true
	}

	return false
}
