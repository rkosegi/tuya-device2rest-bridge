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
	"fmt"
	"maps"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/rkosegi/go-http-commons/body"
	"github.com/rkosegi/tuya-device2rest-bridge/pkg/api"
	"github.com/rkosegi/tuya-proto/client"
	"github.com/rkosegi/tuya-proto/proto"
	"github.com/samber/lo"
)

func (i *impl) ListDevices(w http.ResponseWriter, _ *http.Request) {
	i.rwlock.RLock()
	defer i.rwlock.RUnlock()

	devs := lo.MapToSlice(i.cfg.Devices, func(name string, dev api.DeviceSpec) api.DeviceSpec {
		return dev
	})

	// sort by name, so output is deterministic
	slices.SortFunc(devs, func(d1 api.DeviceSpec, d2 api.DeviceSpec) int {
		return strings.Compare(d1.Name, d2.Name)
	})

	// obfuscate keys
	censored := lo.Map(devs, func(item api.DeviceSpec, _ int) *api.DeviceSpec {
		item.Key = dummyKey
		return &item
	})
	out.SendWithStatus(w, censored, http.StatusOK)
}

func (i *impl) GetDevice(w http.ResponseWriter, _ *http.Request, name api.Device) {
	i.rwlock.RLock()
	defer i.rwlock.RUnlock()

	if dev, exists := i.cfg.Devices[name]; exists {
		dev.Key = dummyKey
		out.SendWithStatus(w, dev, http.StatusOK)
		return
	}
	out.SendWithStatus(w, nil, http.StatusNotFound)
}

func (i *impl) deleteDeviceByName(name string) {
	if dev, exists := (i.mgr.devs)[name]; exists {
		_ = dev.Close()
	}
	delete(i.cfg.Devices, name)
	delete(i.mgr.devs, name)
}

func (i *impl) DeleteDevice(w http.ResponseWriter, _ *http.Request, name api.Device) {
	i.rwlock.Lock()
	defer i.rwlock.Unlock()

	i.deleteDeviceByName(name)
	out.SendWithStatus(w, nil, http.StatusNoContent)
}

func (i *impl) CreateDevice(w http.ResponseWriter, r *http.Request) {
	newDev, err := body.ConsumeAs[api.DeviceSpec](r)
	if err != nil {
		out.SendWithStatus(w, err, http.StatusBadRequest)
		return
	}
	// TODO validate device name
	i.rwlock.Lock()
	defer i.rwlock.Unlock()
	if _, exists := i.mgr.devs[newDev.Name]; exists {
		out.SendWithStatus(w, fmt.Errorf("device with name %s already exists", newDev.Name), http.StatusConflict)
		return
	}
	i.mgr.devs[newDev.Name] = i.mgr.newDevCtx(*newDev)

}

func (i *impl) SendCommand(w http.ResponseWriter, r *http.Request, device api.Device, cmdId int) {

	i.rwlock.RLock()
	defer i.rwlock.RUnlock()

	var err error
	if dev, exists := i.mgr.devs[device]; exists {

		if err = dev.Connect(); err != nil {
			out.SendWithStatus(w, err, http.StatusInternalServerError)
			return
		}
		cl := dev.Client()
		defer func(cl client.BlockingClient) {
			_ = cl.Close()
		}(cl)

		var rb *api.SendCommandJSONRequestBody
		if rb, err = body.ConsumeAs[api.SendCommandJSONRequestBody](r); err != nil {
			out.SendWithStatus(w, err, http.StatusBadRequest)
			return
		}

		payload := make(map[string]any)
		if rb.DefaultRequest != nil && *rb.DefaultRequest == true {
			payload["devId"] = dev.cfg.Id
			payload["gwId"] = dev.cfg.Id
			payload["uid"] = dev.cfg.Id
			payload["t"] = int(time.Now().Unix())
		}

		if rb.Payload != nil {
			maps.Copy(payload, *rb.Payload)
		}

		if err = cl.Send(proto.CmdIdType(cmdId), payload); err != nil {
			out.SendWithStatus(w, err, http.StatusInternalServerError)
			return
		}

		var m map[string]any
		if err = cl.Read(&m); err != nil {
			out.SendWithStatus(w, err, http.StatusInternalServerError)
			return
		}
		out.SendWithStatus(w, m, http.StatusOK)
		return
	}
}

func (i *impl) GetDeviceStats(w http.ResponseWriter, r *http.Request, device api.Device) {
	i.rwlock.RLock()
	defer i.rwlock.RUnlock()

	if dev, exists := i.mgr.devs[device]; exists {
		if dev.cl != nil {
			out.SendWithStatus(w, dev.cl.Stats(), http.StatusOK)
			return
		}
	}
	out.SendWithStatus(w, nil, http.StatusNotFound)
}
