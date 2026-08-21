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
	"log/slog"

	"github.com/rkosegi/go-http-commons/servertypes"
	"github.com/rkosegi/tuya-device2rest-bridge/pkg/api"
)

type Opt func(srv servertypes.RunCloser)

func WithLogger(l *slog.Logger) Opt {
	return func(i servertypes.RunCloser) {
		i.(*impl).l = l
	}
}

func New(cfg *api.AppConfig, opts ...Opt) (servertypes.RunCloser, error) {
	i := &impl{cfg: cfg}
	for _, opt := range append([]Opt{
		WithLogger(slog.Default()),
	}, opts...) {
		opt(i)
	}
	i.mgr = newDevMgr(cfg, i.l)
	return i, nil
}
