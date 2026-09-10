package ci

import (
	"fmt"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"labci/pkg/kube"
	"labci/pkg/kvstore"
	"labci/pkg/labenv"
	"labci/pkg/splitdim"
)

// Lab 05 asks the student to make SplitDim resilient: bound the retries against a failing
// key-value store, shut down gracefully, and expose a health check.

// giveUp is the bound on how long a transfer may take. The retry policy the lab suggests
// (4 attempts, 50ms base, 500ms cap) finishes in about a second, while an unbounded retry
// loop never does.
const giveUp = 5 * time.Second

func TestLab05(t *testing.T) {
	e := labenv.SetupSplitDimKVStore(t)
	defer e.Teardown()

	// The command the lab prescribes: the API suite, the retry budget and the health
	// endpoint, with the mode tag insisting the app really talks to the key-value store.
	t.Run("Suite", func(t *testing.T) {
		e.RunSuite(t, "splitdim", labenv.ResilienceSuiteTags)
	})

	t.Run("HealthProbe", func(t *testing.T) { lab05HealthProbe(t, e) })
	t.Run("GivesUpOnADeadStore", func(t *testing.T) { lab05GivesUp(t, e) })
	t.Run("SurvivesAFlakyStore", func(t *testing.T) { lab05FlakyStore(t, e) })
	t.Run("GracefulShutdown", func(t *testing.T) { lab05GracefulShutdown(t, e) })
}

// lab05HealthProbe checks the half of the health check task that only Kubernetes can
// answer: the endpoint itself is covered by the suite, but the Deployment also has to tell
// Kubernetes to probe it, and the probe has to be pointed somewhere that actually works.
func lab05HealthProbe(t *testing.T, e *labenv.Env) {
	d, err := e.Kube.Deployment(e.Ctx, "splitdim")
	if err != nil {
		t.Fatalf("%s", err)
	}

	ctr, err := kube.Container(&d.Spec.Template.Spec, "splitdim")
	if err != nil {
		t.Fatalf("%s", err)
	}

	probe := ctr.LivenessProbe
	if probe == nil || probe.HTTPGet == nil {
		t.Fatalf("the splitdim container declares no HTTP livenessProbe: Kubernetes has " +
			"no way to tell whether the app is alive")
	}
	assert.Equal(t, "/healthz", probe.HTTPGet.Path, "the liveness probe should hit /healthz")

	// A misconfigured probe restarts the pod every few seconds: with the periodSeconds
	// the lab suggests, the first restart lands within about ten. Watch for a while and
	// stop at the first one rather than sleeping through the whole window.
	before, err := e.Kube.RestartCount(e.Ctx, "app=splitdim")
	if err != nil {
		t.Fatalf("%s", err)
	}

	deadline := time.Now().Add(20 * time.Second)
	for time.Now().Before(deadline) {
		after, err := e.Kube.RestartCount(e.Ctx, "app=splitdim")
		if err != nil {
			t.Fatalf("%s", err)
		}
		if !assert.Equal(t, before, after,
			"a SplitDim pod restarted while idle: the liveness probe is misconfigured, "+
				"it is failing against a healthy app") {
			return
		}

		time.Sleep(2 * time.Second)
	}
}

// lab05GivesUp runs the lab's own "the store is not there" check, with Kubernetes taking
// the store away. That run refuses to start unless the store really is unreachable, so it
// cannot pass by accident.
func lab05GivesUp(t *testing.T, e *labenv.Env) {
	if err := e.Kube.ScaleStatefulSet(e.Ctx, "kvstore", 0); err != nil {
		t.Fatalf("%s", err)
	}
	defer func() {
		if err := e.Kube.ScaleStatefulSet(e.Ctx, "kvstore", 1); err != nil {
			t.Fatalf("%s", err)
		}
		if err := e.KVStore.WaitReady(2 * time.Minute); err != nil {
			t.Fatalf("%s", err)
		}
		// The store answering us is not the same as the app being able to resolve it.
		e.WaitDownstreamReady(2 * time.Minute)
	}()

	if err := e.Kube.WaitPodsReady(e.Ctx, "app=kvstore", 0, 2*time.Minute); err != nil {
		t.Fatalf("stopping the key-value store: %s", err)
	}

	e.RunSuite(t, "splitdim", labenv.FailingKVStoreTags)
}

