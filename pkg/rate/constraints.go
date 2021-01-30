package rate

import "time"

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
