package image

import (
	"archive/tar"
	"context"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"testing"
	"time"

	"labci/pkg/labenv"
)

// TestBuildLabImage builds a real lab image from the student's own Dockerfile and checks
// that the archive it produces carries the tag. The tag is the whole point: a cluster
// importing an untagged archive ends up with an image no pod can refer to.
func TestBuildLabImage(t *testing.T) {
	root, err := labenv.FindRoot()
	if err != nil {
		t.Fatalf("%s", err)
	}
	code := filepath.Join(root, "99-labs", "code")

	module := os.Getenv("LABS_MODULE")
	if module == "" {
		module = "splitdim"
	}
	ref := "localhost/" + module + ":latest"

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	b, err := NewBuilder()
	if err != nil {
		t.Fatalf("%s", err)
	}
	defer b.Close()

	if err := b.Ping(ctx); err != nil {
		t.Fatalf("%s", err)
	}

	dir, cleanup, err := Vendor(ctx, code, module)
	if err != nil {
		t.Fatalf("%s", err)
	}
	defer cleanup()

	// Vendoring has to happen in the scratch copy: go.mod, go.sum and the working tree
	// belong to the student and must come back exactly as they were.
	if _, err := os.Stat(filepath.Join(code, module, "vendor")); err == nil {
		t.Errorf("vendoring wrote into the student's tree at %s/vendor", module)
	}

	archive := filepath.Join(t.TempDir(), module+".tar")
	if err := b.Build(ctx, dir, "deploy/Dockerfile", ref, archive); err != nil {
		t.Fatalf("%s", err)
	}

	info, err := os.Stat(archive)
	if err != nil {
		t.Fatalf("%s", err)
	}
	t.Logf("built %s -> %s (%d bytes)", ref, filepath.Base(archive), info.Size())

	tags := repoTags(t, archive)
	t.Logf("archive RepoTags: %v", tags)

	found := false
	for _, tag := range tags {
		if tag == ref {
			found = true
		}
	}
	if !found {
		t.Fatalf("the archive does not carry the tag %q, a cluster importing it would end "+
			"up with an image no pod can name; got %v", ref, tags)
	}
}

// repoTags reads the tags out of a saved image archive.
func repoTags(t *testing.T, archive string) []string {
	t.Helper()

	f, err := os.Open(archive)
	if err != nil {
		t.Fatalf("%s", err)
	}
	defer f.Close()

	tr := tar.NewReader(f)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatalf("read archive: %s", err)
		}
		if hdr.Name != "manifest.json" {
			continue
		}

		var manifest []struct {
			RepoTags []string `json:"RepoTags"`
		}
		if err := json.NewDecoder(tr).Decode(&manifest); err != nil {
			t.Fatalf("decode manifest.json: %s", err)
		}
		if len(manifest) == 0 {
			return nil
		}

		return manifest[0].RepoTags
	}

	t.Fatal("no manifest.json in the archive")

	return nil
}
