package db

import (
	"time"
)

type Statistics struct {
	Queues []QueueStatistics `json:"queues"`
}

type QueueStatistics struct {
	Name        string             `json:"name"`
	Scheduled   int                `json:"scheduled"`
	Initialized int                `json:"initialized"`
	Failed      int                `json:"failed"`
	Completed   int                `json:"completed"`
	Workers     []WorkerStatistics `json:"workers"`
}
type WorkerStatistics struct {
	ID         string    `json:"id"`
	LastSeenAt time.Time `json:"last_seen_at"`
}
