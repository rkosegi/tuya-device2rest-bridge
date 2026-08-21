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
	"context"
	"log/slog"
	"net/http"
	"sync"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rkosegi/go-http-commons/middlewares"
	"github.com/rkosegi/go-http-commons/openapi"
	"github.com/rkosegi/go-http-commons/output"
	"github.com/rkosegi/tuya-device2rest-bridge/pkg/api"
)

var (
	out = output.NewBuilder().Build()
)

type impl struct {
	l      *slog.Logger
	mgr    *devMgr
	cfg    *api.AppConfig
	preg   prometheus.Gatherer
	rwlock sync.RWMutex
}

func (i *impl) Close() error { return i.mgr.Close() }

func (i *impl) Run(ctx context.Context) error {
	cors := handlers.CORS(
		handlers.AllowedMethods([]string{
			http.MethodGet,
			http.MethodDelete,
			http.MethodPatch,
			http.MethodPut,
			http.MethodPost,
		}),
		handlers.AllowedOrigins(i.cfg.Server.Cors.AllowedOrigins),
		handlers.MaxAge(i.cfg.Server.Cors.MaxAge),
		handlers.AllowedHeaders([]string{"Content-Type"}),
	)

	r := mux.NewRouter()
	r.HandleFunc("/spec/openapi.v1.json", openapi.SpecHandler(api.PathToRawSpec))
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("OK\n"))
	})
	if i.cfg.Server.Telemetry.Enabled {
		r.Handle("/metrics", promhttp.HandlerFor(prometheus.Gatherers{i.preg}, promhttp.HandlerOpts{
			ErrorHandling: promhttp.ContinueOnError,
		}))
	}

	mws := []api.MiddlewareFunc{
		middlewares.NewLoggingBuilder().WithLevel(slog.LevelDebug).WithLogger(i.l).Build(),
	}
	i.l.DebugContext(ctx, "starting server", "listen-address", i.cfg.Server.ListenAddress,
		"api-prefix", i.cfg.Server.APIPrefix)
	return i.cfg.Server.RunUntil(&http.Server{
		Handler: cors(api.HandlerWithOptions(i, api.GorillaServerOptions{
			BaseURL:     i.cfg.Server.APIPrefix,
			BaseRouter:  r,
			Middlewares: mws,
		})),
	}, ctx.Done())
}
