package common_test

import (
	"testing"

	"github.com/pulumi/pulumi-command/provider/pkg/provider/common"
	"github.com/stretchr/testify/assert"
)

func TestShouldLog(t *testing.T) {
	for _, tc := range []struct {
		logging                    common.Logging
		expectStdout, expectStderr bool
	}{
		{common.LogStdoutAndStderr, true, true},
		{common.LogStdout, true, false},
		{common.LogStderr, false, true},
		{common.NoLogging, false, false},
	} {
		t.Run(string(tc.logging), func(t *testing.T) {
			assert.Equal(t, tc.expectStdout, tc.logging.ShouldLogStdout())
			assert.Equal(t, tc.expectStderr, tc.logging.ShouldLogStderr())
		})
	}
}
