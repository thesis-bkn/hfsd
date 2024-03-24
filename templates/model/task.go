package model

import "time"

//go:generate go-enum

// ENUM(
//
//	inference
//	sample
//	finetune
//
// )
type TaskType int

type Task struct {
	ID            string
	SourceModelID string
	OutputModelID string
	TaskType      TaskType
	CreatedAt     time.Time
	HandledAt     *time.Time
	FinishedAt    *time.Time
	HumanPrefs    []*HumanPref
}

type HumanPref struct {
	Pref  bool
	Order int
}
