package ratelimit

import "time"

type Limit struct {
	Calls       []time.Time
	Constraints Constraints
	Max         int
}

func (limit Limit) Ready() bool {
	return limit.Constraints.Ready(limit.Calls)
}

func (limit *Limit) Use() {
	limit.Calls = append(limit.Calls, time.Now())
	if len(limit.Calls) > limit.Max {
		limit.Calls = limit.Calls[len(limit.Calls)-limit.Max:]
	}
}

func (limit *Limit) WithLimit(requests int, duration time.Duration) *Limit {
	limit.addConstraint(Constraint{Requests: requests, Duration: duration})
	return limit
}

func (limit *Limit) addConstraint(constraint Constraint) {
	if constraint.Requests > limit.Max {
		limit.Max = constraint.Requests
	}
	limit.Constraints = append(limit.Constraints, constraint)
}

func New() *Limit { return &Limit{} }
