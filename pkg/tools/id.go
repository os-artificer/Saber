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
	"crypto/sha256"
	"encoding/binary"
	"hash/fnv"
	"sync"
	"time"

	"github.com/denisbrodbeck/machineid"
	"github.com/google/uuid"
)

var (
	sf         *Snowflake
	sfmu       sync.Mutex
	id         string
	idErr      error
	idHash     uint64
	idOnce     ResetOnce
	idHashOnce ResetOnce
)

// Hash Calculate the hash value of any string.
// s    The target string value
// bits The number of digits to be retained.
func Hash(s string, bits uint) uint64 {
	h1 := sha256.Sum256([]byte(s))
	h2 := fnv.New64a()
	h2.Write(h1[:])

	fullHash := binary.BigEndian.Uint64(h1[:8]) ^ h2.Sum64()
	mask := uint64(1<<bits - 1)
	return fullHash & mask
}

// MachineID  return the  machine-id
func MachineID() (string, error) {
	idOnce.Do(func() error {
		id, idErr = machineid.ProtectedID("dbha-v2")
		return idErr
	})
	return id, idErr
}

// NewSequenceID create a new sequence-id
func NewSequenceID() (uint64, error) {
	sfmu.Lock()
	defer sfmu.Unlock()

	idHashOnce.Do(func() error {
		idHash = Hash(id, machineIDBits)
		return nil
	})

	if sf == nil {
		epoch, _ := time.Parse("2006-01-02", "2024-08-01")
		s, err := NewSnowflake(idHash, epoch)
		if err != nil {
			return 0, err
		}
		sf = s
	}

	id, err := sf.NextID()
	if err != nil {
		// After the time rewind, only one retry is allowed.
		s, e := NewSnowflake(idHash, time.Now())
		if e != nil {
			return 0, e
		}

		sf = s
		id, err = sf.NextID()
	}

	return id, err
}

// NewMessageID create a new message-id
func NewMessageID() string {
	id := uuid.New()
	return id.String()
}
