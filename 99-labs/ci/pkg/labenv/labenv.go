// Package labenv sets up and tears down the Kubernetes state a lab test needs: it builds
// the student's manifests into a running deployment and hands the test the endpoints to
// talk to. Bringing up the cluster and the images the manifests refer to belongs to the
// harness in the parent package, not here.
package labenv

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"labci/pkg/kube"
	"labci/pkg/kvstore"
	"labci/pkg/splitdim"
)

const (
	// SuiteTags covers the whole SplitDim API. With no mode tag alongside it the run
	// expects the in-memory data layer, which is what lab 03 asks for.
	SuiteTags = "httphandler,api,localconstructor,reset,transfer,accounts,clear"

	// KVStoreSuiteTags is the same suite plus the assertion that the app really keeps its
	// accounts in the key-value store. This is the command labs 04 and 06 prescribe.
	KVStoreSuiteTags = "kvstoremode," + SuiteTags

	// ResilienceSuiteTags adds the resilience checks, as lab 05 prescribes.
	ResilienceSuiteTags = "kvstoremode,retry,healthz," + SuiteTags

	// FailingKVStoreTags is the run for a key-value store that is deliberately not there:
	// it refuses to start if the store answers, so it cannot pass by accident.
	FailingKVStoreTags = "failingkvstoremode"

	// readyTimeout is how long a workload gets to come up before we call it broken.
	readyTimeout = 4 * time.Minute

	// goneTimeout is how long a teardown waits for objects to disappear.
	goneTimeout = 2 * time.Minute
)

// Env is the live state of one lab under test.
type Env struct {
	T    *testing.T
	Ctx  context.Context
	Kube *kube.Client

	// Root is the repository root.
	Root string

	// SplitDimIP and SplitDimPort locate the SplitDim LoadBalancer, and SplitDim is a
	// client for it. They are empty for labs that do not deploy SplitDim.
	SplitDimIP   string
	SplitDimPort int32
	SplitDim     *splitdim.Client

	// KVStore is a client for the key-value store, reachable through a CI-only
	// LoadBalancer Service. It is nil for labs that do not deploy the store.
	KVStore *kvstore.Client

	// HelloIP and HelloPort locate the helloworld LoadBalancer.
	HelloIP   string
	HelloPort int32

	// KVStoreMode is the KVSTORE_MODE value that selects the key-value store data layer.
	// The labs let students call it kvstore, resilientkvstore or transactionalkvstore, so
	// the CI uses whichever one their own manifest asks for.
	KVStoreMode string

	applied []*kube.Object
}

// Root returns the repository root, failing the test when it cannot be found.
func Root(t *testing.T) string {
	t.Helper()

	root, err := FindRoot()
	if err != nil {
		t.Fatalf("%s", err)
	}

	return root
}

// FindRoot returns the repository root: LABS_SOURCE when the caller sets it, otherwise the
// nearest ancestor directory that contains 99-labs.
func FindRoot() (string, error) {
	if root := os.Getenv("LABS_SOURCE"); root != "" {
		return root, nil
	}

	dir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("cannot determine the working directory: %w", err)
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, "99-labs")); err == nil {
			return dir, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("cannot find the repository root above %q, set LABS_SOURCE", dir)
		}
		dir = parent
	}
}

// New connects to the cluster and returns an environment with nothing deployed yet.
func New(t *testing.T) *Env {
	t.Helper()

	c, err := kube.NewClient()
	if err != nil {
		t.Fatalf("cannot build a Kubernetes client: %s", err)
	}

	ctx := context.Background()
	if err := c.Ping(ctx); err != nil {
		t.Fatalf("%s", err)
	}

	return &Env{T: t, Ctx: ctx, Kube: c, Root: Root(t)}
}

// Path resolves a path relative to the repository root.
func (e *Env) Path(parts ...string) string {
	return filepath.Join(append([]string{e.Root}, parts...)...)
}

// apply applies a manifest file, remembering the objects so that Teardown can remove
// exactly what was created.
func (e *Env) apply(path string) []*kube.Object {
	e.T.Helper()

	if _, err := os.Stat(path); err != nil {
		e.T.Fatalf("missing manifest %s: this lab asks you to write it",
			e.relative(path))
	}

	objs, err := kube.ParseFile(path)
	if err != nil {
		e.T.Fatalf("%s", err)
	}

	e.normalize(objs)

	if err := e.Kube.Apply(e.Ctx, objs); err != nil {
		e.T.Fatalf("applying %s: %s", e.relative(path), err)
	}
	e.applied = append(e.applied, objs...)

	if err := e.Kube.WaitReady(e.Ctx, objs, readyTimeout); err != nil {
		e.T.Fatalf("after applying %s: %s", e.relative(path), err)
	}

	return objs
}

// applyYAML applies an inline manifest owned by the CI itself.
func (e *Env) applyYAML(manifest string) []*kube.Object {
	e.T.Helper()

	objs, err := e.Kube.ApplyYAML(e.Ctx, manifest)
	if err != nil {
		e.T.Fatalf("applying CI manifest: %s", err)
	}
	e.applied = append(e.applied, objs...)

	if err := e.Kube.WaitReady(e.Ctx, objs, readyTimeout); err != nil {
		e.T.Fatalf("CI manifest never became ready: %s", err)
	}

	return objs
}