// lab05FlakyStore replaces the Istio fault injection from the README: a third of the puts
// fail, and the app has to keep both working and consistent.
func lab05FlakyStore(t *testing.T, e *labenv.Env) {
	if err := e.KVStore.ClearFaults(); err != nil {
		t.Fatalf("%s", err)
	}
	defer e.KVStore.ClearFaults()

	if err := e.SplitDim.Reset(); err != nil {
		t.Fatalf("%s", err)
	}

	if err := e.KVStore.Fault(kvstore.FaultConfig{
		Path: "/api/put", AbortPercent: 33,
	}); err != nil {
		t.Fatalf("%s", err)
	}

	for i := 0; i < 12; i++ {
		status, _, err := e.SplitDim.WithTimeout(giveUp).Transfer(splitdim.Transfer{
			Sender: "a", Receiver: "b", Amount: 1,
		})
		if err != nil {
			t.Fatalf("transfer %d never returned: %s", i, err)
		}
		if !assert.Equal(t, http.StatusOK, status,
			"transfer %d failed even though only a third of the puts are being rejected", i) {
			break
		}
	}

	if err := e.KVStore.ClearFaults(); err != nil {
		t.Fatalf("%s", err)
	}

	total, err := e.SplitDim.TotalBalance()
	if err != nil {
		t.Fatalf("%s", err)
	}
	assert.Equal(t, 0, total,
		"the balances no longer add up to zero: a transfer was left applied only halfway")
}

// lab05GracefulShutdown checks that a terminating pod finishes the requests it has already
// accepted. The key-value store is slowed down so that requests are demonstrably in
// flight when the pod is deleted; without server.Shutdown the app dies on SIGTERM at once
// and those clients see a connection reset.
func lab05GracefulShutdown(t *testing.T, e *labenv.Env) {
	d, err := e.Kube.Deployment(e.Ctx, "splitdim")
	if err != nil {
		t.Fatalf("%s", err)
	}

	grace := int64(0)
	if d.Spec.Template.Spec.TerminationGracePeriodSeconds != nil {
		grace = *d.Spec.Template.Spec.TerminationGracePeriodSeconds
	}
	assert.GreaterOrEqual(t, grace, int64(600),
		"the pod spec does not ask Kubernetes to wait for a graceful shutdown: the lab "+
			"asks for terminationGracePeriodSeconds: 600")

	if err := e.KVStore.ClearFaults(); err != nil {
		t.Fatalf("%s", err)
	}
	defer e.KVStore.ClearFaults()

	if err := e.SplitDim.Reset(); err != nil {
		t.Fatalf("%s", err)
	}

	pods, err := e.Kube.PodNames(e.Ctx, "app=splitdim")
	if err != nil {
		t.Fatalf("%s", err)
	}
	if len(pods) != 1 {
		t.Fatalf("expected exactly one SplitDim pod, found %d", len(pods))
	}

	// Every put now takes two seconds, so the transfers below are unambiguously still
	// running when the pod is told to go away.
	if err := e.KVStore.Fault(kvstore.FaultConfig{
		Path: "/api/put", Delay: "2s",
	}); err != nil {
		t.Fatalf("%s", err)
	}

	const inFlight = 8

	var wg sync.WaitGroup
	codes := make([]int, inFlight)
	errs := make([]error, inFlight)

	// Each request uses its own pair of accounts: this test is about draining in-flight
	// work, not about surviving write conflicts, and the slowed-down store would make
	// contending transfers exhaust their retry budget for unrelated reasons.
	for i := 0; i < inFlight; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			codes[i], _, errs[i] = e.SplitDim.WithTimeout(2 * time.Minute).
				Transfer(splitdim.Transfer{
					Sender:   fmt.Sprintf("payer%d", i),
					Receiver: fmt.Sprintf("payee%d", i),
					Amount:   1,
				})
		}(i)
	}

	// Let the requests reach the app and block on the slow key-value store.
	time.Sleep(2 * time.Second)

	if err := e.Kube.DeletePod(e.Ctx, pods[0], grace); err != nil {
		t.Fatalf("%s", err)
	}

	wg.Wait()

	for i := range errs {
		if errs[i] != nil {
			t.Errorf("in-flight request %d did not complete: %s\n"+
				"the pod stopped serving as soon as it received SIGTERM instead of "+
				"draining its outstanding requests", i, errs[i])
			continue
		}
		assert.Equal(t, http.StatusOK, codes[i],
			"in-flight request %d failed while the pod was shutting down", i)
	}

	// Whether the pod lingered for a particular number of seconds is not the point and
	// timing it only makes the check flaky. The requests above completing is the property
	// that matters: a server that ignores SIGTERM drops them.

	if err := e.Kube.WaitRollout(e.Ctx, "splitdim", 3*time.Minute); err != nil {
		t.Fatalf("SplitDim did not come back after the pod was replaced: %s", err)
	}
	if err := e.SplitDim.WaitReady(2 * time.Minute); err != nil {
		t.Fatalf("%s", err)
	}

	if err := e.KVStore.ClearFaults(); err != nil {
		t.Fatalf("%s", err)
	}

	total, err := e.SplitDim.TotalBalance()
	if err != nil {
		t.Fatalf("%s", err)
	}
	assert.Equal(t, 0, total,
		"the balances no longer add up to zero after the shutdown: a transfer was "+
			"interrupted halfway through")
}
