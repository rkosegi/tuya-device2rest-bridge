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

package server

import (
	"time"

	"github.com/rkosegi/go-http-commons/config"
	"github.com/rkosegi/tuya-device2rest-bridge/pkg/api"
)

const (
	dummyKey = "*******"
)

// DefaultConfig gets the default application configuration
func DefaultConfig() *api.AppConfig {
	return &api.AppConfig{
		Server: &config.ServerConfig{
			ListenAddress: ":22005",
			APIPrefix:     "/api/v1",
			Cors: &config.CorsConfig{
				AllowedOrigins: []string{"*"},
				MaxAge:         86400,
			},
			Telemetry: &config.TelemetryConfig{
				Enabled: true,
				Path:    "/metrics",
			},
			ReadTimeout:       new(30 * time.Second),
			ReadHeaderTimeout: new(30 * time.Second),
			WriteTimeout:      new(30 * time.Second),
			IdleTimeout:       new(30 * time.Second),
		},
	}
}
