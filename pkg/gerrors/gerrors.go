/**
 * Copyright 2025 Saber authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
**/

package gerrors

import (
	"errors"
	"fmt"
)

// Code global error code
type Code int

const (
	Unknown Code = iota - 2
	Failure
	Success
	Timeout
	InternalServer
	InvalidParameter
	InvalidConfig
	InvalidNetAddress
	Unimplemented
	NotFound
	AlreadyExists
	QueueFull
	ComponentFailure
)

// GError global error type
type GError struct {
	code    Code
	message string
	cause   error // optional; set by NewE for error chain (Unwrap)
}

// New create a new GError object
func New(c Code, msg string) *GError {
	return &GError{code: c, message: msg}
}

func Newf(c Code, format string, args ...any) *GError {
	msg := fmt.Sprintf(format, args...)
	return &GError{code: c, message: msg}
}

func NewE(c Code, err error) *GError {
	if err == nil {
		return nil
	}
	return &GError{code: c, message: err.Error(), cause: err}
}

// Code only return error code
func (g *GError) Code() Code {
	return g.code
}

// Message only return message
func (g *GError) Message() string {
	return g.message
}

// Error error interface method
func (g *GError) Error() string {
	return fmt.Sprintf("code:%d, errmsg:%s", g.code, g.message)
}

// Unwrap returns the underlying error if any, for errors.Is/errors.As compatibility.
func (g *GError) Unwrap() error {
	return g.cause
}

// Is reports whether the target is considered a match for this error.
// When target is *GError, matches by Code equality so that e.g.
// errors.Is(err, gerrors.New(gerrors.NotFound, "")) matches any GError with NotFound.
func (g *GError) Is(target error) bool {
	t, ok := target.(*GError)
	if !ok {
		return false
	}
	return g.code == t.code
}

// Is reports whether any error in err's chain matches target. Behavior is the same as errors.Is.
func Is(err, target error) bool {
	return errors.Is(err, target)
}

// As finds the first error in err's chain that matches target and assigns it. Behavior is the same as errors.As.
func As(err error, target any) bool {
	return errors.As(err, target)
}
