package raft

import (
	"errors"
	"fmt"
)

// None is a placeholder node ID used when there is no leader.
const None uint64 = 0

var errNoLeader = errors.New("no leader")

// Possible values for StateType.
const (
	StateFollower StateType = iota
	StateCandidate
	StateLeader
)

// StateType represents the role of a node in a cluster.
type StateType uint64

var stamp = [...]string{
	"StateFolloer",
	"StateCandidate",
	"StaetLeader",
}

type Progress struct {
	Match, Next uint64
	Wait        int
}

func (pr *Progress) update(n uint64) {
	pr.waitReset()
	if pr.Match < n {
		pr.Match = n
	}
	if pr.Next < n+1 {
		pr.Next = n + 1
	}
}

func (pr *Progress) optimisticUpdate(n uint64) { pr.Next = n + 1 }

// maybyDecrTo returns false if the given to index comes from an out of order message.
// Otherwise it decreases the progress next indeex to min(rejected, last) and returns true.
func (pr *Progress) maybyDecrTo(rejected, last uint64) bool {
	pr.waitReset()
	if pr.Match != 0 {
		// the rejection must be stale if the progress has matched and "relected"
		// is smaller than "match"
		if rejected <= pr.Match {
			return false
		}
		// directly decrease next to match + 1
		pr.Next = pr.Match + 1
		return true
	}

	// the rejection must be stale if "rejected" does not match next -1
	if pr.Next-1 != rejected {
		return false
	}

	if pr.Next = min(rejected, last+1); pr.Next < 1 {
		pr.Next = 1
	}
	return true
}

func (pr *Progress) waitDecr(i int) {
	pr.Wait -= i
	if pr.Wait < 0 {
		pr.Wait = 0
	}
}

func (pr *Progress) waitSet(w int)    { pr.Wait = w }
func (pr *Progress) waitReset()       { pr.Wait = 0 }
func (pr *Progress) shouldWait() bool { return pr.Match == 0 && pr.Wait > 0 }

func (pr *Progress) String() string {
	return fmt.Sprintf("next = %d, match = %d, wait = %v", pr.Next, pr.Match, pr.Wait)
}
