package util

import "time"

func TimeToString(t time.Time) string {
	return t.UTC().Format(time.RFC3339)
}

func StringToTime(s string) (time.Time, error) {
	return time.Parse(time.RFC3339, s)
}

func TimeNow() time.Time {
	return time.Now().UTC()
}