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

package tools

import (
	"fmt"
	"runtime"
	"runtime/debug"

	"os-artificer/saber/pkg/logger"
)

// PanicHandler is called when a panic is recovered. It receives the panic value.
// Handlers must not panic; they run after recover and must not crash the process.
type PanicHandler func(panicValue any)

// defaultPanicHandler logs the panic and stack trace.
func defaultPanicHandler(r any) {
	const size = 64 << 10
	stack := make([]byte, size)
	stack = stack[:runtime.Stack(stack, false)]
	logger.Errorf("goroutine panic recovered (process not exiting): panic=%v\nstack:\n%s", r, string(stack))
}

// Go runs fn in a new goroutine. If fn panics, the panic is recovered and logged
// (with stack trace); the process does not exit. No retry is performed.
// Use this to start goroutines that must not crash the process on panic.
func Go(fn func()) {
	go runWithRecover(fn, defaultPanicHandler)
}

// GoWithRecover runs fn in a new goroutine. If fn panics, onPanic is invoked with
// the panic value and the process does not exit. No retry. If onPanic is nil,
// the default handler (log + stack) is used.
func GoWithRecover(fn func(), onPanic PanicHandler) {
	handler := onPanic
	if handler == nil {
		handler = defaultPanicHandler
	}
	go runWithRecover(fn, handler)
}

// RecoverToError is meant to be used with defer inside a goroutine that returns an error.
// If a panic occurs, it sets *err to an error that includes the panic value and call stack,
// and does not re-panic.
//
//	errPtr := new(error)
//	go func() {
//	    defer tools.RecoverToError(errPtr)
//	    *errPtr = doWork()
//	}()
func RecoverToError(err *error) {
	if r := recover(); r != nil {
		const size = 64 << 10
		stack := make([]byte, size)
		stack = stack[:runtime.Stack(stack, false)]
		*err = fmt.Errorf("recovered from panic: panic=%v (original err=%v)\nstack:\n%s", r, *err, string(stack))
	}
}

// RecoverWithHandler returns a function that should be used with defer. When the
// surrounding function panics, it recovers, calls handler with the panic value,
// and does not re-panic. Use when you need custom recovery inside a goroutine
// without starting it via Go/GoWithRecover.
//
//	go func() {
//	    defer tools.RecoverWithHandler(func(r any) { log.Printf("panic: %v", r) })()
//	    riskyWork()
//	}()
func RecoverWithHandler(handler PanicHandler) func() {
	return func() {
		if r := recover(); r != nil {
			if handler != nil {
				handler(r)
			} else {
				defaultPanicHandler(r)
			}
		}
	}
}

// RecoverNoop recovers a panic and discards it; process does not exit. Use for
// fire-and-forget goroutines where no logging or error reporting is needed.
func RecoverNoop() {
	_ = recover()
}

// RecoverAndStackHandler returns a function to be used with defer. When the
// surrounding function panics, it recovers, invokes handler with the panic message
// and stack trace, and does not re-panic. recover() must be called from the
// deferred function, so this returns the deferred func rather than being called inside one.
//
//	var msg, stack string
//	go func() {
//	    defer tools.RecoverAndStackHandler(func(m, s string) { msg, stack = m, s })()
//	    risky()
//	}()
func RecoverAndStackHandler(handler func(panicMsg, stackTrace string)) func() {
	return func() {
		if r := recover(); r != nil {
			const size = 64 << 10
			buf := make([]byte, size)
			buf = buf[:runtime.Stack(buf, false)]
			handler(fmt.Sprint(r), string(buf))
		}
	}
}

// RecoverWithDebugStackHandler is like RecoverAndStackHandler but passes
// debug.Stack() (full goroutine dump) to the handler. Useful for debugging.
func RecoverWithDebugStackHandler(handler func(panicMsg, fullStack string)) func() {
	return func() {
		if r := recover(); r != nil {
			handler(fmt.Sprint(r), string(debug.Stack()))
		}
	}
}

// runWithRecover executes fn and defers a recover that calls handler and does not re-panic.
func runWithRecover(fn func(), handler PanicHandler) {
	defer func() {
		if r := recover(); r != nil {
			handler(r)
		}
	}()
	fn()
}
