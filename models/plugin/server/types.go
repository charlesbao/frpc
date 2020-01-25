// Copyright 2019 fatedier, fatedier@gmail.com
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

package plugin

import (
	"github.com/charlesbao/frpc/models/msg"
)

type Request struct {
	Version string      `json:"version"`
	Op      string      `json:"op"`
	Content interface{} `json:"content"`
}

type Response struct {
	Reject       bool        `json:"reject"`
	RejectReason string      `json:"reject_reason"`
	Unchange     bool        `json:"unchange"`
	Content      interface{} `json:"content"`
}

type LoginContent struct {
	msg.Login
}

type UserInfo struct {
	User  string            `json:"user"`
	Metas map[string]string `json:"metas"`
}

type NewProxyContent struct {
	User UserInfo `json:"user"`
	msg.NewProxy
}
