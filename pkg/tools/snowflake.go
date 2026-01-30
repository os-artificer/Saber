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

package tools

import (
	"fmt"
	"sync"
	"time"

	"os-artificer/saber/pkg/gerrors"
)

const (
	timeBits              = 41  // timestamp bit count
	machineIDBits         = 10  // machine id bit count
	sequenceBits          = 12  // serial number bit count
	maxRollbackTimeMillis = 100 // Max rollback time
	maxMachineID          = -1 ^ (-1 << machineIDBits)
	maxSequence           = -1 ^ (-1 << sequenceBits)
	timeShift             = machineIDBits + sequenceBits
	machineShift          = sequenceBits
)

type Snowflake struct {
	mu            sync.Mutex
	epoch         uint64 // Customize era time(in milliseconds).
	machineID     uint64 // Machine ID.
	sequence      uint64 // Serial number.
	lastTimestamp uint64 // The timestamp when the ID was last generated.
	timeBackward  bool   // Clock rewind indicator.
}

// NewSnowflake create new snowflake object
func NewSnowflake(machineID uint64, epoch time.Time) (*Snowflake, error) {
	if machineID > maxMachineID {
		return nil, gerrors.New(gerrors.InvalidParameter, "machine-id out of range")
	}

	return &Snowflake{
		epoch:     uint64(epoch.UnixMilli()),
		machineID: machineID,
	}, nil
}

func (s *Snowflake) NextID() (uint64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	current := uint64(time.Now().UnixMilli()) - s.epoch

	if current < s.lastTimestamp {
		s.timeBackward = true
		offset := s.lastTimestamp - current

		if offset > maxRollbackTimeMillis {
			// Allow for clock rollbacks within 100 milliseconds.
			return 0, gerrors.New(gerrors.Failure, "clock moved backwards too much")
		}

		// Wait for the clock to catch up.
		time.Sleep(time.Duration(offset) * time.Millisecond)
		current = uint64(time.Now().UnixMilli()) - s.epoch

		if current < s.lastTimestamp {
			return 0, gerrors.New(gerrors.Failure, "clock moved backwards after waiting")
		}
	}

	if current == s.lastTimestamp {
		s.sequence = (s.sequence + 1) & maxSequence
		if s.sequence == 0 {
			// The current millisecond sequence number has been exhausted.
			// Waitint for the nexe millisecond.
			current = uint64(s.waitNextMillis())
		}
	} else {
		fmt.Println("current:", current)
		s.sequence = current & maxSequence
	}

	s.lastTimestamp = current
	return (current << timeShift) | (s.machineID << machineShift) | s.sequence, nil
}

func (s *Snowflake) waitNextMillis() uint64 {
	current := uint64(time.Now().UnixMilli()) - s.epoch

	for current <= s.lastTimestamp {
		time.Sleep(100 * time.Microsecond)
		current = uint64(time.Now().UnixMilli()) - s.epoch
	}

	return current
}

func (s *Snowflake) ParseID(id uint64) (timestamp, machineID, sequence uint64) {
	timestamp = (id >> timeShift) + uint64(s.epoch)
	machineID = (id >> machineShift) & maxMachineID
	sequence = id & maxSequence
	return
}
