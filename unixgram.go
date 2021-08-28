// Copyright (c) 2017 Dave Pifke.
//
// Redistribution and use in source and binary forms, with or without
// modification, is permitted provided that the following conditions are met:
//
// 1. Redistributions of source code must retain the above copyright notice,
//    this list of conditions and the following disclaimer.
//
// 2. Redistributions in binary form must reproduce the above copyright notice,
//    this list of conditions and the following disclaimer in the documentation
//    and/or other materials provided with the distribution.
//
// 3. Neither the name of the copyright holder nor the names of its
//    contributors may be used to endorse or promote products derived from
//    this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
// AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
// IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
// ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE
// LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
// CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
// SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
// INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN
// CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
// ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
// POSSIBILITY OF SUCH DAMAGE.

package wpasupplicant

import (
	"bytes"
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"net"
	"os"
	"path"
	"strconv"
	"strings"
	"sync"
	"time"
)

// message is a queued response (or read error) from the wpa_supplicant
// daemon.  Messages may be either solicited or unsolicited.
type message struct {
	priority int
	data     []byte
	err      error
}

// unixgram is the implementation of Conn for the AF_UNIX SOCK_DGRAM
// control interface.
//
// See https://w1.fi/wpa_supplicant/devel/ctrl_iface_page.html.
type unixgram struct {
	ctx                    context.Context
	local, remote          string
	conn                   *net.UnixConn
	solicited, unsolicited chan message
	wpaEvents              chan WPAEvent
	lock                   sync.Mutex
}

var _ Conn = (*unixgram)(nil)

// Connect returns a connection to wpa_supplicant for the specified
// interface, using the socket-based control interface.
func Connect(ctx context.Context, iface string, options ...Option) (Conn, error) {
	return ConnectPath(ctx, stdSocketPath, iface, options...)
}

// ConnectPath connects to iface within ctrlPath and returns a connection.
func ConnectPath(ctx context.Context, ctrlPath string, iface string, options ...Option) (Conn, error) {
	var err error
	uc := &unixgram{
		ctx:         ctx,
		solicited:   make(chan message),
		unsolicited: make(chan message),
		wpaEvents:   make(chan WPAEvent),
	}

	local, err := createLocalPath(iface)
	if err != nil {
		return nil, err
	}

	defaults := []Option{CustomUnixgram(local, path.Join(ctrlPath, iface))}
	defaults = append(defaults, options...)

	for _, fn := range defaults {
		if err = fn(uc); err != nil {
			return nil, err
		}
	}

	go uc.readLoop()
	go uc.readUnsolicited()

	// Issue an ATTACH command to start receiving unsolicited events.
	err = uc.runCommand("ATTACH")
	if err != nil {
		return nil, err
	}

	return uc, nil
}

// readLoop is spawned after we connect.  It receives messages from the
// socket, and routes them to the appropriate channel based on whether they
// are solicited (in response to a request) or unsolicited.
func (uc *unixgram) readLoop() {
	buf := make([]byte, 2048)
	for {
		select {
		case <-uc.ctx.Done():
			return
		default:
			b, oob := make([]byte, 2048), make([]byte, 2048)
			n, oobn, _, _, err := uc.conn.ReadMsgUnix(b, oob)
			if err != nil {
				uc.solicited <- message{
					err: err,
				}
				continue
			}

			// rebuild data buffer
			b = append(buf, b[:n]...)
			buf = oob[:oobn]

			// Unsolicited messages are preceded by a priority
			// specification, e.g. "<1>message".  If there's no priority,
			// default to 2 (info) and assume it's the response to
			// whatever command was last issued.
			var p int
			var c chan message
			if len(b) >= 3 && b[0] == '<' && b[2] == '>' {
				switch b[1] {
				case '0', '1', '2', '3', '4':
					c = uc.unsolicited
					p, _ = strconv.Atoi(string(b[1]))
					b = b[3:]
				default:
					c = uc.solicited
					p = 2
				}
			} else {
				c = uc.solicited
				p = 2
			}

			c <- message{
				priority: p,
				data:     b,
			}
		}
	}
}

// readUnsolicited handles messages sent to the unsolicited channel and parse them
// into a WPAEvent. At the moment we only handle `CTRL-EVENT-*` events and only events
// where the 'payload' is formatted with key=val.
func (uc *unixgram) readUnsolicited() {
	for {
		select {
		case <-uc.ctx.Done():
			return
		default:
			mgs := <-uc.unsolicited
			data := bytes.NewBuffer(mgs.data).String()

			parts := strings.Split(data, " ")
			if len(parts) == 0 {
				continue
			}

			var e WPAEvent
			if strings.Index(parts[0], "CTRL-") != 0 {
				e = WPAEvent{
					Event: "MESSAGE",
					Line:  data,
				}
			} else {
				e = WPAEvent{
					Event:     strings.TrimPrefix(parts[0], "CTRL-EVENT-"),
					Arguments: make(map[string]string),
					Line:      data,
				}

				for _, args := range parts[1:] {
					if strings.Contains(args, "=") {
						keyval := strings.Split(args, "=")
						if len(keyval) != 2 {
							continue
						}
						e.Arguments[keyval[0]] = keyval[1]
					}
				}
			}

			select {
			case uc.wpaEvents <- e:
			case <-time.After(time.Millisecond):
			}
		}
	}
}

