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

package utils

import (
	"github.com/google/uuid"
)

// NewUUID generates a new version 7 UUID.  It returns an error if UUID generation fails.
func NewUUID() (uuid.UUID, error) {
	return uuid.NewV7()
}

// MustUUID generates a new version 7 UUID. It panics if UUID generation fails.
//
// The only condition under which it could theoretically return an error is if the underlying system's source of
// randomness is completely broken or unavailable.  This is an exceptionally rare and serious system-level problem.  It
// would indicate a much deeper issue than just UUID generation. In practice, it almost certainly never going to fail.
// That's why the MustUUID function, which panics on error, is generally considered acceptable in this specific context.
// The panic implies a catastrophic failure of the system's random number generator, which is far more severe than a
// simple UUID generation failure.  A crash due to this problem is arguably preferable to silently generating a
// non-unique or predictable UUID, leading to subtle and hard-to-debug issues.
func MustUUID() uuid.UUID {
	id, err := NewUUID()
	if err != nil {
		panic(err)
	}

	return id
}

// ParseUUID converts a string into a uuid.UUID and returns an error if invalid.
func ParseUUID(input string) (uuid.UUID, error) {
	parsed, err := uuid.Parse(input)
	if err != nil {
		return uuid.Nil, err
	}

	return parsed, nil
}
