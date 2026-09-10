// Package image builds the lab container images from the students' own Dockerfiles and
// saves them as archives a cluster can import.
//
// It talks to a container engine over the Docker API, through the Moby client. That names
// the protocol, not the engine: Podman serves the same API, so pointing DOCKER_HOST at its
// socket builds with Podman instead, with no change here.
package image

import (
	"archive/tar"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/moby/moby/client"
)

// Builder builds lab images through a container engine.
type Builder struct {
	cli *client.Client
}

// NewBuilder connects to the engine named by DOCKER_HOST, or to the local socket.
func NewBuilder() (*Builder, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, fmt.Errorf("connect to the container engine: %w", err)
	}

	return &Builder{cli: cli}, nil
}

// Close releases the connection to the engine.
func (b *Builder) Close() error { return b.cli.Close() }

// Ping reports whether the engine is reachable, so that a missing socket is a clear
// message rather than a confusing build failure.
func (b *Builder) Ping(ctx context.Context) error {
	if _, err := b.cli.Ping(ctx, client.PingOptions{}); err != nil {
		return fmt.Errorf("no container engine at %s: %w", b.cli.DaemonHost(), err)
	}

	return nil
}

// Build builds the image described by dockerfile within contextDir, tags it as ref, and
// writes it to an archive at archivePath. The archive carries the tag, so a cluster
// importing it ends up with an image under exactly that name.
func (b *Builder) Build(ctx context.Context, contextDir, dockerfile, ref, archivePath string) error {
	tarball, err := tarDirectory(contextDir)
	if err != nil {
		return fmt.Errorf("pack the build context %q: %w", contextDir, err)
	}
	defer tarball.Close()
	defer os.Remove(tarball.Name())

	res, err := b.cli.ImageBuild(ctx, tarball, client.ImageBuildOptions{
		Dockerfile: dockerfile,
		Tags:       []string{ref},
		Remove:     true,
	})
	if err != nil {
		return fmt.Errorf("build %s: %w", ref, err)
	}
	defer res.Body.Close()

	if err := readBuildLog(res.Body); err != nil {
		return fmt.Errorf("build %s from %s: %w", ref, filepath.Join(contextDir, dockerfile), err)
	}

	return b.save(ctx, ref, archivePath)
}

// save writes an image out as a tar archive.
func (b *Builder) save(ctx context.Context, ref, archivePath string) error {
	stream, err := b.cli.ImageSave(ctx, []string{ref})
	if err != nil {
		return fmt.Errorf("save %s: %w", ref, err)
	}
	defer stream.Close()

	if err := os.MkdirAll(filepath.Dir(archivePath), 0o755); err != nil {
		return err
	}

	f, err := os.Create(archivePath)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := io.Copy(f, stream); err != nil {
		return fmt.Errorf("write %s: %w", archivePath, err)
	}

	return nil
}

// readBuildLog drains the build stream and reports what went wrong. The engine answers a
// failed build with a successful HTTP response whose body carries the error, so ignoring
// the stream would turn a broken Dockerfile into a silent success.
func readBuildLog(r io.Reader) error {
	type message struct {
		Stream string `json:"stream"`
		Error  string `json:"error"`
	}

	var log strings.Builder
	dec := json.NewDecoder(r)

	for {
		m := message{}
		if err := dec.Decode(&m); err != nil {
			if err == io.EOF {
				return nil
			}
			return fmt.Errorf("reading the build log: %w", err)
		}

		log.WriteString(m.Stream)

		if m.Error != "" {
			return fmt.Errorf("%s\n%s", m.Error, tail(log.String(), 20))
		}
	}
}

// tail returns the last n lines, which is the useful part of a failed build log.
func tail(s string, n int) string {
	lines := strings.Split(strings.TrimRight(s, "\n"), "\n")
	if len(lines) > n {
		lines = lines[len(lines)-n:]
	}

	return strings.Join(lines, "\n")
}

// tarDirectory packs a directory into a temporary tar file, which is how a build context
// reaches the engine.
func tarDirectory(dir string) (*os.File, error) {
	f, err := os.CreateTemp("", "labs-context-*.tar")
	if err != nil {
		return nil, err
	}

	tw := tar.NewWriter(f)

	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		rel, err := filepath.Rel(dir, path)
		if err != nil {
			return err
		}
		if rel == "." {
			return nil
		}

		// Symlinks in a build context are more trouble than they are worth here.
		if !info.Mode().IsRegular() && !info.IsDir() {
			return nil
		}

		hdr, err := tar.FileInfoHeader(info, "")
		if err != nil {
			return err
		}
		hdr.Name = filepath.ToSlash(rel)

		if err := tw.WriteHeader(hdr); err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		src, err := os.Open(path)
		if err != nil {
			return err
		}
		defer src.Close()

		_, err = io.Copy(tw, src)

		return err
	})
	if err != nil {
		f.Close()
		os.Remove(f.Name())

		return nil, err
	}

	if err := tw.Close(); err != nil {
		f.Close()
		os.Remove(f.Name())

		return nil, err
	}

	if _, err := f.Seek(0, io.SeekStart); err != nil {
		f.Close()
		os.Remove(f.Name())

		return nil, err
	}

	return f, nil
}

// Vendor copies the lab modules to a scratch directory and resolves the named module's
// dependencies into vendor/ there.
//
// The copy matters: the lab Dockerfiles build from vendor/, because the modules refer to
// each other through local replace directives and a Dockerfile cannot reach outside its
// build context. Doing it in a scratch copy keeps go.mod, go.sum and the working tree
// exactly as the student left them.
func Vendor(ctx context.Context, codeDir, module string) (string, func(), error) {
	scratch, err := os.MkdirTemp("", "labs-vendor-*")
	if err != nil {
		return "", func() {}, err
	}
	cleanup := func() { os.RemoveAll(scratch) }

	if err := copyTree(codeDir, scratch); err != nil {
		cleanup()
		return "", func() {}, fmt.Errorf("copy the lab modules: %w", err)
	}

	dir := filepath.Join(scratch, module)

	cmd := exec.CommandContext(ctx, "go", "mod", "vendor")
	cmd.Dir = dir

	if out, err := cmd.CombinedOutput(); err != nil {
		cleanup()
		return "", func() {}, fmt.Errorf("go mod vendor in %s: %w\n%s", module, err, out)
	}

	return dir, cleanup, nil
}

func copyTree(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		target := filepath.Join(dst, rel)

		switch {
		case info.IsDir():
			return os.MkdirAll(target, 0o755)
		case !info.Mode().IsRegular():
			return nil
		}

		in, err := os.Open(path)
		if err != nil {
			return err
		}
		defer in.Close()

		out, err := os.OpenFile(target, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, info.Mode())
		if err != nil {
			return err
		}
		defer out.Close()

		_, err = io.Copy(out, in)

		return err
	})
}
