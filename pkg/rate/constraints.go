package rate

import (
	"math"
	"time"
)

type Constraints []Constraint

func (constraints Constraints) Ready(hits []time.Time) bool {
	for _, constraint := range constraints {
		ready := constraint.Ready(hits)
		if !ready {
			return false
		}
	}
	return true
}

func (constraints Constraints) When(hits []time.Time) time.Duration {
	best := time.Duration(math.MaxInt64)
	for _, constraint := range constraints {
		when := constraint.When(hits)
		if when < best {
			best = when
		}
	}
	return best
}

type Constraint struct {
	Requests int
	Duration time.Duration
}

func (constraint Constraint) Ready(hits []time.Time) bool {
	if len(hits) < constraint.Requests {
		return true
	}
	hit := hits[len(hits)-constraint.Requests]
	return time.Since(hit) > constraint.Duration
}

func (constraint Constraint) When(hits []time.Time) time.Duration {
	if constraint.Ready(hits) {
		return time.Duration(0)
	}
	hit := hits[len(hits)-constraint.Requests]
	return constraint.Duration - time.Since(hit)
}
