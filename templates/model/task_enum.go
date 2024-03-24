// Code generated by go-enum DO NOT EDIT.
// Version: 0.6.0
// Revision: 919e61c0174b91303753ee3898569a01abb32c97
// Build Date: 2023-12-18T15:54:43Z
// Built By: goreleaser

package model

import (
	"errors"
	"fmt"
)

const (
	// TaskTypeInference is a TaskType of type Inference.
	TaskTypeInference TaskType = iota
	// TaskTypeSample is a TaskType of type Sample.
	TaskTypeSample
	// TaskTypeFinetune is a TaskType of type Finetune.
	TaskTypeFinetune
)

var ErrInvalidTaskType = errors.New("not a valid TaskType")

const _TaskTypeName = "inferencesamplefinetune"

var _TaskTypeMap = map[TaskType]string{
	TaskTypeInference: _TaskTypeName[0:9],
	TaskTypeSample:    _TaskTypeName[9:15],
	TaskTypeFinetune:  _TaskTypeName[15:23],
}

// String implements the Stringer interface.
func (x TaskType) String() string {
	if str, ok := _TaskTypeMap[x]; ok {
		return str
	}
	return fmt.Sprintf("TaskType(%d)", x)
}

// IsValid provides a quick way to determine if the typed value is
// part of the allowed enumerated values
func (x TaskType) IsValid() bool {
	_, ok := _TaskTypeMap[x]
	return ok
}

var _TaskTypeValue = map[string]TaskType{
	_TaskTypeName[0:9]:   TaskTypeInference,
	_TaskTypeName[9:15]:  TaskTypeSample,
	_TaskTypeName[15:23]: TaskTypeFinetune,
}

// ParseTaskType attempts to convert a string to a TaskType.
func ParseTaskType(name string) (TaskType, error) {
	if x, ok := _TaskTypeValue[name]; ok {
		return x, nil
	}
	return TaskType(0), fmt.Errorf("%s is %w", name, ErrInvalidTaskType)
}
