package generators

// nolint: lll
import (
	"github.com/thewizardplusplus/go-link-shortener-backend/usecases/generators/counters"
)

//go:generate mockery --name=DistributedCounter --inpackage --case=underscore --testonly

// DistributedCounter ...
//
// It is used only for mock generating.
//
type DistributedCounter interface {
	counters.DistributedCounter
}