// applyObjects applies objects the CI has already parsed and modified.
func (e *Env) applyObjects(objs []*kube.Object, what string) {
	e.T.Helper()

	e.normalize(objs)

	if err := e.Kube.Apply(e.Ctx, objs); err != nil {
		e.T.Fatalf("applying %s: %s", what, err)
	}
	e.applied = append(e.applied, objs...)

	if err := e.Kube.WaitReady(e.Ctx, objs, readyTimeout); err != nil {
		e.T.Fatalf("after applying %s: %s", what, err)
	}
}

// normalize adjusts the manifests for the differences between the CI cluster and the
// Minikube setup the labs are written for, and says in the log what it touched.
func (e *Env) normalize(objs []*kube.Object) {
	e.T.Helper()

	for _, obj := range objs {
		changed, err := kube.NormalizeLocalImages(obj)
		if err != nil {
			e.T.Fatalf("%s", err)
		}
		for _, c := range changed {
			e.T.Logf("note: %s", c)
		}
	}
}

func (e *Env) relative(path string) string {
	rel, err := filepath.Rel(e.Root, path)
	if err != nil {
		return path
	}
	return rel
}

// Teardown removes everything this environment created. When the test has failed it first
// dumps the cluster state, so a CI log explains itself.
func (e *Env) Teardown() {
	if e.T.Failed() {
		for _, obj := range e.applied {
			switch obj.GetKind() {
			case "Deployment", "StatefulSet":
				e.T.Logf("%s", e.Kube.Diagnose(e.Ctx, obj))
			}
		}
	}

	if err := e.Kube.Delete(e.Ctx, e.applied); err != nil {
		e.T.Logf("teardown: %s", err)
	}

	if err := e.Kube.WaitGone(e.Ctx, e.applied, goneTimeout); err != nil {
		e.T.Logf("teardown: %s", err)
	}

	// A deleted Deployment or StatefulSet is reported gone before its pods have finished
	// terminating, and a deleted LoadBalancer Service keeps its address until the
	// provider tears the forwarding down. The labs share one cluster, so the next one
	// must not inherit either.
	if err := e.Kube.WaitWorkloadsDrained(e.Ctx, e.applied, goneTimeout); err != nil {
		e.T.Logf("teardown: %s", err)
	}

	// StatefulSet volumes outlive their owner, and the next lab must not inherit the
	// previous lab's transaction log.
	if err := e.Kube.DeletePVCs(e.Ctx, "app=kvstore"); err != nil {
		e.T.Logf("teardown: %s", err)
	}

	e.applied = nil
}

// RunSuite runs one of the students' own build-tagged test suites against the deployment,
// which is exactly the command the lab README asks them to run by hand.
//
// It reports on the t it is given rather than on the environment's own, so that a suite run
// from inside a subtest fails that subtest. Reporting on the parent instead would print the
// subtest as passing and fail the whole lab with nothing to point at.
func (e *Env) RunSuite(t *testing.T, module, tags string) {
	t.Helper()
	e.goTest(t, module, "suite "+tags,
		[]string{"test", "./...", "--tags=" + tags, "-v", "-count", "1"})
}

// RunPackage runs the tests of a single package inside a lab module, for the checks a lab
// prescribes on its own, like the key-value store client library in lab 04.
func (e *Env) RunPackage(t *testing.T, module, pkg string) {
	t.Helper()
	e.goTest(t, module, pkg, []string{"test", pkg, "-v", "-count", "1"})
}

func (e *Env) goTest(t *testing.T, module, what string, args []string) {
	t.Helper()

	dir := e.Path("99-labs", "code", module)

	cmd := exec.CommandContext(e.Ctx, "go", args...)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), e.suiteEnv()...)

	out, err := cmd.CombinedOutput()

	// Go reports "no test files" as success but it would mean the run silently did
	// nothing, which must not read as a pass.
	if !strings.Contains(string(out), "PASS") && !strings.Contains(string(out), "ok  ") {
		t.Errorf("%s in %s produced no passing tests", what, module)
	}

	if err != nil {
		t.Errorf("go test %s in %s failed: %s\n%s", what, module, err, out)
		return
	}

	t.Logf("go test %s in %s passed", what, module)
}

// suiteEnv points the students' tests at the deployment under test.
func (e *Env) suiteEnv() []string {
	env := []string{}

	if e.SplitDimIP != "" {
		env = append(env,
			"EXTERNAL_IP="+e.SplitDimIP,
			fmt.Sprintf("EXTERNAL_PORT=%d", e.SplitDimPort))
	}
	if e.HelloIP != "" {
		env = append(env,
			"EXTERNAL_IP="+e.HelloIP,
			fmt.Sprintf("EXTERNAL_PORT=%d", e.HelloPort))
	}
	if e.KVStore != nil {
		env = append(env, "KVSTORE_URL="+e.KVStore.BaseURL)
	}

	return env
}
