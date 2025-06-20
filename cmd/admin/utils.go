package main

import (
	"strconv"
	"time"
)

// Constants for seconds
const (
	oneMinute        = 60
	fiveMinutes      = 300
	fifteenMinutes   = 900
	thirtyMinutes    = 1800
	fortyfiveMinutes = 2500
	oneHour          = 3600
	threeHours       = 10800
	sixHours         = 21600
	eightHours       = 28800
	twelveHours      = 43200
	fifteenHours     = 54000
	twentyHours      = 72000
	oneDay           = 86400
	twoDays          = 172800
	sevenDays        = 604800
	fifteenDays      = 1296000
)

// Helper to get a string based on the difference of two times
func stringifyTime(seconds int) string {
	var timeStr string
	w := make(map[int]string)
	w[oneDay] = "day"
	w[oneHour] = "hour"
	w[oneMinute] = "minute"
	// Ordering the values will prevent bad values
	var ww [3]int
	ww[0] = oneDay
	ww[1] = oneHour
	ww[2] = oneMinute
	for _, v := range ww {
		if seconds >= v {
			d := seconds / v
			dStr := strconv.Itoa(d)
			timeStr = dStr + " " + w[v]
			if d > 1 {
				timeStr += "s"
			}
			break
		}
	}
	return timeStr
}

// Helper to format past times only returning one value (minute, hour, day)
func pastTimeAgo(t time.Time) string {
	if t.IsZero() {
		return "Never"
	}
	now := time.Now()
	seconds := int(now.Sub(t).Seconds())
	if seconds < 2 {
		return "Just Now"
	}
	if seconds < oneMinute {
		return strconv.Itoa(seconds) + " seconds ago"
	}
	if seconds > fifteenDays {
		return "Since " + t.Format("Mon Jan 02 15:04:05 MST 2006")
	}
	return stringifyTime(seconds) + " ago"
}

