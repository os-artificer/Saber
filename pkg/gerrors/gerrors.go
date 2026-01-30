/**
 * Copyright 2025 saber authors.
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

import "fmt"

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
}

// New create a new GError object
func New(c Code, msg string) *GError {
	return &GError{code: c, message: msg}
}

func Newf(c Code, format string, args ...interface{}) *GError {
	msg := fmt.Sprintf(format, args...)
	return &GError{code: c, message: msg}
}

func NewE(c Code, err error) *GError {
	if err == nil {
		return nil
	}
	return &GError{code: c, message: err.Error()}
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
