package remote

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/pulumi/pulumi-go-provider/infer/types"
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
	srcDirName := "src"
	srcDir := filepath.Join(baseDir, srcDirName)

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

	// Our test directory structure:
	// file1
	// one/
	//   file2
	//   two/
	//     file3
	require.NoError(t, os.MkdirAll(filepath.Join(srcDir, "one", "two"), 0755))
	_, err := os.Create(filepath.Join(srcDir, "file1"))
	require.NoError(t, err)
	_, err = os.Create(filepath.Join(srcDir, "one", "file2"))
	require.NoError(t, err)
	_, err = os.Create(filepath.Join(srcDir, "one", "two", "file3"))
	require.NoError(t, err)

	assertDirectoryTree := func(t *testing.T, base string) {
		assert.FileExists(t, filepath.Join(base, "file1"))
		assert.FileExists(t, filepath.Join(base, "one", "file2"))
		assert.FileExists(t, filepath.Join(base, "one", "two", "file3"))
	}

	sshClient, err := xssh.Dial("tcp", serverAddr, &xssh.ClientConfig{
		HostKeyCallback: xssh.InsecureIgnoreHostKey(),
	})
	require.NoError(t, err)
	sftpClient, err := sftp.NewClient(sshClient)
	require.NoError(t, err)

	t.Run("copy dir recursively", func(t *testing.T) {
		dest := filepath.Join(t.TempDir(), "dest")
		require.NoError(t, sftpCopy(sftpClient, srcDir, dest))
		assertDirectoryTree(t, filepath.Join(dest, srcDirName))
	})

	t.Run("copy dir contents recursively", func(t *testing.T) {
		dest := filepath.Join(t.TempDir(), "dest")
		require.NoError(t, sftpCopy(sftpClient, srcDir+"/", dest))
		assertDirectoryTree(t, dest)
	})

	t.Run("copy file into directory", func(t *testing.T) {
		dest := t.TempDir()
		require.NoError(t, sftpCopy(sftpClient, filepath.Join(srcDir, "file1"), dest))
		assert.FileExists(t, filepath.Join(dest, "file1"))
	})

	t.Run("copy file to given remote path", func(t *testing.T) {
		dest := filepath.Join(t.TempDir(), "remoteFile")
		require.NoError(t, sftpCopy(sftpClient, filepath.Join(srcDir, "file1"), dest))
		assert.FileExists(t, dest)
	})

	t.Run("overwrite file", func(t *testing.T) {
		tmp := t.TempDir()

		// create a file
		srcDir := filepath.Join(tmp, "src")
		srcFile := filepath.Join(srcDir, "file1")
		require.NoError(t, os.Mkdir(srcDir, 0755))
		_, err = os.Create(srcFile)
		require.NoError(t, err)

		// copy it to remote
		dest := filepath.Join(tmp, "dest")
		require.NoError(t, sftpCopy(sftpClient, srcDir+"/", dest))
		assert.FileExists(t, filepath.Join(dest, "file1"))

		// modify the file
		require.NoError(t, os.WriteFile(srcFile, []byte("new content"), 0644))

		// copy it to remote again
		require.NoError(t, sftpCopy(sftpClient, srcDir, dest))
		content, err := os.ReadFile(srcFile)
		require.NoError(t, err)
		assert.Equal(t, "new content", string(content))
	})

	t.Run("don't replace file with directory", func(t *testing.T) {
		dest := t.TempDir()
		require.NoError(t, sftpCopy(sftpClient, srcDir, dest))
		assertDirectoryTree(t, filepath.Join(dest, "src"))

		fileTwo := filepath.Join(dest, "src", "one", "two")
		require.NoError(t, os.RemoveAll(fileTwo))
		require.NoError(t, os.WriteFile(fileTwo, []byte("dir turned to file"), 0644))

		require.Error(t, sftpCopy(sftpClient, srcDir, dest))
	})

	t.Run("wildcards are not supported", func(t *testing.T) {
		require.Error(t, sftpCopy(sftpClient, filepath.Join(srcDir, "file*"), t.TempDir()))
	})
}

func TestCheck(t *testing.T) {
	host := "host"
	validConnection := &Connection{
		connectionBase: connectionBase{Host: &host},
	}

	copy := &Copy{}

	makeNewInput := func(asset *asset.Asset, archive *archive.Archive) CopyInputs {
		aa := types.AssetOrArchive{}
		if asset != nil {
			aa.Asset = asset
		} else if archive != nil {
			aa.Archive = archive
		}
		return CopyInputs{
			Connection: validConnection,
			Source:     aa,
			RemotePath: "path/to/remote",
		}
	}

	checkNoError := func(news CopyInputs) []p.CheckFailure {
		newsRaw := resource.NewPropertyMap(news)
		_, failures, err := copy.Check(context.Background(), "urn", nil, newsRaw)
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
