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

package source

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// ConnMeta holds optional metadata for a tracked connection (e.g. client ID from first message).
type ConnMeta struct {
	ClientID    string
	RemoteAddr  string
	ConnectedAt time.Time
}

// ConnectionManager tracks long-lived gRPC client connections (PushData streams).
// It is safe for concurrent use and supports high concurrency.
type ConnectionManager struct {
	mu       sync.RWMutex
	conns    map[string]*ConnMeta
	maxConns int
	nextID   atomic.Uint64
}

// NewConnectionManager creates a manager that tracks active connections.
// If maxConns > 0, Register returns an error when the count would exceed maxConns.
func NewConnectionManager(maxConns int) *ConnectionManager {
	return &ConnectionManager{
		conns:    make(map[string]*ConnMeta),
		maxConns: maxConns,
	}
}

// Register adds a connection with the given id and optional metadata.
// Returns an error if maxConns is set and the active count would exceed it.
func (m *ConnectionManager) Register(id string, meta *ConnMeta) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.maxConns > 0 && len(m.conns) >= m.maxConns {
		return fmt.Errorf("max connections reached: %d", m.maxConns)
	}

	if _, exists := m.conns[id]; exists {
		return fmt.Errorf("connection already registered: %s", id)
	}

	if meta == nil {
		meta = &ConnMeta{ConnectedAt: time.Now()}
	} else if meta.ConnectedAt.IsZero() {
		meta.ConnectedAt = time.Now()
	}

	m.conns[id] = meta
	return nil
}

// Unregister removes a connection by id. It is safe to call multiple times.
func (m *ConnectionManager) Unregister(id string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.conns, id)
}

// UpdateMeta updates metadata for an existing connection (e.g. ClientID from first message).
func (m *ConnectionManager) UpdateMeta(id string, meta *ConnMeta) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if existing, ok := m.conns[id]; ok && meta != nil {
		if meta.ClientID != "" {
			existing.ClientID = meta.ClientID
		}

		if meta.RemoteAddr != "" {
			existing.RemoteAddr = meta.RemoteAddr
		}
	}
}

// ActiveCount returns the number of currently tracked connections.
func (m *ConnectionManager) ActiveCount() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.conns)
}

// NextID returns a unique numeric id for use in connection ids (e.g. peer + id).
func (m *ConnectionManager) NextID() uint64 {
	return m.nextID.Add(1)
}
