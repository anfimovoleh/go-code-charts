package main

import "time"

type FileChanges map[string][]time.Time

type TimelinePoint struct {
	Key             string
	AmountOfChanges int
}

type Chart struct {
	Filename       string
	TimelinePoints []TimelinePoint
}
