// Copyright 2021 Comcast Cable Communications Management, LLC
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

package bit

import "fmt"

type Mask uint

func (b *Mask) Set(flag Mask) {
	*b = *b | flag
}

func (b *Mask) Clear(flag Mask) {
	*b = *b & ^flag
}

func (b *Mask) Flip(flag Mask) {
	*b ^= flag
}

func (b Mask) IsSet(flag Mask) bool {
	return b&flag != 0
}

func (b Mask) String() string {
	return fmt.Sprintf("%b", b)
}