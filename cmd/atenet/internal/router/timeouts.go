// Copyright 2026 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package router

import (
	"os"
	"strconv"
	"time"
)

// Resume / forwarding timeouts and their environment overrides.
//
// Resume-on-demand of a suspended actor performs a cold gVisor restore, which
// for a large snapshot (tens of MiB) can take tens of seconds. The default
// route / ext_proc / background-resume timeouts are sized for steady-state
// traffic and can cancel an in-flight restore, surfacing as a 504. These knobs
// let an operator running heavy actors raise the ceilings without a code change.
//
// Defaults are unchanged from prior behavior, so this is purely additive.
//
// Longer term, suspend-safe actor networking (agent-substrate/substrate#465)
// should remove most of the need to tune these.
const (
	resumeTimeoutEnv  = "ATE_RESUME_TIMEOUT_SECONDS"
	routeTimeoutEnv   = "ATE_ROUTE_TIMEOUT_SECONDS"
	extProcTimeoutEnv = "ATE_EXTPROC_TIMEOUT_SECONDS"

	defaultResumeTimeout  = 15 * time.Second
	defaultRouteTimeout   = 10 * time.Second
	defaultExtProcTimeout = 5 * time.Second
)

// timeoutFromEnv returns the duration from a whole-seconds env var, falling back
// to def when the var is unset or not a positive integer.
func timeoutFromEnv(key string, def time.Duration) time.Duration {
	if v := os.Getenv(key); v != "" {
		if secs, err := strconv.Atoi(v); err == nil && secs > 0 {
			return time.Duration(secs) * time.Second
		}
	}
	return def
}

// resumeTimeout bounds a detached, background resume operation.
func resumeTimeout() time.Duration { return timeoutFromEnv(resumeTimeoutEnv, defaultResumeTimeout) }

// routeTimeout bounds a forwarded upstream request (e.g. a long LLM turn).
func routeTimeout() time.Duration { return timeoutFromEnv(routeTimeoutEnv, defaultRouteTimeout) }

// extProcTimeout bounds the ext_proc round-trip, which can drive a cold restore.
func extProcTimeout() time.Duration { return timeoutFromEnv(extProcTimeoutEnv, defaultExtProcTimeout) }
