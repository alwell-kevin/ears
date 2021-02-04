// Copyright 2020 Comcast Cable Communications Management, LLC
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

package debug

import (
	"github.com/xmidt-org/ears/pkg/receiver"
	rdebug "github.com/xmidt-org/ears/pkg/receiver/debug"
	"github.com/xmidt-org/ears/pkg/sender"
	sdebug "github.com/xmidt-org/ears/pkg/sender/debug"

	"github.com/xmidt-org/ears/pkg/plugin"
)

var (
	Name    = "debug"
	Version = "v0.0.0"
	Commit  = ""
)

func NewPlugin() (*plugin.Plugin, error) {
	return NewPluginVersion(Name, Version, Commit)
}

func NewPluginVersion(name string, version string, commitID string) (*plugin.Plugin, error) {
	return plugin.NewPlugin(
		plugin.WithName(name),
		plugin.WithVersion(version),
		plugin.WithCommitID(commitID),
		plugin.WithNewReceiver(NewReceiver),
		plugin.WithNewSender(NewSender),
	)
}

func NewReceiver(config interface{}) (receiver.Receiver, error) {
	return rdebug.NewReceiver(config)
}

func NewSender(config interface{}) (sender.Sender, error) {
	return sdebug.NewSender(config)
}