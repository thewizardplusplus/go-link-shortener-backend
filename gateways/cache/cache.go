package cache

import (
	"github.com/thewizardplusplus/go-link-shortener/entities"
)

// KeyExtractor ...
type KeyExtractor func(link entities.Link) string