// cmd executes a command and waits for a reply.
func (uc *unixgram) cmd(cmd string) ([]byte, error) {
	uc.lock.Lock()
	defer uc.lock.Unlock()

	_, err := uc.conn.Write([]byte(cmd))
	if err != nil {
		return nil, err
	}

	msg := <-uc.solicited
	return msg.data, msg.err
}

// runCommand is a wrapper around the uc.cmd command which makes sure the
// command returned a successful (OK) response.
func (uc *unixgram) runCommand(cmd string) error {
	resp, err := uc.cmd(cmd)
	if err != nil {
		return err
	}

	if bytes.Equal(resp, []byte("OK\n")) {
		return nil
	}

	return &ParseError{Line: string(resp)}
}

func (uc *unixgram) EventQueue() chan WPAEvent {
	return uc.wpaEvents
}

func (uc *unixgram) Ping() error {
	resp, err := uc.cmd("PING")
	if err != nil {
		return err
	}

	if bytes.Equal(resp, []byte("PONG\n")) {
		return nil
	}
	return &ParseError{Line: string(resp)}
}

func (uc *unixgram) AddNetwork() (int, error) {
	resp, err := uc.cmd("ADD_NETWORK")
	if err != nil {
		return -1, err
	}

	b := bytes.NewBuffer(resp)
	return strconv.Atoi(strings.Trim(b.String(), "\n"))
}

func (uc *unixgram) EnableNetwork(networkID int) error {
	return uc.runCommand(fmt.Sprintf("ENABLE_NETWORK %d", networkID))
}

func (uc *unixgram) EnableAllNetworks() error {
	return uc.runCommand("ENABLE_NETWORK all")
}

func (uc *unixgram) SelectNetwork(networkID int) error {
	return uc.runCommand(fmt.Sprintf("SELECT_NETWORK %d", networkID))
}

func (uc *unixgram) DisableNetwork(networkID int) error {
	return uc.runCommand(fmt.Sprintf("DISABLE_NETWORK %d", networkID))
}

func (uc *unixgram) RemoveNetwork(networkID int) error {
	return uc.runCommand(fmt.Sprintf("REMOVE_NETWORK %d", networkID))
}

func (uc *unixgram) RemoveAllNetworks() error {
	return uc.runCommand("REMOVE_NETWORK all")
}

func (uc *unixgram) SetNetwork(networkID int, variable string, value interface{}) error {
	b := strings.Builder{}
	b.WriteString("SET_NETWORK")
	b.WriteString(" ")
	b.WriteString(strconv.Itoa(networkID))
	b.WriteString(" ")
	b.WriteString(variable)
	b.WriteString(" ")

	// Since key_mgmt expects the value to not be wrapped in "" we do a little check here.
	// Update: since we have to support AP mode, we need to support integer value (and hex value that just for non-ascii ssid)
	switch v := value.(type) {
	case string:
		switch variable {
		case "key_mgmt":
			b.WriteString(v)
		default:
			b.WriteString("\"")
			b.WriteString(v)
			b.WriteString("\"")
		}
	case int:
		b.WriteString(strconv.Itoa(v))
	case []byte:
		b.WriteString(hex.EncodeToString(v))
	default:
		return errors.New("unsupported value type")
	}

	return uc.runCommand(b.String())
}

func (uc *unixgram) SaveConfig() error {
	return uc.runCommand("SAVE_CONFIG")
}

func (uc *unixgram) Reconfigure() error {
	return uc.runCommand("RECONFIGURE")
}

func (uc *unixgram) Reassociate() error {
	return uc.runCommand("REASSOCIATE")
}

func (uc *unixgram) Reconnect() error {
	return uc.runCommand("RECONNECT")
}

func (uc *unixgram) Scan() error {
	return uc.runCommand("SCAN")
}

func (uc *unixgram) ScanResults() ([]ScanResult, []error) {
	resp, err := uc.cmd("SCAN_RESULTS")
	if err != nil {
		return nil, []error{err}
	}

	return parseScanResults(bytes.NewBuffer(resp))
}

func (uc *unixgram) Status() (StatusResult, error) {
	resp, err := uc.cmd("STATUS")
	if err != nil {
		return nil, err
	}

	return parseStatusResults(bytes.NewBuffer(resp))
}

func (uc *unixgram) ListNetworks() ([]ConfiguredNetwork, error) {
	resp, err := uc.cmd("LIST_NETWORKS")
	if err != nil {
		return nil, err
	}

	return parseListNetworksResult(bytes.NewBuffer(resp))
}

func (uc *unixgram) Close() error {
	defer os.Remove(uc.local)

	if err := uc.runCommand("DETACH"); err != nil {
		return err
	}

	return uc.conn.Close()
}
