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
	"sync"
	"sync/atomic"
)

type ResetOnce struct {
	done uint32
	m    sync.Mutex
}

func (o *ResetOnce) Do(f func() error) error {
	if atomic.LoadUint32(&o.done) != 0 {
		return nil
	}

	return o.do(f)
}

func (o *ResetOnce) do(f func() error) error {
	o.m.Lock()
	defer o.m.Unlock()

	if o.done != 0 {
		return nil
	}

	if err := f(); err != nil {
		return err
	}

	atomic.StoreUint32(&o.done, 1)
	return nil
}

func (o *ResetOnce) Reset() {
	atomic.StoreUint32(&o.done, 0)
}
