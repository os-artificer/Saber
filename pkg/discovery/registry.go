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

package discovery

import (
	"context"
	"fmt"
	"math"
	"strings"
	"sync"
	"time"

	"os-artificer/saber/pkg/gerrors"
	"os-artificer/saber/pkg/logger"

	clientv3 "go.etcd.io/etcd/client/v3"
)

type Registry struct {
	serviceId     string
	rootKey       string
	ttl           int64
	wg            sync.WaitGroup
	mu            sync.Mutex
	client        *clientv3.Client
	leaseId       clientv3.LeaseID
	exit          chan struct{}
	eventChan     chan *Event
	keepAliveChan <-chan *clientv3.LeaseKeepAliveResponse
}

func (r *Registry) grant(ctx context.Context) error {

	leaseResp, err := r.client.Grant(ctx, r.ttl)
	if err != nil {
		return gerrors.New(gerrors.Failure, err.Error())
	}
	r.leaseId = leaseResp.ID

	keepAliveChan, err := r.client.KeepAlive(ctx, r.leaseId)
	if err != nil {
		return gerrors.New(gerrors.Failure, err.Error())
	}
	r.keepAliveChan = keepAliveChan

	if r.exit == nil {
		r.exit = make(chan struct{})
	}

	r.wg.Add(2)
	go r.monitorKeepalive(ctx)
	go r.checkLeaseTTL(ctx)

	return nil
}

func (r *Registry) monitorKeepalive(ctx context.Context) {

	defer r.wg.Done()

	for {
		select {
		case <-r.exit:
			logger.Info("exit registry monitor keepalive")
			return

		case <-ctx.Done():
			logger.Info("exit registry monitor keepalive")
			return

		case _, ok := <-r.keepAliveChan:
			if !ok {
				r.wg.Add(1)
				go func(ctx context.Context) {
					defer r.wg.Done()
					r.recoverLease(ctx)
				}(ctx)
			}
		}
	}

}

func (r *Registry) checkLeaseTTL(ctx context.Context) {

	defer r.wg.Done()

	ttl := time.Duration(math.Floor(float64(r.ttl) / 2))
	ticker := time.NewTicker(ttl * time.Second)

	defer ticker.Stop()

	for {
		select {
		case <-r.exit:
			logger.Info("exit registry check lease ttl")
			return

		case <-ctx.Done():
			logger.Info("exit registry check lease ttl")
			return

		case <-ticker.C:
			ttlResp, err := r.client.TimeToLive(ctx, r.leaseId)
			if err != nil || ttlResp.TTL <= 0 {
				r.wg.Add(1)
				go func(ctx context.Context) {
					defer r.wg.Done()
					r.recoverLease(ctx)
				}(ctx)
			}
		}
	}
}

func (r *Registry) recoverLease(ctx context.Context) error {

	r.mu.Lock()
	defer r.mu.Unlock()

	ttlResp, err := r.client.Lease.TimeToLive(ctx, r.leaseId)
	if err == nil && ttlResp.TTL > 0 {
		r.wg.Add(2)
		return r.grant(ctx)
	}

	return nil
}

func (r *Registry) SetService(ctx context.Context, value string) error {

	value = strings.TrimSpace(value)

	if r.serviceId == "" {
		return gerrors.New(gerrors.InvalidParameter, "serviceId is required")
	}

	if r.leaseId == 0 {
		if err := r.grant(ctx); err != nil {
			return err
		}
	}

	_, err := r.client.Put(ctx, r.rootKey, value, clientv3.WithLease(r.leaseId))
	if err != nil {
		return gerrors.New(gerrors.Failure, err.Error())
	}

	return nil
}

func (r *Registry) Events() chan *Event {
	return r.eventChan
}

func (r *Registry) Set(ctx context.Context, key, value string) error {

	key = strings.TrimSpace(key)
	value = strings.TrimSpace(value)

	if key == "" {
		return gerrors.New(gerrors.InvalidParameter, "key is required")
	}

	if !strings.HasPrefix(key, r.rootKey) {
		key = fmt.Sprintf("%s/%s", r.rootKey, key)
	}

	_, err := r.client.Put(ctx, key, value, clientv3.WithLease(r.leaseId))
	if err != nil {
		return gerrors.New(gerrors.Failure, err.Error())
	}

	return nil
}

func (r *Registry) Close() {
	close(r.exit)
	r.exit = nil
	r.wg.Wait()
}
