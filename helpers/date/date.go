package date

import "github.com/golang-module/carbon/v2"

func GetNewYorkZone() carbon.Carbon {
	return carbon.SetTimezone("America/New_York")
}
