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

package gcp

import (
	"context"
	"net/http"
	"regexp"
)

type (
	// TracingContext defines a type for the Cloud Trace Context key.
	TracingContext string
)

const (
	// CloudTraceContextHeader is the header name for the Cloud Trace Context.
	CloudTraceContextHeader = "X-Cloud-Trace-Context"

	// TraceContextKey is the key for storing the trace ID in the request context.
	TraceContextKey TracingContext = "trace"

	// SpanContextKey is the key for storing the span ID in the request context.
	SpanContextKey TracingContext = "span"

	// SampledContextKey is the key for storing the trace sampling flag in the request context.
	SampledContextKey TracingContext = "trace_sampled"
)

var (
	// match defines a regular expression for parsing the Cloud Trace Context header.
	match = regexp.MustCompile(
		// Matches on "TRACE_ID"
		`([a-f\d]+)?` +
			// Matches on "/SPAN_ID"
			`(?:/([a-f\d]+))?` +
			// Matches on ";0=TRACE_TRUE"
			`(?:;o=(\d))?`)
)

// CloudTraceReporter extracts trace, span, and sampling flags from a HTTP request
// header and stores them in the context.
//
// It parses the `X-Cloud-Trace-Context` header and stores the extracted
// values in the context using the `TraceContextKey`, `SpanContextKey`, and
// `SampledContextKey` keys, respectively.
func CloudTraceReporter(req *http.Request) *http.Request {
	header := req.Header.Get(CloudTraceContextHeader)
	if header != "" && match != nil {
		matches := match.FindStringSubmatch(header)
		trace, span, sampled := matches[1], matches[2], matches[3] == "1"

		if span == "0" {
			span = ""
		}

		ctx := req.Context()
		ctx = context.WithValue(ctx, TraceContextKey, trace)
		ctx = context.WithValue(ctx, SpanContextKey, span)
		ctx = context.WithValue(ctx, SampledContextKey, sampled)

		req = req.WithContext(ctx)
	}

	return req
}
