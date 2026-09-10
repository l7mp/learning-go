package ci

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"labci/pkg/labenv"
	"labci/pkg/splitdim"
)

// Lab 04 asks the student to move the account database out of the app and into the
// key-value store, so that SplitDim becomes stateless and can be scaled and restarted.

// TestLab04 checks the externalised data layer: that the balances really live in the
// key-value store, that several replicas agree, and that the data survives a restart.
func TestLab04(t *testing.T) {
	e := labenv.SetupSplitDimKVStore(t)
	defer e.Teardown()

	// The client library the lab has you write, checked by its own test, which needs no
	// cluster and starts a key-value store of its own.
	t.Run("ClientLibrary", func(t *testing.T) {
		e.RunPackage(t, "kvstore", "./pkg/client/...")
	})

	t.Run("APISuite", func(t *testing.T) {
		e.RunSuite(t, "splitdim", labenv.KVStoreSuiteTags)
	})

	// The suite above passes just as happily against the in-memory data layer. This is
	// the check that the state was actually externalised.
	t.Run("StateLivesInTheKeyValueStore", func(t *testing.T) {
		clearStore(t, e.KVStore.BaseURL)

		if err := e.SplitDim.Reset(); err != nil {
			t.Fatalf("%s", err)
		}

		status, _, err := e.SplitDim.Transfer(splitdim.Transfer{
			Sender: "alice", Receiver: "bob", Amount: 7,
		})
		if err != nil {
			t.Fatalf("%s", err)
		}
		assert.Equal(t, http.StatusOK, status, "transfer")

		balances := storeBalances(t, e.KVStore.BaseURL)

		assert.Equal(t, 7, balances["alice"],
			"the key-value store does not hold alice's balance: SplitDim is still keeping "+
				"the accounts in memory")
		assert.Equal(t, -7, balances["bob"], "the key-value store does not hold bob's balance")
	})

	// Replicas backed by their own in-memory map disagree; replicas backed by a shared
	// store do not.
	t.Run("ReplicasAgree", func(t *testing.T) {
		scaleSplitDim(t, e, 3)
		defer scaleSplitDim(t, e, 1)

		if err := e.SplitDim.Reset(); err != nil {
			t.Fatalf("%s", err)
		}

		status, _, err := e.SplitDim.Transfer(splitdim.Transfer{
			Sender: "carol", Receiver: "dave", Amount: 3,
		})
		if err != nil {
			t.Fatalf("%s", err)
		}
		assert.Equal(t, http.StatusOK, status, "transfer")

		want := []splitdim.Account{
			{Holder: "carol", Balance: 3},
			{Holder: "dave", Balance: -3},
		}

		// The Service spreads these over all three pods, so a disagreement shows up as a
		// reply that does not match.
		for i := 0; i < 30; i++ {
			got, err := e.SplitDim.Accounts()
			if err != nil {
				t.Fatalf("%s", err)
			}
			if !assert.Equal(t, want, got,
				"request %d got a different answer: the replicas do not share their state", i) {
				break
			}
		}
	})

	// A handful of transfers issued at once against a scaled-out deployment. Each uses
	// its own pair of accounts, so this only asks the simple question the lab is about:
	// do several replicas sharing one store apply every transfer exactly once? Making
	// concurrent transfers touch the *same* account is a different subject, transactions,
	// which the lab treats as optional and which we deliberately do not test here.
	t.Run("TransfersAcrossReplicas", func(t *testing.T) {
		scaleSplitDim(t, e, 3)
		defer scaleSplitDim(t, e, 1)

		if err := e.SplitDim.Reset(); err != nil {
			t.Fatalf("%s", err)
		}

		const rounds = 10

		var wg sync.WaitGroup
		codes := make([]int, rounds)
		errs := make([]error, rounds)

		for i := 0; i < rounds; i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				codes[i], _, errs[i] = e.SplitDim.Transfer(splitdim.Transfer{
					Sender:   fmt.Sprintf("payer%d", i),
					Receiver: fmt.Sprintf("payee%d", i),
					Amount:   1,
				})
			}(i)
		}
		wg.Wait()

		for i := range errs {
			if errs[i] != nil {
				t.Fatalf("transfer %d never returned: %s", i, errs[i])
			}
			assert.Equal(t, http.StatusOK, codes[i], "transfer %d", i)
		}

		want := []splitdim.Account{}
		for i := 0; i < rounds; i++ {
			want = append(want,
				splitdim.Account{Holder: fmt.Sprintf("payee%d", i), Balance: -1},
				splitdim.Account{Holder: fmt.Sprintf("payer%d", i), Balance: 1})
		}
		sort.Slice(want, func(i, j int) bool { return want[i].Holder < want[j].Holder })

		got, err := e.SplitDim.Accounts()
		if err != nil {
			t.Fatalf("%s", err)
		}
		assert.Equal(t, want, got, "every transfer should be applied exactly once")
	})

	// The manual "Persistence" section of the lab, automated.
	t.Run("SurvivesARestart", func(t *testing.T) {
		if err := e.SplitDim.Reset(); err != nil {
			t.Fatalf("%s", err)
		}

		for _, tr := range []splitdim.Transfer{
			{Sender: "a", Receiver: "b", Amount: 1},
			{Sender: "b", Receiver: "c", Amount: 1},
		} {
			if _, _, err := e.SplitDim.Transfer(tr); err != nil {
				t.Fatalf("%s", err)
			}
		}

		before, err := e.SplitDim.Clear()
		if err != nil {
			t.Fatalf("%s", err)
		}

		if err := e.Kube.RolloutRestartStatefulSet(e.Ctx, "kvstore"); err != nil {
			t.Fatalf("%s", err)
		}
		e.RestartSplitDim(1)

		if err := e.KVStore.WaitReady(2 * time.Minute); err != nil {
			t.Fatalf("%s", err)
		}

		after, err := e.SplitDim.Clear()
		if err != nil {
			t.Fatalf("after restarting both components: %s", err)
		}

		assert.Equal(t, before, after,
			"the ledger changed across a restart of SplitDim and the key-value store: the "+
				"state was not persisted")
	})
}

