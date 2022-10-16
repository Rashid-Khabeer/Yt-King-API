package base

import (
	"time"
)

type Timestamped struct {
	CreatedAt time.Time `bson:"createdAt" json:"createdAt,omitempty"`
	UpdatedAt time.Time `bson:"updatedAt" json:"updatedAt,omitempty"`
}
