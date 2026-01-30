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

package discovery_test

import (
	"context"
	"log"
	"os"
	"strings"
	"testing"

	"os-artificer/saber/pkg/discovery"

	clientv3 "go.etcd.io/etcd/client/v3"
)

var (
	client *clientv3.Client
	reg    *discovery.Registry
	dis    *discovery.Discovery
)

func setup() {

	endpoints := os.Getenv("YLG_ETCD_ENDPOINTS")
	user := os.Getenv("YLG_ETCDTL_USER")
	password := os.Getenv("YLG_ETCDTL_PASSWORLD")

	log.Println("endpoints:", endpoints)
	log.Println("user:", user)
	log.Println("password:", password)

	cli, err := discovery.NewClient(&discovery.ClientOptions{
		Endpoints: strings.Split(endpoints, ","),
		User:      user,
		Password:  password,
	})

	if err != nil {
		log.Fatalf("failed to create etcd client. errmsg:%v", err)
	}
	client = cli

	registry, err := discovery.NewRegistry(client, "test-service-id", 10)
	if err != nil {
		log.Fatalf("failed to create registry. errmsg:%v", err)
	}

	reg = registry

	discover, err := discovery.NewDiscovery(client)
	if err != nil {
		log.Fatalf("failed to create discovery. errmsg:%v", err)
	}

	dis = discover

	reg.SetService(context.Background(), "test-id")
}

func tear() {
	reg.Close()
	dis.Close()
	client.Close()
}

func TestSetService(t *testing.T) {

	err := reg.SetService(context.Background(), "test-id")
	if err != nil {
		t.Fatalf("failed to set registry service. errmsg:%v", err)
	}
}

func TestGetWithPrefix(t *testing.T) {

	resp, err := dis.GetWithPrefix(context.Background(), "/")
	if err != nil {
		t.Fatalf("failed to get with prefix. errmsg:%v", err)
	}

	for key, value := range resp {
		t.Logf("key:%s, value:%s", key, string(value))
	}
}

func TestMain(m *testing.M) {

	setup()
	code := m.Run()
	tear()
	os.Exit(code)
}
