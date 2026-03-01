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

package sbmodels

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// JSONValue is a generic wrapper for JSON columns that may be null. V == nil represents null.
// It implements sql.Scanner, driver.Valuer, json.Marshaler, and json.Unmarshaler so that
// DB NULL and JSON "null" are handled correctly for serialization and deserialization.
type JSONValue[T any] struct {
	V *T
}

// Scan implements sql.Scanner. Accepts nil or []byte (JSON). Sets j.V to nil for null.
func (j *JSONValue[T]) Scan(value any) error {
	if value == nil {
		j.V = nil
		return nil
	}
	data, ok := value.([]byte)
	if !ok {
		return errors.New("json_value: scan source is not []byte")
	}
	if len(data) == 0 {
		j.V = nil
		return nil
	}
	var t T
	if err := json.Unmarshal(data, &t); err != nil {
		return err
	}
	j.V = &t
	return nil
}

// Value implements driver.Valuer. Returns nil when j.V is nil, otherwise JSON bytes.
func (j JSONValue[T]) Value() (driver.Value, error) {
	if j.V == nil {
		return nil, nil
	}
	return json.Marshal(j.V)
}

// MarshalJSON implements json.Marshaler. Outputs "null" when V is nil.
func (j JSONValue[T]) MarshalJSON() ([]byte, error) {
	if j.V == nil {
		return []byte("null"), nil
	}
	return json.Marshal(j.V)
}

// UnmarshalJSON implements json.Unmarshaler. Sets V to nil when data is "null".
func (j *JSONValue[T]) UnmarshalJSON(data []byte) error {
	if len(data) == 0 {
		j.V = nil
		return nil
	}
	if isJSONNull(data) {
		j.V = nil
		return nil
	}
	var t T
	if err := json.Unmarshal(data, &t); err != nil {
		return err
	}
	j.V = &t
	return nil
}

func isJSONNull(data []byte) bool {
	return len(data) == 4 && data[0] == 'n' && data[1] == 'u' && data[2] == 'l' && data[3] == 'l'
}

// Ptr returns the inner *T. Returns nil when the value is null.
func (j JSONValue[T]) Ptr() *T {
	return j.V
}

// JSONValueOf returns a JSONValue holding v. Use when writing a non-null value.
func JSONValueOf[T any](v *T) JSONValue[T] {
	return JSONValue[T]{V: v}
}
