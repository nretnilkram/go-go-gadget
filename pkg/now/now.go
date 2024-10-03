package now

import (
	"strconv"
	"time"
)

// Return the String format of Date Time as requested
func ShowDateTime(format string, showTime bool) (dateTime string) {
	current_time := time.Now()
	var year = current_time.Year()
	var month = int(current_time.Month())
	var day = current_time.Day()
	var hour = current_time.Hour()
	var minute = current_time.Minute()
	var separator = "-"

	switch format {
	case "colon":
		separator = ":"
	case "dash":
		separator = "-"
	case "dot":
		separator = "."
	case "slash":
		separator = "/"
	default:
		return current_time.String()
	}

	dateTime = strconv.Itoa(year) + separator + strconv.Itoa(month) + separator + strconv.Itoa(day)
	if showTime {
		dateTime = dateTime + " " + strconv.Itoa(hour) + ":" + strconv.Itoa(minute)
	}

	return dateTime
}
