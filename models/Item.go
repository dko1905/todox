package models

import (
	"fmt"
	"hash"
	"time"

	"github.com/google/uuid"
)

type Task struct {
	ID        uuid.UUID `json:"id"`
	Title     string    `json:"title"`
	Done      bool      `json:"done"`
	CreatedTS time.Time `json:"createdTS"`
}

func (t Task) HashWrite32(h *hash.Hash32) {
	(*h).Write([]byte(t.ID.String()))
	(*h).Write([]byte(t.Title))
	if t.Done {
		(*h).Write([]byte{1})
	} else {
		(*h).Write([]byte{0})
	}
	(*h).Write([]byte(fmt.Sprintf("%d", t.CreatedTS.Unix())))
}
