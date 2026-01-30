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

package discovery

import (
	"fmt"
	"strings"
	"time"

	"os-artificer/saber/pkg/gerrors"

	clientv3 "go.etcd.io/etcd/client/v3"
)

type EventType int

const (
	EventTypePut EventType = iota
	EventTypeDelete
	EventTypeRecover
)

const (
	defaultTTL                   = 6
	defaultDialTimeout           = 5 * time.Second
	defaultAutoSyncInterval      = 60 * time.Second
	defaultKeepAliveTime         = 30 * time.Second
	defaultKeepAliveTimeout      = 10 * time.Second
	defautlRegistryRootKeyPrefix = "/os-artificer/saber/registry"
)

type Event struct {
	Type  EventType
	Key   string
	Value []byte
}

type ClientOptions struct {
	Endpoints            []string
	User                 string
	Password             string
	DialTimeout          time.Duration
	AutoSyncInterval     time.Duration
	DialKeepAliveTime    time.Duration
	DialKeepAliveTimeout time.Duration
}

func NewClient(opt *ClientOptions) (*clientv3.Client, error) {

	if opt == nil {
		return nil, gerrors.New(gerrors.InvalidParameter, "")
	}

	if len(opt.Endpoints) == 0 {
		return nil, gerrors.New(gerrors.InvalidParameter, "endpoints are required")
	}

	if opt.User == "" {
		return nil, gerrors.New(gerrors.InvalidParameter, "user is required")
	}

	if opt.Password == "" {
		return nil, gerrors.New(gerrors.InvalidParameter, "password is required")
	}

	if opt.DialTimeout == 0 {
		opt.DialTimeout = defaultDialTimeout
	}

	if opt.AutoSyncInterval == 0 {
		opt.AutoSyncInterval = defaultAutoSyncInterval
	}

	if opt.DialKeepAliveTime == 0 {
		opt.DialKeepAliveTime = defaultKeepAliveTime
	}

	if opt.DialKeepAliveTimeout == 0 {
		opt.DialKeepAliveTimeout = defaultKeepAliveTimeout
	}

	cli, err := clientv3.New(clientv3.Config{
		Username:             opt.User,
		Password:             opt.Password,
		Endpoints:            opt.Endpoints,
		DialTimeout:          opt.DialTimeout,
		AutoSyncInterval:     opt.AutoSyncInterval,
		DialKeepAliveTime:    opt.DialKeepAliveTime,
		DialKeepAliveTimeout: opt.DialKeepAliveTimeout,
	})

	if err != nil {
		return nil, gerrors.New(gerrors.ComponentFailure, err.Error())
	}

	return cli, nil
}

func NewRegistry(c *clientv3.Client, serviceId string, ttl int64) (*Registry, error) {
	if c == nil {
		return nil, gerrors.New(gerrors.InvalidParameter, "client is required")
	}

	serviceId = strings.TrimSpace(serviceId)
	if serviceId == "" {
		return nil, gerrors.New(gerrors.InvalidParameter, "service id is required")
	}

	if ttl < defaultTTL {
		ttl = defaultTTL
	}

	registry := &Registry{
		serviceId: serviceId,
		eventChan: make(chan *Event, 1024),
		rootKey:   fmt.Sprintf("%s/%s", defautlRegistryRootKeyPrefix, serviceId),
		ttl:       ttl,
		client:    c,
	}

	return registry, nil
}

func NewDiscovery(c *clientv3.Client) (*Discovery, error) {

	if c == nil {
		return nil, gerrors.New(gerrors.InvalidParameter, "client is required")
	}

	discovery := &Discovery{
		exit:   make(chan struct{}),
		client: c,
	}

	return discovery, nil
}
