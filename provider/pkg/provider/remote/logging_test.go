package remote_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/pulumi/pulumi-command/provider/pkg/provider/remote"
)

func TestShouldLog(t *testing.T) {
	for _, tc := range []struct {
		logging                    remote.Logging
		expectStdout, expectStderr bool
	}{
		{remote.LogStdoutAndStderr, true, true},
		{remote.LogStdout, true, false},
		{remote.LogStderr, false, true},
		{remote.NoLogging, false, false},
	} {
		t.Run(string(tc.logging), func(t *testing.T) {
			assert.Equal(t, tc.expectStdout, tc.logging.ShouldLogStdout())
			assert.Equal(t, tc.expectStderr, tc.logging.ShouldLogStderr())
		})
	}
}
