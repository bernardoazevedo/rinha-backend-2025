package dates

import (
	"fmt"
	"time"
)

func FormatYearMonthDay(data time.Time) string {
	year := data.Year()
	month := data.Month()
	day := data.Day()

	return fmt.Sprintf("%d-%d-%d", year, month, day)
}

func YearMonthDay() string {
	return FormatYearMonthDay(time.Now())
}

func HourMinuteSecond() string {
	hour, minute, second := time.Now().Clock()
	return fmt.Sprintf("%d:%d:%d", hour, minute, second)
}