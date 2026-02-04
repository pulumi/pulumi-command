package local_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/pulumi/pulumi-command/provider/pkg/provider/local"
)

func TestShouldLog(t *testing.T) {
	for _, tc := range []struct {
		logging                    local.Logging
		expectStdout, expectStderr bool
	}{
		{local.LogStdoutAndStderr, true, true},
		{local.LogStdout, true, false},
		{local.LogStderr, false, true},
		{local.NoLogging, false, false},
	} {
		t.Run(string(tc.logging), func(t *testing.T) {
			assert.Equal(t, tc.expectStdout, tc.logging.ShouldLogStdout())
			assert.Equal(t, tc.expectStderr, tc.logging.ShouldLogStderr())
		})
	}
}
