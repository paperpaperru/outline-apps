// Copyright 2019 The Outline Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package xrayMobile

import (
	"errors"
	"github.com/Jigsaw-Code/outline-apps/client/src/tun2socks/tunnel"
	"github.com/eycorsican/go-tun2socks/core"
	"github.com/eycorsican/go-tun2socks/proxy/socks"
	"io"
	"runtime/debug"
	"time"
)

type Tunnel interface {
	tunnel.Tunnel
}

// TunWriter is an interface that allows for outputting packets to the TUN (VPN).
type TunWriter interface {
	io.WriteCloser
}

func init() {
	// Apple VPN extensions have a memory limit of 15MB. Conserve memory by increasing garbage
	// collection frequency and returning memory to the OS every minute.
	debug.SetGCPercent(10)
	// TODO: Check if this is still needed in go 1.13, which returns memory to the OS
	// automatically.
	ticker := time.NewTicker(time.Minute * 1)
	go func() {
		for range ticker.C {
			debug.FreeOSMemory()
		}
	}()
}

func newTunnel(tunWriter io.WriteCloser) (tunnel.Tunnel, error) {
	if tunWriter == nil {
		return nil, errors.New("Must provide a TUN writer")
	}
	core.RegisterOutputFn(func(data []byte) (int, error) {
		return tunWriter.Write(data)
	})
	lwipStack := core.NewLWIPStack()
	core.RegisterTCPConnHandler(socks.NewTCPHandler("127.0.0.1", 1080))
	core.RegisterUDPConnHandler(socks.NewUDPHandler("127.0.0.1", 1080, 30*time.Second))

	// Copy packets from tun device to lwip stack, it's the main loop.

	return tunnel.NewTunnel(tunWriter, lwipStack), nil
}

// ConnectLocalSocksTunnel reads packets from a TUN device and routes it to a local socks proxy server.
// Returns a Tunnel instance that should be used to input packets to the tunnel.

func ConnectLocalSocksTunnel(tunWriter TunWriter) (Tunnel, error) {
	if tunWriter == nil {
		return nil, errors.New("must provide a TunWriter")
	}
	return newTunnel(tunWriter)
}
