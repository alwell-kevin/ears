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

package app

import (
	"fmt"
	"github.com/rs/zerolog"
	"github.com/xmidt-org/ears/internal/pkg/config"
	"github.com/xmidt-org/ears/internal/pkg/rtsemconv"
	"os"
)

func ProvideLogger(config config.Config) (*zerolog.Logger, error) {
	logLevel, err := zerolog.ParseLevel(config.GetString("ears.logLevel"))
	if err != nil {
		return nil, &InvalidOptionError{
			Option: fmt.Sprintf("loglevel %s is not valid", config.GetString("ears.logLevel")),
		}
	}
	hostname := config.GetString("ears.hostname")
	if hostname == "" {
		hostname, _ = os.Hostname()
	}
	logger := zerolog.New(os.Stdout).Level(logLevel).With().
		Str(rtsemconv.EarsLogHostnameKey, hostname).
		Timestamp().Logger()
	zerolog.LevelFieldName = "log.level"
	return &logger, nil
}
