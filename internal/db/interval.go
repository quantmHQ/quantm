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

package db

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"google.golang.org/protobuf/types/known/durationpb"
)

// IntervalToDuration converts a pgtype.Interval to a time.Duration.
func IntervalToDuration(interval pgtype.Interval) time.Duration {
	ms := interval.Microseconds +
		int64(interval.Days*24*60*60*1000*1000) +
		int64(interval.Months*30*24*60*60*1000*1000)

	return time.Duration(ms) * time.Microsecond
}

// IntervalToProto converts a pgtype.Interval to a durationpb.Duration.
func IntervalToProto(interval pgtype.Interval) *durationpb.Duration {
	return durationpb.New(IntervalToDuration(interval))
}

// DurationToInterval converts a time.Duration to a pgtype.Interval.
func DurationToInterval(d time.Duration) pgtype.Interval {
	return pgtype.Interval{
		Microseconds: int64(d / time.Microsecond),
	}
}

// ProtoToInterval converts a durationpb.Duration to a pgtype.Interval.
func ProtoToInterval(d *durationpb.Duration) pgtype.Interval {
	return pgtype.Interval{
		Microseconds: int64(d.AsDuration() / time.Microsecond),
	}
}
