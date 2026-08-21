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
	"sync"

	"github.com/rkosegi/tuya-device2rest-bridge/pkg/api"
	"github.com/rkosegi/tuya-proto/client"
	"github.com/rkosegi/tuya-proto/proto"
	"github.com/samber/lo"
)

type devCtx struct {
	cfg      api.DeviceSpec
	protoVer proto.Version
	l        *slog.Logger
	cl       client.BlockingClient
}

func (d *devCtx) Client() client.BlockingClient {
	return d.cl
}

func (d *devCtx) Connect() error {
	d.cl = client.NewBlockingClient(d.protoVer, d.cfg.Address, []byte(d.cfg.Key),
		client.WithTimeout(d.cfg.ConnectTimeout),
		client.WithReadTimeout(d.cfg.ReadTimeout),
		client.WithWriteTimeout(d.cfg.WriteTimeout),
		client.WithLogger(d.l))
	d.l.Debug("Connect")
	return d.cl.Connect()
}

func (d *devCtx) Close() (err error) {
	if d.cl != nil {
		d.l.Debug("closing client")
		err = d.cl.Close()
		d.cl = nil
	}
	return
}

func (d *devMgr) newDevCtx(dev api.DeviceSpec) *devCtx {
	return &devCtx{
		cfg:      dev,
		protoVer: dev.Protocol.AsProto(),
		l:        d.l.With("device", dev.Name),
	}
}

type devMgr struct {
	l    *slog.Logger
	lock sync.RWMutex
	devs map[string]*devCtx
}

func (d *devMgr) Close() error {
	lo.ForEach(lo.Values(d.devs), func(dc *devCtx, _ int) {
		_ = dc.Close()
	})
	clear(d.devs)
	return nil
}

func newDevMgr(cfg *api.AppConfig, l *slog.Logger) *devMgr {
	dm := &devMgr{
		devs: make(map[string]*devCtx),
		l:    l,
	}
	for _, dev := range cfg.Devices {
		dm.devs[dev.Name] = dm.newDevCtx(dev)
	}
	return dm
}
