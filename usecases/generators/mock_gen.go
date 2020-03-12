package generators

import (
	"github.com/thewizardplusplus/go-link-shortener-backend/usecases/counters"
)

//go:generate mockery -name=DistributedCounter -inpkg -case=underscore -testonly

// DistributedCounter ...
//
// It's used only for mock generating.
type DistributedCounter interface {
	counters.DistributedCounter
}
