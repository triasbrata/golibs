package chronos

import "time"

type t_cronos struct {
}

// Now implements Chronos.
func (*t_cronos) Now() time.Time {
	return time.Now()
}

func New() Chronos {
	return &t_cronos{}
}
