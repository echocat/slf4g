package support

import "time"

func PString(v string) *string {
	return &v
}

func PTime(v time.Time) *time.Time {
	return &v
}
