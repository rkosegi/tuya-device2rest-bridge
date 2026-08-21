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

package main

import (
	"context"
	"errors"
	"flag"
	"net/http"

	"github.com/rkosegi/go-http-commons/servertypes"
	xlog "github.com/rkosegi/slog-config"
	"github.com/rkosegi/tuya-device2rest-bridge/pkg/api"
	"github.com/rkosegi/tuya-device2rest-bridge/pkg/server"
	"github.com/rkosegi/yaml-toolkit/fluent"
)

func main() {
	var cfgFile string
	sc := xlog.MustNew("info", xlog.LogFormatLogFmt)
	flag.StringVar(&cfgFile, "config", "config.yaml", "config file")
	sc.AddFlags(flag.CommandLine)
	flag.Parse()
	logger := sc.Logger()
	logger.Debug("loading config", "file", cfgFile)
	cfg := fluent.NewConfigHelper[api.AppConfig]().
		Add(server.DefaultConfig()).
		Load(cfgFile).Result()

	for dn := range cfg.Devices {
		dev := cfg.Devices[dn]
		dev.Name = dn
		dev.Normalize()
		cfg.Devices[dn] = dev
	}

	var (
		srv servertypes.RunCloser
		err error
	)
	if srv, err = server.New(cfg, server.WithLogger(logger)); err != nil {
		panic(err)
	}

	defer func(srv servertypes.RunCloser) {
		_ = srv.Close()
	}(srv)

	if err = srv.Run(context.Background()); err != nil && !errors.Is(err, http.ErrServerClosed) {
		panic(err)
	}
}
