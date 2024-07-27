package state

import (
	"github.com/pkg/errors"
	"github.com/samber/lo"
)

var (
	ErrStateUnknown = errors.New("unknown task state")
)

type State int

const (
	Pending State = iota
	Scheduled
	Running
	Completed
	Failed
)

func (s State) String() string {
	switch s {
	case Pending:
		return "Pending"
	case Scheduled:
		return "Scheduled"
	case Running:
		return "Running"
	case Completed:
		return "Completed"
	case Failed:
		return "Failed"
	default:
		return "Unknown"
	}
}

func Parse(s string) (State, error) {
	switch s {
	case "Pending":
		return Pending, nil
	case "Scheduled":
		return Scheduled, nil
	case "Running":
		return Running, nil
	case "Completed":
		return Completed, nil
	case "Failed":
		return Failed, nil
	default:
		return -1, ErrStateUnknown
	}
}

var stateTransitionTable = map[State][]State{
	Pending:   {Scheduled},
	Scheduled: {Scheduled, Running, Failed},
	Running:   {Running, Completed, Failed},
	Completed: {},
	Failed:    {},
}

func ValidStateTransition(src State, dst State) bool {
	return lo.Contains(stateTransitionTable[src], dst)
}
