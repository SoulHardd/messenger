package domain

import "time"

type Cursor struct {
	Time time.Time
	Id   int64
}