// scaleSplitDim resizes the SplitDim Deployment and waits for the rollout to settle.
func scaleSplitDim(t *testing.T, e *labenv.Env, replicas int) {
	t.Helper()

	if err := e.Kube.Scale(e.Ctx, "splitdim", int32(replicas)); err != nil {
		t.Fatalf("%s", err)
	}
	if err := e.Kube.WaitRollout(e.Ctx, "splitdim", 3*time.Minute); err != nil {
		t.Fatalf("scaling SplitDim to %d replicas: %s", replicas, err)
	}

	// The rollout being complete is not the same as the Service forwarding to it: the
	// load balancer is reprogrammed a moment later, and a request sent in between is
	// refused.
	if err := e.SplitDim.WaitReady(2 * time.Minute); err != nil {
		t.Fatalf("%s", err)
	}
}

// The two calls below read and clear the key-value store over plain HTTP rather than
// through a client library. Lab 04 asks students to write exactly such a library, so a
// finished one in this repository would hand them the answer; these follow the style the
// slides use for talking to the store.

// storeBalances returns the numeric entries the key-value store holds.
func storeBalances(t *testing.T, baseURL string) map[string]int {
	t.Helper()

	r, err := http.Get(baseURL + "/api/list")
	if err != nil {
		t.Fatalf("reading the key-value store: %s", err)
	}
	defer r.Body.Close()

	if r.StatusCode != http.StatusOK {
		t.Fatalf("reading the key-value store: HTTP status %d", r.StatusCode)
	}

	entries := []struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}{}
	if err := json.NewDecoder(r.Body).Decode(&entries); err != nil {
		t.Fatalf("decoding the key-value store contents: %s", err)
	}

	// The store is generic and only SplitDim gives its values a meaning, so anything that
	// is not a number is not an account balance.
	balances := map[string]int{}
	for _, e := range entries {
		if n, err := strconv.Atoi(e.Value); err == nil {
			balances[e.Key] = n
		}
	}

	return balances
}

// clearStore empties the key-value store.
func clearStore(t *testing.T, baseURL string) {
	t.Helper()

	r, err := http.Get(baseURL + "/api/reset")
	if err != nil {
		t.Fatalf("clearing the key-value store: %s", err)
	}
	defer r.Body.Close()

	if r.StatusCode != http.StatusOK {
		t.Fatalf("clearing the key-value store: HTTP status %d", r.StatusCode)
	}
}
