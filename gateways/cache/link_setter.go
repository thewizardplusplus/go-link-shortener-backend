package cache

import (
	"time"
)

// LinkSetter ...
type LinkSetter struct {
	KeyExtractor KeyExtractor
	Client       Client
	Expiration   time.Duration
}
