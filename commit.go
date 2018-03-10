package hatena

import (
	"time"

	"github.com/podhmo/commithistory"
)

type Commit = commithistory.Commit

// NewCommit creates and initializes a new Commit object.
func NewCommit(id string, alias string, action string) *Commit {
	now := time.Now()
	return &Commit{
		ID:        id,
		CreatedAt: now,
		Alias:     alias,
		Action:    action,
	}
}
