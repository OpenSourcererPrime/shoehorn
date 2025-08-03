package config

import (
	"errors"
	"fmt"
)

var ErrorParseConfig = errors.New("failed to parse config")

// ErrorInvalidStrategy is returned when an invalid strategy is specified
type ErrorInvalidStrategy struct {
	Strategy string
	Name     string
}

func (e *ErrorInvalidStrategy) Error() string {
	return fmt.Sprintf("invalid strategy '%s' for config '%s'. Must be 'append' or 'template'", e.Strategy, e.Name)
}

// ErrorMissingTemplate is returned when a template path is required but not provided
type ErrorMissingTemplate struct {
	Name string
}

func (e *ErrorMissingTemplate) Error() string {
	return fmt.Sprintf("template path must be provided when strategy is 'template' for '%s'", e.Name)
}

// ErrorInvalidReloadMethod is returned when an invalid reload method is specified
type ErrorInvalidReloadMethod struct {
	Method string
}

func (e *ErrorInvalidReloadMethod) Error() string {
	return fmt.Sprintf("invalid reload method '%s'. Must be 'restart' or 'signal'", e.Method)
}

// ErrorMissingSignal is returned when a signal is required but not provided
type ErrorMissingSignal struct{}

func (e *ErrorMissingSignal) Error() string {
	return "signal must be provided when reload method is 'signal'"
}
