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

package defs

import (
	"encoding/json"
	"strings"
	"time"
)

type (
	Timestamp time.Time // Timestamp is hack around github's funky use of time
)

func (t *Timestamp) UnmarshalJSON(data []byte) error {
	var raw any
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	switch v := raw.(type) {
	case float64:
		*t = Timestamp(time.Unix(int64(v), 0))
	case string:
		if strings.HasSuffix(v, "Z") {
			t_, err := time.Parse("2006-01-02T15:04:05Z", v)
			if err != nil {
				return err
			}

			*t = Timestamp(t_)
		} else {
			t_, err := time.Parse(time.RFC3339, v)
			if err != nil {
				return err
			}

			*t = Timestamp(t_)
		}
	}

	return nil
}

func (t Timestamp) MarshalJSON() ([]byte, error) {
	t_ := time.Time(t)
	return json.Marshal(t_.Format(time.RFC3339))
}

func (t Timestamp) Time() time.Time {
	return time.Time(t)
}
