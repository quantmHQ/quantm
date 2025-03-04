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

package erratic

type (
	// Hints represents a map of key-value pairs providing additional information about an error.
	//
	// Example:
	//
	//  info := Hints{"field": "invalid value"}
	//  fmt.Println(info) // Output: map[string]string{"field": "invalid value"}
	Hints map[string]string
)

func NewHints(args ...string) Hints {
	odd := false

	if len(args)%2 != 0 {
		odd = true
	}

	details := make(Hints)

	for i := 0; i < len(args); i += 2 {
		details[args[i]] = args[i+1]
	}

	if odd {
		details["unknown"] = args[len(args)-1]
	}

	return details
}
