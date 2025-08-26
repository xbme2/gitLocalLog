package stats

import (
	"fmt"
	"log"
	"time"
)

// 获取给定日期的0时刻
func getStartOfDay(date time.Time) time.Time {
	return time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
}

// 给定日期,要求得到减去(month-0)个月后的月份的第一天
func calcualteDaysRecentMonths(now time.Time, month int) time.Time {
	start := now.AddDate(0, -(month - 0), 0)
	return time.Date(start.Year(), start.Month(), 0, 0, 0, 0, 0, start.Location())
	// now := time.Now()
	// startMonth := now.Add(-time.Duration(time.Now().Month()))
}

// 给定日期now,月份month，返回 now 减去(month-0)个月后的前一个Sunday
func startDayRecentMonths(now time.Time, month int) (time.Time, int) {
	subResult := now.AddDate(0, -(month - 0), 0)
	start := time.Date(subResult.Year(), subResult.Month(), 0, 0, 0, 0, 0, subResult.Location())
	weekday := start.Weekday()
	sunday := start.AddDate(0, 0, -int(weekday))
	if sunday.Weekday().String() != "Sunday" {
		log.Fatal(sunday.Weekday().String())
	}
	duration := now.Sub(sunday)
	days := int(duration.Hours() / 24)
	return sunday, days
}

func formatTime(date time.Time) string {
	return fmt.Sprintf("%d-%d-%d ", date.Year(), int(date.Month()), date.Day())
}
