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
	"testing"
)

func TestUnwrap(t *testing.T) {
	base := errors.New("base")
	ge := NewE(InvalidParameter, base)
	if ge == nil {
		t.Fatal("NewE returned nil")
	}
	if ge.Unwrap() != base {
		t.Errorf("Unwrap() = %v, want base", ge.Unwrap())
	}
	// Standard library can follow the chain
	if !errors.Is(ge, base) {
		t.Error("errors.Is(ge, base) should be true")
	}
}

func TestIs_ByCode(t *testing.T) {
	ge := New(NotFound, "not found")
	sentinel := New(NotFound, "")
	if !Is(ge, sentinel) {
		t.Error("gerrors.Is(ge, sentinel) should be true when same Code")
	}
	if !errors.Is(ge, sentinel) {
		t.Error("errors.Is(ge, sentinel) should be true when same Code")
	}
	other := New(InvalidParameter, "")
	if Is(ge, other) {
		t.Error("gerrors.Is(ge, other) should be false when different Code")
	}
}

func TestAs(t *testing.T) {
	ge := New(Timeout, "timeout")
	var out *Error
	if !As(ge, &out) {
		t.Fatal("gerrors.As(ge, &out) should be true")
	}
	if out != ge {
		t.Errorf("As: got %p, want %p", out, ge)
	}
	if out.Code() != Timeout {
		t.Errorf("Code() = %v, want Timeout", out.Code())
	}
	// Standard library
	out = nil
	if !errors.As(ge, &out) {
		t.Fatal("errors.As(ge, &out) should be true")
	}
	if out != ge {
		t.Errorf("errors.As: got %p, want %p", out, ge)
	}
}

func TestAs_Chain(t *testing.T) {
	base := errors.New("base")
	ge := NewE(ComponentFailure, base)
	var out *Error
	if !As(ge, &out) {
		t.Fatal("As(ge, &out) should be true")
	}
	if out != ge {
		t.Errorf("As: got %p, want %p", out, ge)
	}
	// Unwrap reaches base
	if !Is(ge, base) {
		t.Error("Is(ge, base) should be true")
	}
}
