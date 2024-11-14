// Copyright 2024 Daytona Platforms Inc.
// SPDX-License-Identifier: Apache-2.0

package workspace

import "errors"

type Store interface {
	List(filter *Filter) ([]*Job, error)
	Find(filter *Filter) (*Job, error)
	Save(job *Job) error
	Delete(job *Job) error
}

type Filter struct {
	Id      *string
	States  *[]JobState
	Actions *[]JobAction
}

func (f *Filter) StatesToInterface() []interface{} {
	args := make([]interface{}, len(*f.States))
	for i, v := range *f.States {
		args[i] = v
	}
	return args
}

func (f *Filter) ActionsToInterface() []interface{} {
	args := make([]interface{}, len(*f.Actions))
	for i, v := range *f.Actions {
		args[i] = v
	}
	return args
}

var (
	ErrJobNotFound = errors.New("workspace job not found")
)

func IsJobNotFound(err error) bool {
	return err.Error() == ErrJobNotFound.Error()
}
