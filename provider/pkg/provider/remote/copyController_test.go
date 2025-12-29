package remote

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/gliderlabs/ssh"
	"github.com/pkg/sftp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	xssh "golang.org/x/crypto/ssh"

	p "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi-go-provider/infer"
	"github.com/pulumi/pulumi-go-provider/infer/types"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource/archive"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource/asset"
	"github.com/pulumi/pulumi/sdk/v3/go/property"
)

func testSftpHandler(t *testing.T, baseDir string, sess ssh.Session) {
	server, err := sftp.NewServer(sess, sftp.WithServerWorkingDirectory(baseDir))
	require.NoError(t, err)

	if err := server.Serve(); err == io.EOF {
		server.Close()
		fmt.Println("sftp client exited session.")
	}
}

// Start a local SSH and SFTP server that writes files to the local file system, under baseDir.
func startSSHServer(t *testing.T, baseDir string) *sftp.Client {
	serverAddr := "127.0.0.1:3333"

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

	// Wait until SSH server is up
	var sshClient *xssh.Client
	var err error
	for i := 0; i < 20; i++ {
		sshClient, err = xssh.Dial("tcp", serverAddr, &xssh.ClientConfig{
			//nolint:gosec // G106: InsecureIgnoreHostKey is acceptable in tests
			HostKeyCallback: xssh.InsecureIgnoreHostKey(),
		})
		if err == nil {
			fmt.Printf("SSH server is up at attempt %d\n", i)
			break
		}
		time.Sleep(100 * time.Millisecond)
	}
	require.NoError(t, err)

	sftpClient, err := sftp.NewClient(sshClient)
	require.NoError(t, err)
	return sftpClient
}

func initCopyTest(t *testing.T) (srcDir, destDir string, sftpClient *sftp.Client) {
	baseDir := t.TempDir()

	destDir = filepath.Join(baseDir, "dest")
	require.NoError(t, os.Mkdir(destDir, 0o755))

	sftpClient = startSSHServer(t, destDir)

	srcDirName := "src"
	srcDir = filepath.Join(baseDir, srcDirName)

	// Our test directory structure:
	// file1
	// one/
	//   file2
	//   two/
	//     file3
	require.NoError(t, os.MkdirAll(filepath.Join(srcDir, "one", "two"), 0o755))
	_, err := os.Create(filepath.Join(srcDir, "file1"))
	require.NoError(t, err)
	_, err = os.Create(filepath.Join(srcDir, "one", "file2"))
	require.NoError(t, err)
	_, err = os.Create(filepath.Join(srcDir, "one", "two", "file3"))
	require.NoError(t, err)

	return srcDir, destDir, sftpClient
}

// assertDirectoryTree asserts that the directory structure under baseDir matches the structure
// that initCopyTest creates.
func assertDirectoryTree(t *testing.T, baseDir string) {
	assert.FileExists(t, filepath.Join(baseDir, "file1"))
	assert.FileExists(t, filepath.Join(baseDir, "one", "file2"))
	assert.FileExists(t, filepath.Join(baseDir, "one", "two", "file3"))

	// No other files or directories should exist.
	b, err := os.ReadDir(baseDir)
	require.NoError(t, err)
	assert.Len(t, b, 2)

	b, err = os.ReadDir(filepath.Join(baseDir, "one"))
	require.NoError(t, err)
	assert.Len(t, b, 2)

	b, err = os.ReadDir(filepath.Join(baseDir, "one", "two"))
	require.NoError(t, err)
	assert.Len(t, b, 1)
}

