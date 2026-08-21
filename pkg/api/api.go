/*
Copyright 2026 Richard Kosegi

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package api

import (
	"time"

	"github.com/rkosegi/tuya-proto/proto"
)

//go:generate go tool oapi-codegen --config=openapi-config.yaml openapi.yaml

func (dp DeviceProtocol) AsProto() proto.Version {
	switch dp {
	case Tuya31:
		return proto.Version31
	case Tuya34:
		return proto.Version34
	default:
		return proto.Version31
	}
}

func (dev *DeviceSpec) Normalize() {
	if dev.ConnectTimeout == 0 {
		dev.ConnectTimeout = 10 * time.Second
	}
	if dev.ReadTimeout == 0 {
		dev.ReadTimeout = 10 * time.Second
	}
	if dev.WriteTimeout == 0 {
		dev.WriteTimeout = 10 * time.Second
	}
}
