/*
funcs.go contains arbitrary functions, unattached to any struct.
*/
package apis

import (
	"fmt"
	"time"
)

// getDateNowString returns today's date in the format of YYYYMMDD.
// [0:4] for the year, [4:6] for the month, and [6:8] for the day.
func GetDateNowString() string {
	y, m, d := time.Now().Date()

	currentDate := fmt.Sprintf("%d", y)

	if m < 10 {
		currentDate += fmt.Sprintf("0%d", m)
	} else {
		currentDate += fmt.Sprintf("%d", m)
	}

	if d < 10 {
		currentDate += fmt.Sprintf("0%d", d)
	} else {
		currentDate += fmt.Sprintf("%d", d)
	}

	return currentDate
}
