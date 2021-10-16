package main

import (
	"fmt"
	"time"
)

func makeKey(month string, year int) string {
	return fmt.Sprintf("%s %d", month, year)
}

func chartsFromFileChanges(fileChanges FileChanges) []Chart {
	var charts []Chart
	for fileName, times := range fileChanges {
		var timelinePoints []TimelinePoint
		beginningTS := earliestTimestamp(times)
		latestTS := latestTimestamp(times)

		timelinePoints = TimelinePoints(beginningTS, latestTS, times)
		charts = append(charts, Chart{
			Filename:       fileName,
			TimelinePoints: timelinePoints,
		})
	}

	return charts
}

func TimelinePoints(beginningTS time.Time, latestTS time.Time, times []time.Time) []TimelinePoint {
	var timelinePoints []TimelinePoint
	for year := beginningTS.Year(); year <= latestTS.Year(); year++ {
		for month := beginningTS.Month(); month <= 12; month++ {
			// get amount of timestamps in current month and year
			amountOfChanges := 0
			for _, t := range times {
				if t.Year() == year && t.Month() == month {
					amountOfChanges++
				}
			}

			timelinePoints = append(timelinePoints, TimelinePoint{
				Key:             makeKey(month.String(), year),
				AmountOfChanges: amountOfChanges,
			})

			// end of calculation
			if year == latestTS.Year() && month == latestTS.Month() {
				return timelinePoints
			}

			if month == time.December && year != latestTS.Year() {
				beginningTS = time.Date(year, time.January, 1,
					0, 0, 0, 0, &time.Location{})
			}
		}
	}

	return timelinePoints
}

func earliestTimestamp(tss []time.Time) time.Time {
	var ts = time.Now()
	for i := range tss {
		if tss[i].Before(ts) {
			ts = tss[i]
		}
	}

	return ts
}

func latestTimestamp(tss []time.Time) time.Time {
	var ts time.Time
	for i := range tss {
		if tss[i].After(ts) {
			ts = tss[i]
		}
	}

	return ts
}
