package time

import "time"

func StartOfMonth(date time.Time) time.Time {
	return time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, date.Location())
}

func EndOfMonth(date time.Time) time.Time {
	firstDayOfNextMonth := StartOfMonth(date).AddDate(0, 1, 0)
	return firstDayOfNextMonth.Add(-time.Second)
}

func StartOfDayOfWeek(date time.Time) time.Time {
	daysSinceSunday := int(date.Weekday())
	return date.AddDate(0, 0, -daysSinceSunday)
}

func EndOfDayOfWeek(date time.Time) time.Time {
	daysUntilSaturday := 6 - int(date.Weekday())
	return date.AddDate(0, 0, daysUntilSaturday)
}

func StartAndEndOfWeeksOfMonth(year, month int) []struct{ Start, End time.Time } {
	startOfMonth := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	weeks := make([]struct{ Start, End time.Time }, 0)

	for current := startOfMonth; current.Month() == time.Month(month); current = current.AddDate(0, 0, 7) {
		startOfWeek := StartOfDayOfWeek(current)
		endOfWeek := EndOfDayOfWeek(current)

		if endOfWeek.Month() != time.Month(month) {
			endOfWeek = EndOfMonth(current)
		}
		weeks = append(weeks, struct{ Start, End time.Time }{startOfWeek, endOfWeek})
	}

	return weeks
}

func WeekNumberInMonth(date time.Time) int {
	startOfMonth := StartOfMonth(date)
	_, week := date.ISOWeek()
	_, startWeek := startOfMonth.ISOWeek()
	return week - startWeek + 1
}

func StartOfYear(date time.Time) time.Time {
	return time.Date(date.Year(), time.January, 1, 0, 0, 0, 0, date.Location())
}

func EndOfYear(date time.Time) time.Time {
	startOfNextYear := StartOfYear(date).AddDate(1, 0, 0)
	return startOfNextYear.Add(-time.Second)
}

func StartOfQuarter(date time.Time) time.Time {
	// you can directly use 0, 1, 2, 3 quarter
	quarter := (int(date.Month()) - 1) / 3
	startMonth := time.Month(quarter*3 + 1)
	return time.Date(date.Year(), startMonth, 1, 0, 0, 0, 0, date.Location())
}

func EndOfQuarter(date time.Time) time.Time {
	startOfNextQuarter := StartOfQuarter(date).AddDate(0, 3, 0)
	return startOfNextQuarter.Add(-time.Second)
}

func CurrentWeekRange(timeZone string) (startOfWeek, endOfWeek time.Time) {
	loc, _ := time.LoadLocation(timeZone)

	now := time.Now().In(loc)
	startOfWeek = StartOfDayOfWeek(now)
	endOfWeek = EndOfDayOfWeek(now)

	return startOfWeek, endOfWeek
}

func DurationBetween(start, end time.Time) time.Duration {
	return end.Sub(start)
}

func GetDatesForDayOfWeek(year, month int, day time.Weekday) []time.Time {
	var dates []time.Time

	firstDayOfMonth := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	diff := int(day) - int(firstDayOfMonth.Weekday())
	if diff < 0 {
		diff += 7
	}

	firstDay := firstDayOfMonth.AddDate(0, 0, diff)
	for current := firstDay; current.Month() == time.Month(month); current = current.AddDate(0, 0, 7) {
		dates = append(dates, current)
	}

	return dates
}

func AddBusinessDays(startDate time.Time, daysToAdd int) time.Time {
	currentDate := startDate
	for i := 0; i < daysToAdd; {
		currentDate = currentDate.AddDate(0, 0, 1)
		if currentDate.Weekday() != time.Saturday && currentDate.Weekday() != time.Sunday {
			i++
		}
	}
	return currentDate
}

func NumberMonthInRange(start, end time.Time) int {
	return int(end.Month() - start.Month())
}
