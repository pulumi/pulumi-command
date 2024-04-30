package remote

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/gliderlabs/ssh"
	"github.com/pkg/sftp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	xssh "golang.org/x/crypto/ssh"
)

func testSftpHandler(t *testing.T, baseDir string, sess ssh.Session) {
	server, err := sftp.NewServer(sess, sftp.WithServerWorkingDirectory(baseDir))
	require.NoError(t, err)

	if err := server.Serve(); err == io.EOF {
		server.Close()
		fmt.Println("sftp client exited session.")
	}
}

func TestCopyDirectories(t *testing.T) {
	// Start a local SSH and SFTP server that writes files to the local file system, under baseDir.
	const serverAddr = "127.0.0.1:3333"
	baseDir := t.TempDir()

	server := ssh.Server{
		Addr: serverAddr,
		SubsystemHandlers: map[string]ssh.SubsystemHandler{
			"sftp": func(s ssh.Session) { testSftpHandler(t, baseDir, s) },
		},
	}
	go func() {
		// "ListenAndServe always returns a non-nil error."
		_ = server.ListenAndServe()
	}()
	t.Cleanup(func() {
		_ = server.Close()
	})

	err := os.MkdirAll(filepath.Join(baseDir, "src", "one", "two"), 0755)
	require.NoError(t, err)
	_, err = os.Create(filepath.Join(baseDir, "src", "file1"))
	require.NoError(t, err)
	_, err = os.Create(filepath.Join(baseDir, "src", "one", "file2"))
	require.NoError(t, err)
	_, err = os.Create(filepath.Join(baseDir, "src", "one", "two", "file3"))
	require.NoError(t, err)

	sshClient, err := xssh.Dial("tcp", serverAddr, &xssh.ClientConfig{
		HostKeyCallback: xssh.InsecureIgnoreHostKey(),
	})
	require.NoError(t, err)
	sftpClient, err := sftp.NewClient(sshClient)
	require.NoError(t, err)

	err = copyDir(sftpClient, filepath.Join(baseDir, "src"), filepath.Join(baseDir, "dest"))
	require.NoError(t, err)

	assert.FileExists(t, filepath.Join(baseDir, "dest", "file1"))
	assert.FileExists(t, filepath.Join(baseDir, "dest", "one", "file2"))
	assert.FileExists(t, filepath.Join(baseDir, "dest", "one", "two", "file3"))
}
