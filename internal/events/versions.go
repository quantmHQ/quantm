// Crafted with ❤ at Breu, Inc. <info@breu.io>, Copyright © 2024.
//
// Functional Source License, Version 1.1, Apache 2.0 Future License
//
// We hereby irrevocably grant you an additional license to use the Software under the Apache License, Version 2.0 that
// is effective on the second anniversary of the date we make the Software available. On or after that date, you may use
// the Software under the Apache License, Version 2.0, in which case the following will apply:
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with
// the License.
//
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on
// an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the
// specific language governing permissions and limitations under the License.

package events

type (
	EventVersion string // Version is the version of the event.
)

// String returns the string representation of the Version.
func (ev EventVersion) String() string {
	return string(ev)
}

const (
	Version_0_1_0 EventVersion = "0.1.0" // version 0.1.0.
	Version_0_1_1 EventVersion = "0.1.1" // version 0.1.1.
)

const (
	// EventVersionDefault alias for the default version. This allows for easy versioning without chaniging the code base.
	EventVersionDefault = Version_0_1_0
)
