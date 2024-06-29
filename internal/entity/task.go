package entity

import "time"

type Task interface {
	TaskContent() string
	TaskType() string
	Estimate() time.Duration
}
