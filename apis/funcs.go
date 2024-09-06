/*
funcs.go contains arbitrary functions, unattached to any struct.
*/
package apis

import (
	"fmt"
	"time"
)

// getDateNowString returns today's date in the format of YYYY-MM-DD.
func GetDateNowString() string {
	y, m, d := time.Now().Date()

	currentDate := fmt.Sprintf("%d-", y)

	if m < 10 {
		currentDate += fmt.Sprintf("0%d-", m)
	} else {
		currentDate += fmt.Sprintf("%d-", m)
	}

	if d < 10 {
		currentDate += fmt.Sprintf("0%d", d)
	} else {
		currentDate += fmt.Sprintf("%d", d)
	}

	return currentDate
}
