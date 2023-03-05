package helpers

import "time"
import "fmt"

func RelativeDateText(nextDate time.Time) string {
	var relDate = ""
	if Now().Format("2006-02-01") == nextDate.Format("2006-02-01") {
		relDate = nextDate.Format("15:04")
	} else if Now().Add(24*time.Hour).Format("2006-02-01") == nextDate.Format("2006-02-01") {
		relDate = fmt.Sprintf("%s завтра", nextDate.Format("15:04"))
	} else if Now().Add(-24*time.Hour).Format("2006-02-01") == nextDate.Format("2006-02-01") {
		relDate = fmt.Sprintf("%s вчера", nextDate.Format("15:04"))
	} else {
		relDate = nextDate.Format("02.01 15:04")
	}
	return relDate
}
