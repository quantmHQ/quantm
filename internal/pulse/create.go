// Crafted with ❤ at Breu, Inc. <info@breu.io>, Copyright © 2024, 2025.
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

package pulse

import (
	"context"
	"fmt"

	"github.com/gobeam/stringy"
)

const (
	statement__events__create = `
CREATE TABLE IF NOT EXISTS %s (
  version String,
  id UUID,
  parents Array(UUID),
  hook Int32,
  scope String,
  action String,
  source String,
  subject_id UUID,
  subject_name String,
  user_id UUID,
  team_id UUID,
  org_id UUID,
  timestamp DateTime
)
ENGINE = MergeTree()
PARTITION BY toYYYYMM(timestamp)
ORDER BY (toStartOfWeek(timestamp), toStartOfMonth(timestamp), timestamp, id);
`
)

// table_name returns the table name for the given kind and slug.
func table_name(kind, slug string) string {
	table := fmt.Sprintf("%s_%s", kind, slug)

	return stringy.New(table).SnakeCase().Get()
}

func CreateEventsTable(ctx context.Context, slug string) error {
	table := table_name("events", slug)
	stmt := fmt.Sprintf(statement__events__create, table)

	return Get().Connection().Exec(ctx, stmt)
}