func TestCopyDirectories(t *testing.T) {
	t.Run("copy file into directory", func(t *testing.T) {
		srcDir, destDir, sftpClient := initCopyTest(t)
		require.NoError(t, sftpCopy(sftpClient, filepath.Join(srcDir, "file1"), destDir))
		assert.FileExists(t, filepath.Join(destDir, "file1"))
	})

	t.Run("copy file to file", func(t *testing.T) {
		srcDir, destDir, sftpClient := initCopyTest(t)
		dest := filepath.Join(destDir, "remoteFile")
		require.NoError(t, sftpCopy(sftpClient, filepath.Join(srcDir, "file1"), dest))
		assert.FileExists(t, dest)
	})

	t.Run("copy dir recursively", func(t *testing.T) {
		srcDir, destDir, sftpClient := initCopyTest(t)
		require.NoError(t, sftpCopy(sftpClient, srcDir, destDir))
		assertDirectoryTree(t, filepath.Join(destDir, filepath.Base(srcDir)))
	})

	t.Run("copy dir contents recursively", func(t *testing.T) {
		srcDir, destDir, sftpClient := initCopyTest(t)
		require.NoError(t, sftpCopy(sftpClient, srcDir+"/", destDir))
		assertDirectoryTree(t, destDir)
	})

	t.Run("copy dir then no-op update", func(t *testing.T) {
		srcDir, destDir, sftpClient := initCopyTest(t)
		require.NoError(t, sftpCopy(sftpClient, srcDir, destDir))
		assertDirectoryTree(t, filepath.Join(destDir, filepath.Base(srcDir)))

		require.NoError(t, sftpCopy(sftpClient, srcDir, destDir))
	})

	t.Run("don't replace file with directory", func(t *testing.T) {
		srcDir, destDir, sftpClient := initCopyTest(t)
		require.NoError(t, sftpCopy(sftpClient, srcDir, destDir))
		assertDirectoryTree(t, filepath.Join(destDir, filepath.Base(srcDir)))

		fileTwo := filepath.Join(destDir, "src", "one", "two")
		require.NoError(t, os.RemoveAll(fileTwo))
		require.NoError(t, os.WriteFile(fileTwo, []byte("dir turned to file"), 0o600))

		require.Error(t, sftpCopy(sftpClient, srcDir, destDir))
	})

	t.Run("wildcards are not supported", func(t *testing.T) {
		srcDir, destDir, sftpClient := initCopyTest(t)
		require.Error(t, sftpCopy(sftpClient, filepath.Join(srcDir, "file*"), destDir))
	})

	t.Run("overwrite file", func(t *testing.T) {
		srcDir, destDir, sftpClient := initCopyTest(t)

		require.NoError(t, sftpCopy(sftpClient, srcDir, destDir))
		destFile := filepath.Join(destDir, "src", "file1")
		assert.FileExists(t, destFile)

		// modify the file
		srcFile := filepath.Join(srcDir, "file1")
		require.NoError(t, os.WriteFile(srcFile, []byte("new content"), 0o600))

		// copy it to remote again
		require.NoError(t, sftpCopy(sftpClient, srcFile, destFile))
		content, err := os.ReadFile(destFile)
		require.NoError(t, err)
		assert.Equal(t, "new content", string(content))
	})

	t.Run("overwrite file copying dir", func(t *testing.T) {
		srcDir, destDir, sftpClient := initCopyTest(t)

		require.NoError(t, sftpCopy(sftpClient, srcDir, destDir))
		destFile := filepath.Join(destDir, "src", "file1")
		assert.FileExists(t, destFile)

		// modify the file
		srcFile := filepath.Join(srcDir, "file1")
		require.NoError(t, os.WriteFile(srcFile, []byte("new content"), 0o600))

		// copy it to remote again
		require.NoError(t, sftpCopy(sftpClient, srcDir, destDir))
		content, err := os.ReadFile(destFile)
		require.NoError(t, err)
		assert.Equal(t, "new content", string(content))
	})

	t.Run("overwrite file copying dir contents", func(t *testing.T) {
		srcDir, destDir, sftpClient := initCopyTest(t)

		require.NoError(t, sftpCopy(sftpClient, srcDir+"/", destDir))
		destFile := filepath.Join(destDir, "file1")
		assert.FileExists(t, destFile)

		// modify the file
		srcFile := filepath.Join(srcDir, "file1")
		require.NoError(t, os.WriteFile(srcFile, []byte("new content"), 0o600))

		// copy it to remote again
		require.NoError(t, sftpCopy(sftpClient, srcDir+"/", destDir))
		content, err := os.ReadFile(destFile)
		require.NoError(t, err)
		assert.Equal(t, "new content", string(content))
	})
}

