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
	"strings"
	"sync"

	"os-artificer/saber/pkg/gerrors"
	"os-artificer/saber/pkg/logger"

	clientv3 "go.etcd.io/etcd/client/v3"
)

type Discovery struct {
	exit   chan struct{}
	wg     sync.WaitGroup
	client *clientv3.Client
}

func (d *Discovery) Watch(ctx context.Context, key string) (chan *Event, error) {

	key = strings.TrimSpace(key)
	if key == "" {
		return nil, gerrors.New(gerrors.InvalidParameter, "the watched key is required")
	}

	watchChan := d.client.Watch(ctx, key)
	eventChan := make(chan *Event, 1024)

	d.wg.Add(1)

	go func(ctx context.Context) {

		defer d.wg.Done()

		for {
			select {
			case <-d.exit:
				logger.Infof("exit the watcher. key:%s", key)
				return

			case <-ctx.Done():
				logger.Infof("exit the watcher. key:%s", key)
				return

			case watchResp := <-watchChan:
				for _, event := range watchResp.Events {
					switch event.Type {
					case clientv3.EventTypePut:
						event := &Event{
							Type:  EventTypePut,
							Key:   string(event.Kv.Key),
							Value: []byte(event.Kv.Value),
						}
						eventChan <- event

					case clientv3.EventTypeDelete:
						event := &Event{
							Type:  EventTypeDelete,
							Key:   string(event.Kv.Key),
							Value: []byte(event.Kv.Value),
						}
						eventChan <- event
					}
				}
			}
		}
	}(ctx)

	return eventChan, nil
}

func (d *Discovery) WatchWithPrefix(ctx context.Context, keyPrefix string) (chan *Event, error) {

	keyPrefix = strings.TrimSpace(keyPrefix)
	if keyPrefix == "" {
		return nil, gerrors.New(gerrors.InvalidParameter, "the watched key is required")
	}

	watchChan := d.client.Watch(ctx, keyPrefix, clientv3.WithPrefix())
	eventChan := make(chan *Event, 1024)

	d.wg.Add(1)
	go func(ctx context.Context) {
		defer d.wg.Done()

		for {
			select {
			case <-d.exit:
				logger.Infof("exit the watcher. key prefix:%s", keyPrefix)
				return

			case <-ctx.Done():
				logger.Infof("exit the watcher. key prefix:%s", keyPrefix)
				return

			case watchResp := <-watchChan:
				for _, event := range watchResp.Events {
					switch event.Type {
					case clientv3.EventTypePut:
						event := &Event{
							Type:  EventTypePut,
							Key:   string(event.Kv.Key),
							Value: []byte(event.Kv.Value),
						}
						eventChan <- event

					case clientv3.EventTypeDelete:
						event := &Event{
							Type:  EventTypeDelete,
							Key:   string(event.Kv.Key),
							Value: []byte(event.Kv.Value),
						}
						eventChan <- event
					}
				}
			}
		}
	}(ctx)

	return eventChan, nil
}

func (d *Discovery) Get(ctx context.Context, key string) ([]byte, error) {

	key = strings.TrimSpace(key)
	if key == "" {
		return nil, gerrors.New(gerrors.InvalidParameter, "the key prefix is required")
	}

	resp, err := d.client.Get(ctx, key)
	if err != nil {
		return nil, gerrors.New(gerrors.ComponentFailure, err.Error())
	}

	var value []byte
	for _, kv := range resp.Kvs {
		value = []byte(kv.Value)
	}

	return value, nil
}

func (d *Discovery) GetWithPrefix(ctx context.Context, keyPrefix string) (map[string][]byte, error) {

	keyPrefix = strings.TrimSpace(keyPrefix)
	if keyPrefix == "" {
		return nil, gerrors.New(gerrors.InvalidParameter, "the key prefix is required")
	}

	resp, err := d.client.Get(ctx, keyPrefix, clientv3.WithPrefix())
	if err != nil {
		return nil, gerrors.New(gerrors.ComponentFailure, err.Error())
	}

	kvs := make(map[string][]byte, len(resp.Kvs))
	for _, kv := range resp.Kvs {
		kvs[string(kv.Key)] = []byte(kv.Value)
	}

	return kvs, nil
}

func (d *Discovery) Close() {

	close(d.exit)
	d.exit = nil
	d.wg.Wait()
}
