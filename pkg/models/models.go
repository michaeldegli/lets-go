package models

import (
	"time"
)

// Snippet ...
type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

type Snippets []*Snippet