func TestCheck(t *testing.T) {
	makeNewInput := func(asset *asset.Asset, archive *archive.Archive) property.Map {
		m := map[string]any{
			"connection": map[string]any{
				"host": "myhost",
			},
			"remotePath": "path/to/remote",
			// AssetOrArchive has special handling in pulumi-go-provider and is kept as a primitive.
			"source": types.AssetOrArchive{
				Asset:   asset,
				Archive: archive,
			},
		}
		pm := resource.NewPropertyMapFromMap(m)
		return resource.FromResourcePropertyMap(pm)
	}

	check := func(news property.Map) []p.CheckFailure {
		ctr := &CopyToRemote{}
		resp, err := ctr.Check(context.Background(), infer.CheckRequest{Name: "name", NewInputs: news})
		require.NoError(t, err)
		return resp.Failures
	}

	t.Run("happy path, asset", func(t *testing.T) {
		news := makeNewInput(&asset.Asset{Path: "path/to/file"}, nil)
		failures := check(news)
		assert.Empty(t, failures)
	})

	t.Run("happy path, archive", func(t *testing.T) {
		news := makeNewInput(nil, &archive.Archive{Path: "path/to/file"})
		failures := check(news)
		assert.Empty(t, failures)
	})

	t.Run("asset or archive, not both", func(t *testing.T) {
		news := makeNewInput(&asset.Asset{Path: "path/to/file"}, &archive.Archive{Path: "path/to/file"})
		failures := check(news)
		assert.Len(t, failures, 1)
	})

	t.Run("need asset or archive", func(t *testing.T) {
		news := makeNewInput(nil, nil)
		failures := check(news)
		assert.Len(t, failures, 1)
	})

	t.Run("asset must be path-based", func(t *testing.T) {
		news := makeNewInput(&asset.Asset{URI: "http://example.com"}, nil)
		failures := check(news)
		assert.Len(t, failures, 1)
	})

	t.Run("archive must be path-based", func(t *testing.T) {
		news := makeNewInput(nil, &archive.Archive{URI: "http://example.com"})
		failures := check(news)
		assert.Len(t, failures, 1)
	})

	t.Run("can diagnose multiple issues", func(t *testing.T) {
		news := makeNewInput(&asset.Asset{URI: "http://example.com"}, &archive.Archive{URI: "http://example.com"})
		failures := check(news)
		assert.Len(t, failures, 3)
	})

	t.Run("happy path, text asset", func(t *testing.T) {
		news := makeNewInput(&asset.Asset{Text: "hello world"}, nil)
		failures := check(news)
		assert.Empty(t, failures)
	})
}

func TestCopyTextContent(t *testing.T) {
	t.Run("copy text content to file", func(t *testing.T) {
		baseDir := t.TempDir()
		destDir := filepath.Join(baseDir, "dest")
		require.NoError(t, os.Mkdir(destDir, 0o755))

		sftpClient := startSshServer(t, destDir)

		textContent := "hello from text asset"
		destFile := "textfile.txt"
		require.NoError(t, copyTextContent(sftpClient, textContent, destFile))

		content, err := os.ReadFile(filepath.Join(destDir, destFile))
		require.NoError(t, err)
		assert.Equal(t, textContent, string(content))
	})

	t.Run("overwrite existing file with text content", func(t *testing.T) {
		baseDir := t.TempDir()
		destDir := filepath.Join(baseDir, "dest")
		require.NoError(t, os.Mkdir(destDir, 0o755))

		sftpClient := startSshServer(t, destDir)

		// Create an existing file
		destFile := "existing.txt"
		require.NoError(t, os.WriteFile(filepath.Join(destDir, destFile), []byte("old content"), 0o644))

		// Overwrite with text content
		textContent := "new content from text asset"
		require.NoError(t, copyTextContent(sftpClient, textContent, destFile))

		content, err := os.ReadFile(filepath.Join(destDir, destFile))
		require.NoError(t, err)
		assert.Equal(t, textContent, string(content))
	})

	t.Run("error when destination is directory", func(t *testing.T) {
		baseDir := t.TempDir()
		destDir := filepath.Join(baseDir, "dest")
		require.NoError(t, os.Mkdir(destDir, 0o755))

		sftpClient := startSshServer(t, destDir)

		// Create a subdirectory
		subDir := "subdir"
		require.NoError(t, os.Mkdir(filepath.Join(destDir, subDir), 0o755))

		// Trying to copy text content to a directory should fail
		textContent := "hello"
		err := copyTextContent(sftpClient, textContent, subDir)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "is a directory")
	})
}
