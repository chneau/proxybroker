package rate

import "time"

type Limit struct {
	Calls       []time.Time
	Constraints Constraints
	Max         int
}

func (limit Limit) Ready() bool { return limit.Constraints.Ready(limit.Calls) }

func (limit Limit) When() time.Duration { return limit.Constraints.When(limit.Calls) }

func (limit Limit) Clone() *Limit {
	copy := limit
	copy.Calls = append(limit.Calls[:0:0], limit.Calls...)
	return &copy
}

func (limit *Limit) Use() {
	limit.Calls = append(limit.Calls, time.Now())
	if len(limit.Calls) > limit.Max {
		limit.Calls = limit.Calls[len(limit.Calls)-limit.Max:]
	}
}

func (limit *Limit) WithConstraint(requests int, duration time.Duration) *Limit {
	limit.addConstraint(Constraint{Requests: requests, Duration: duration})
	return limit
}

func (limit *Limit) addConstraint(constraint Constraint) {
	if constraint.Requests > limit.Max {
		limit.Max = constraint.Requests
	}
	limit.Constraints = append(limit.Constraints, constraint)
}

func NewLimit() *Limit { return &Limit{} }
