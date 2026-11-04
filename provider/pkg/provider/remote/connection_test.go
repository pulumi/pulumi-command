// Copyright 2026, Pulumi Corporation.
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

package remote

import (
	"context"
	"errors"
	"net"
	"sync/atomic"
	"testing"
	"testing/synctest"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func TestDialWithRetry(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		wait       *int
		limit      int
		succeedOn  int
		wantDelays []time.Duration
	}{
		{
			name: "omitted preserves default backoff and cap", limit: 12,
			wantDelays: []time.Duration{
				100 * time.Millisecond, 150 * time.Millisecond, 225 * time.Millisecond,
				337500 * time.Microsecond, 506250 * time.Microsecond, 759375 * time.Microsecond,
				1139062500 * time.Nanosecond, 1708593750 * time.Nanosecond,
				2562890625 * time.Nanosecond, 3844335937 * time.Nanosecond,
				5 * time.Second, 5 * time.Second,
			},
		},
		{
			name: "fixed wait", wait: pulumi.IntRef(5), limit: 3,
			wantDelays: []time.Duration{5 * time.Second, 5 * time.Second, 5 * time.Second},
		},
		{
			name: "fixed wait above default cap", wait: pulumi.IntRef(10), limit: 2,
			wantDelays: []time.Duration{10 * time.Second, 10 * time.Second},
		},
		{
			name: "zero retries immediately", wait: pulumi.IntRef(0), limit: 3,
			wantDelays: []time.Duration{0, 0, 0},
		},
		{name: "success on first attempt", wait: pulumi.IntRef(5), limit: 3, succeedOn: 1},
		{
			name: "success after retries", wait: pulumi.IntRef(2), limit: 3, succeedOn: 3,
			wantDelays: []time.Duration{2 * time.Second, 2 * time.Second},
		},
		{name: "zero error limit stops on first failure", wait: pulumi.IntRef(5), limit: 0},
		{
			name: "unlimited errors", wait: pulumi.IntRef(1), limit: -1, succeedOn: 13,
			wantDelays: []time.Duration{
				time.Second, time.Second, time.Second, time.Second, time.Second, time.Second,
				time.Second, time.Second, time.Second, time.Second, time.Second, time.Second,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			synctest.Test(t, func(t *testing.T) {
				dialErr := errors.New("dial failed")
				start := time.Now()
				var attempts []time.Time
				result, err := dialWithRetry(t.Context(), "Dial", tt.limit, tt.wait, func() (int, error) {
					attempts = append(attempts, time.Now())
					if len(attempts) == tt.succeedOn {
						return 42, nil
					}
					return 0, dialErr
				})
				if tt.succeedOn > 0 {
					require.NoError(t, err)
					require.Equal(t, 42, result)
				} else {
					require.ErrorIs(t, err, dialErr)
					require.Equal(t, tt.limit+1, len(attempts))
				}
				require.Len(t, attempts, len(tt.wantDelays)+1)
				require.Equal(t, start, attempts[0], "first attempt must not wait")
				for i, delay := range tt.wantDelays {
					require.Equal(t, delay, attempts[i+1].Sub(attempts[i]))
				}
				require.Equal(t, attempts[len(attempts)-1], time.Now(), "must not wait after the last attempt")
			})
		})
	}
}

func TestDialWithRetryCancellation(t *testing.T) {
	t.Parallel()
	for _, wait := range []*int{nil, pulumi.IntRef(5)} {
		synctest.Test(t, func(t *testing.T) {
			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()
			go func() {
				time.Sleep(50 * time.Millisecond)
				cancel()
			}()
			start := time.Now()
			attempts := 0
			_, err := dialWithRetry(ctx, "Dial", -1, wait, func() (int, error) {
				attempts++
				return 0, errors.New("dial failed")
			})
			require.ErrorIs(t, err, context.Canceled)
			require.Equal(t, 1, attempts)
			require.Equal(t, 50*time.Millisecond, time.Since(start))
		})
	}
}

func TestConnectionDialRetryWait(t *testing.T) {
	t.Parallel()
	for _, tt := range []struct {
		name         string
		proxy        bool
		wait         *int
		wantDuration time.Duration
	}{
		{name: "direct", wait: pulumi.IntRef(2), wantDuration: 6 * time.Second},
		{name: "proxy", proxy: true, wait: pulumi.IntRef(2), wantDuration: 6 * time.Second},
		{name: "proxy omitted", proxy: true, wantDuration: 475 * time.Millisecond},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			// Keep real network I/O outside the bubble so it does not prevent the virtual clock from advancing.
			listener, err := net.Listen("tcp", "127.0.0.1:0")
			require.NoError(t, err)
			defer listener.Close()
			var attempts atomic.Int32
			go func() {
				for {
					conn, err := listener.Accept()
					if err != nil {
						return
					}
					attempts.Add(1)
					_ = conn.Close()
				}
			}()
			synctest.Test(t, func(t *testing.T) {
				base := connectionBase{
					Host: pulumi.StringRef("127.0.0.1"), User: pulumi.StringRef("user"),
					Port:           pulumi.Float64Ref(float64(listener.Addr().(*net.TCPAddr).Port)),
					PerDialTimeout: pulumi.IntRef(1), DialErrorLimit: pulumi.IntRef(3), DialRetryWait: tt.wait,
				}
				connection := Connection{connectionBase: base}
				if tt.proxy {
					connection.Proxy = &ProxyConnection{connectionBase: base}
					connection.DialErrorLimit = pulumi.IntRef(0)
					connection.DialRetryWait = pulumi.IntRef(5)
				}
				config, err := connection.SSHConfig()
				require.NoError(t, err)
				require.Equal(t, time.Second, config.Timeout)
				start := time.Now()
				_, err = connection.Dial(t.Context())
				require.Error(t, err)
				if tt.proxy {
					require.ErrorContains(t, err, "proxy:")
				}
				require.Equal(t, int32(4), attempts.Load())
				require.Equal(t, tt.wantDuration, time.Since(start))
			})
		})
	}
}
