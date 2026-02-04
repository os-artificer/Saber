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
	"errors"
	"strings"
	"sync"
	"testing"
)

func TestGo_RecoversPanic(t *testing.T) {
	done := make(chan struct{})
	Go(func() {
		defer close(done)
		panic("test panic")
	})
	<-done
	// If we get here, process did not exit — panic was recovered
}

func TestGoWithRecover_CustomHandler(t *testing.T) {
	var got any
	var wg sync.WaitGroup
	wg.Add(1)
	GoWithRecover(func() {
		defer wg.Done()
		panic("custom")
	}, func(r any) {
		got = r
	})
	wg.Wait()
	if got != "custom" {
		t.Errorf("handler got %v, want \"custom\"", got)
	}
}

func TestGoWithRecover_NilHandlerUsesDefault(t *testing.T) {
	done := make(chan struct{})
	GoWithRecover(func() {
		defer close(done)
		panic("nil handler")
	}, nil)
	<-done
}

func TestRecoverToError(t *testing.T) {
	var err error
	done := make(chan struct{})
	go func() {
		defer close(done)
		defer RecoverToError(&err)
		panic("recover to err")
	}()
	<-done
	if err == nil {
		t.Fatal("expected err to be set after panic")
	}
	if !strings.Contains(err.Error(), "recover to err") {
		t.Errorf("error should mention panic: %v", err)
	}
}

func TestRecoverToError_PreservesOriginalError(t *testing.T) {
	origErr := errors.New("original")
	var err error = origErr
	done := make(chan struct{})
	go func() {
		defer close(done)
		defer RecoverToError(&err)
		panic("panic after err")
	}()
	<-done
	if err == nil {
		t.Fatal("expected err to be set")
	}
	if !strings.Contains(err.Error(), "original") || !strings.Contains(err.Error(), "panic after err") {
		t.Errorf("error should mention both original and panic: %v", err)
	}
}

func TestRecoverWithHandler(t *testing.T) {
	var got any
	done := make(chan struct{})
	go func() {
		defer close(done)
		defer RecoverWithHandler(func(r any) { got = r })()
		panic("handler")
	}()
	<-done
	if got != "handler" {
		t.Errorf("got %v, want \"handler\"", got)
	}
}

func TestRecoverNoop(t *testing.T) {
	done := make(chan struct{})
	go func() {
		defer close(done)
		defer RecoverNoop()
		panic("noop")
	}()
	<-done
}

func TestRecoverAndStackHandler(t *testing.T) {
	var msg, stack string
	done := make(chan struct{})
	go func() {
		defer close(done)
		defer RecoverAndStackHandler(func(m, s string) { msg, stack = m, s })()
		panic("stack test")
	}()
	<-done
	if msg != "stack test" || stack == "" {
		t.Errorf("RecoverAndStackHandler: msg=%q stack empty=%v", msg, stack == "")
	}
}

func TestRecoverWithDebugStackHandler(t *testing.T) {
	var msg, full string
	done := make(chan struct{})
	go func() {
		defer close(done)
		defer RecoverWithDebugStackHandler(func(m, s string) { msg, full = m, s })()
		panic("debug stack")
	}()
	<-done
	if msg != "debug stack" || full == "" {
		t.Errorf("RecoverWithDebugStackHandler: msg=%q full empty=%v", msg, full == "")
	}
}

