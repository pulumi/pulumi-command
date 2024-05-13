package remote

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource/archive"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource/asset"

	"github.com/gliderlabs/ssh"
	"github.com/pkg/sftp"
	p "github.com/pulumi/pulumi-go-provider"
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

func TestCheck(t *testing.T) {
	host := "host"
	validConnection := &Connection{
		connectionBase: connectionBase{Host: &host},
	}

	copy := &Copy{}

	makeNewInput := func(asset *asset.Asset, archive *archive.Archive) CopyInputs {
		return CopyInputs{
			Connection:   validConnection,
			LocalAsset:   asset,
			LocalArchive: archive,
			RemotePath:   "path/to/remote",
		}
	}

	checkNoError := func(news CopyInputs) []p.CheckFailure {
		newsRaw := resource.NewPropertyMap(news)
		_, failures, err := copy.Check(nil, "urn", nil, newsRaw)
		require.NoError(t, err)
		return failures
	}

	t.Run("happy path, asset", func(t *testing.T) {
		news := makeNewInput(&asset.Asset{Path: "path/to/file"}, nil)
		failures := checkNoError(news)
		assert.Empty(t, failures)
	})

	t.Run("happy path, archive", func(t *testing.T) {
		news := makeNewInput(nil, &archive.Archive{Path: "path/to/file"})
		failures := checkNoError(news)
		assert.Empty(t, failures)
	})

	t.Run("asset or archive, not both", func(t *testing.T) {
		news := makeNewInput(&asset.Asset{Path: "path/to/file"}, &archive.Archive{Path: "path/to/file"})
		failures := checkNoError(news)
		assert.Len(t, failures, 1)
	})

	t.Run("need asset or archive", func(t *testing.T) {
		news := makeNewInput(nil, nil)
		failures := checkNoError(news)
		assert.Len(t, failures, 1)
	})

	t.Run("asset must be path-based", func(t *testing.T) {
		news := makeNewInput(&asset.Asset{URI: "http://example.com"}, nil)
		failures := checkNoError(news)
		assert.Len(t, failures, 1)
	})

	t.Run("archive must be path-based", func(t *testing.T) {
		news := makeNewInput(nil, &archive.Archive{URI: "http://example.com"})
		failures := checkNoError(news)
		assert.Len(t, failures, 1)
	})

	t.Run("can diagnose multiple issues", func(t *testing.T) {
		news := makeNewInput(&asset.Asset{URI: "http://example.com"}, &archive.Archive{URI: "http://example.com"})
		failures := checkNoError(news)
		assert.Len(t, failures, 3)
	})
}
