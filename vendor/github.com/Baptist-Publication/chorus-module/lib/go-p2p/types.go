// Copyright 2017 Baptist-Publication Information Technology Services Co.,Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package p2p

import (
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/Baptist-Publication/chorus-module/lib/go-crypto"
)

const maxNodeInfoSize = 10240 // 10Kb

type NodeInfo struct {
	PubKey      crypto.PubKeyEd25519 `json:"pub_key"`
	SigndPubKey string               `json:"signd_pub_key"`
	Moniker     string               `json:"moniker"`
	Network     string               `json:"network"`
	RemoteAddr  string               `json:"remote_addr"`
	ListenAddr  string               `json:"listen_addr"`
	Version     string               `json:"version"` // major.minor.revision
	Other       []string             `json:"other"`   // other application specific data
}

// CONTRACT: two nodes are compatible if the major/minor versions match and network match
func (info *NodeInfo) CompatibleWith(other *NodeInfo) error {
	iMajor, iMinor, _, iErr := splitVersion(info.Version)
	oMajor, oMinor, _, oErr := splitVersion(other.Version)

	// if our own version number is not formatted right, we messed up
	if iErr != nil {
		return iErr
	}

	// version number must be formatted correctly ("x.x.x")
	if oErr != nil {
		return oErr
	}

	// major version must match
	if iMajor != oMajor {
		return fmt.Errorf("Peer is on a different major version. Got %v, expected %v", oMajor, iMajor)
	}

	// minor version must match
	if iMinor != oMinor {
		return fmt.Errorf("Peer is on a different minor version. Got %v, expected %v", oMinor, iMinor)
	}

	// nodes must be on the same network
	if info.Network != other.Network {
		return fmt.Errorf("Peer is on a different network. Got %v, expected %v", other.Network, info.Network)
	}

	return nil
}

func (info *NodeInfo) ListenHost() string {
	host, _, _ := net.SplitHostPort(info.ListenAddr)
	return host
}

func (info *NodeInfo) ListenPort() int {
	_, port, _ := net.SplitHostPort(info.ListenAddr)
	port_i, err := strconv.Atoi(port)
	if err != nil {
		return -1
	}
	return port_i
}

func splitVersion(version string) (string, string, string, error) {
	spl := strings.Split(version, ".")
	if len(spl) != 3 {
		return "", "", "", fmt.Errorf("Invalid version format %v", version)
	}
	return spl[0], spl[1], spl[2], nil
}
