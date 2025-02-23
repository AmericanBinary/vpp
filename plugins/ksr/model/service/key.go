// Copyright (c) 2017 Cisco and/or its affiliates.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at:
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package service

import (
	"github.com/americanbinary/vpp/plugins/ksr/model/ksrkey"
)

const (
	// ServiceKeyword defines the keyword identifying Service data.
	ServiceKeyword = "service"
)

// KeyPrefix returns the key prefix identifying all K8s services in the
// data store.
func KeyPrefix() string {
	return ksrkey.KeyPrefix(ServiceKeyword)
}

// ParseServiceFromKey parses service and namespace ids from the associated
// data-store key.
func ParseServiceFromKey(key string) (service string, namespace string, err error) {
	return ksrkey.ParseNameFromKey(ServiceKeyword, key)
}

// Key returns the key under which the given service is stored in the data-store.
func Key(name string, namespace string) string {
	return ksrkey.Key(ServiceKeyword, name, namespace)
}
