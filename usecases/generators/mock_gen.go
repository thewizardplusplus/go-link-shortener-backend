package generators

// nolint: lll
import (
	"github.com/thewizardplusplus/go-link-shortener-backend/usecases/generators/counters"
)

//go:generate mockery -name=DistributedCounter -inpkg -case=underscore -testonly

// DistributedCounter ...
//
// It's used only for mock generating.
type DistributedCounter interface {
	counters.DistributedCounter
}
