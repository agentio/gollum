package common

import "time"

const ISO8601 = "2006-01-02T15:04:05.000Z"

func Now() string {
	return time.Now().UTC().Format(ISO8601)
}
