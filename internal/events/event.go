package events

import "time"

type Event interface {
    EventType() string
    OccurredAt() time.Time
}


