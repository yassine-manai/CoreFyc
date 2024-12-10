package functions

import "time"

func GetFormatedLocalTime() string {
	var currentTime = time.Now().UTC()
	var FormattedTime = currentTime.Format("2006-01-02 15:04:05")

	return FormattedTime

}
