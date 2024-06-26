package chronos

import "time"

type Chronos interface {
	Now() time.Time
}
