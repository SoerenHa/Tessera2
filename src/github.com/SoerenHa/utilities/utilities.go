/*

 */

package utilities

import (
	"fmt"
	"time"
	"github.com/kelvins/sunrisesunset"
)

type Times struct {
	Sunrise	time.Time
	Sunset	time.Time
}

// Combines a date and time string and returns a time object
func ParseTime(d, t string) time.Time {
	timeString := d + "T" + t + ":00+01:00"
	dateTime, _ := time.Parse(time.RFC3339, timeString)
	return dateTime
}
func ExecutionIsNow(simTime, sceneTime time.Time, timestep time.Duration, daily bool) bool {
	if daily {
		return simTime.Unix() >= sceneTime.Unix() &&
				simTime.Add(-timestep).Unix() % (24*60*60) <= sceneTime.Unix() % (24*60*60) &&
				simTime.Unix() % (24*60*60) > sceneTime.Unix() % (24*60*60)
	} else {
		return simTime.Add(-timestep).Unix() <= sceneTime.Unix() && sceneTime.Unix() < simTime.Unix()
	}
}

// Not in use
func GetTimes () Times {
	timeData := sunrisesunset.Parameters{
		Latitude:  54.774727032665766,
		Longitude: 9.447391927242279,
		UtcOffset: 1.0,
		Date:      time.Now(),
	}

	sunrise, sunset, Err := timeData.GetSunriseSunset()

	if Err != nil {
		fmt.Print(Err.Error())
	}

	return Times{
		Sunrise: sunrise,
		Sunset:	 sunset,
	}
}