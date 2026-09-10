package ci

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"labci/pkg/kube"
	"labci/pkg/labenv"
	"labci/pkg/splitdim"
)

// Lab 06 asks the student to make the data layer selectable: through command line flags,
// and through a ConfigMap mapped onto the KVSTORE_MODE and KVSTORE_ADDR environment
// variables.
//
// Both checks below work the same way: flip the ConfigMap, register a transfer, and look at
// whether the balances landed in the key-value store. The SplitDim API behaves identically
// whichever data layer is live, so this is the only thing that actually distinguishes them,
// and a student who ignores the ConfigMap entirely passes anything else.
//
// The command line flags are the lab's other half, and they are checked where the lab puts
// them, on localhost: see the "Command line parameters" section of its README.

const (
	// configMapKeys are the ConfigMap entries the lab prescribes.
	modeKey = "kvstoreMode"
	addrKey = "kvstoreAddr"

	// kvstoreAddr is where the key-value store answers inside the cluster.
	kvstoreAddr = "kvstore.default:8081"
)

func TestLab06(t *testing.T) {
	e := labenv.SetupSplitDimKVStore(t)
	defer e.Teardown()

	ref := lab06ConfigMapWiring(t, e)

	t.Run("ConfigMapSelectsTheLocalDataLayer", func(t *testing.T) {
		applyConfig(t, e, ref.Name, map[string]string{modeKey: "local"})
		e.RestartSplitDim(1)

		assert.Empty(t, transferAndReadStore(t, e),
			"the key-value store was written to even though the ConfigMap asked for the "+
				"local data layer: KVSTORE_MODE is not reaching the app")
	})

	t.Run("ConfigMapSelectsTheKeyValueStore", func(t *testing.T) {
		applyConfig(t, e, ref.Name, map[string]string{
			modeKey: e.KVStoreMode,
			addrKey: kvstoreAddr,
		})
		e.RestartSplitDim(1)

		balances := transferAndReadStore(t, e)
		assert.Equal(t, 7, balances["alice"],
			"the balances did not reach the key-value store even though the ConfigMap "+
				"selected the %q data layer", e.KVStoreMode)
		assert.Equal(t, -7, balances["bob"], "bob's balance did not reach the store")

		e.RunSuite(t, "splitdim", labenv.KVStoreSuiteTags)
	})

}

// lab06ConfigMapWiring checks that the Deployment sources KVSTORE_MODE and KVSTORE_ADDR
// from a ConfigMap, and returns how it is wired.
func lab06ConfigMapWiring(t *testing.T, e *labenv.Env) *kube.ConfigMapRef {
	t.Helper()

	d, err := e.Kube.Deployment(e.Ctx, "splitdim")
	if err != nil {
		t.Fatalf("%s", err)
	}

	ctr, err := kube.Container(&d.Spec.Template.Spec, "splitdim")
	if err != nil {
		t.Fatalf("%s", err)
	}

	mode := kube.ConfigMapRefFor(ctr, "KVSTORE_MODE")
	if mode == nil {
		t.Fatalf("the splitdim container does not take KVSTORE_MODE from a ConfigMap: the " +
			"lab asks you to map the splitdim-config ConfigMap onto the environment " +
			"with a configMapKeyRef")
	}

	addr := kube.ConfigMapRefFor(ctr, "KVSTORE_ADDR")
	if addr == nil {
		t.Fatalf("the splitdim container does not take KVSTORE_ADDR from a ConfigMap")
	}

	assert.Equal(t, modeKey, mode.Key,
		"the lab prescribes the ConfigMap entry %q for the data layer mode", modeKey)
	assert.Equal(t, addrKey, addr.Key,
		"the lab prescribes the ConfigMap entry %q for the key-value store address", addrKey)
	assert.Equal(t, mode.Name, addr.Name,
		"KVSTORE_MODE and KVSTORE_ADDR should come from the same ConfigMap")

	t.Logf("splitdim reads its config from ConfigMap %q (%s -> %s, %s -> %s)",
		mode.Name, mode.Key, mode.EnvVar, addr.Key, addr.EnvVar)

	return mode
}

// applyConfig replaces the ConfigMap the student's Deployment reads.
func applyConfig(t *testing.T, e *labenv.Env, name string, data map[string]string) {
	t.Helper()

	manifest := fmt.Sprintf("apiVersion: v1\nkind: ConfigMap\nmetadata:\n  name: %s\ndata:\n", name)
	for k, v := range data {
		manifest += fmt.Sprintf("  %s: %q\n", k, v)
	}

	if _, err := e.Kube.ApplyYAML(e.Ctx, manifest); err != nil {
		t.Fatalf("applying ConfigMap %q: %s", name, err)
	}
}

// transferAndReadStore empties both stores, registers one transfer through SplitDim, and
// returns what ended up in the key-value store.
func transferAndReadStore(t *testing.T, e *labenv.Env) map[string]int {
	t.Helper()

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
	if status != http.StatusOK {
		t.Fatalf("transfer returned HTTP %d", status)
	}

	return storeBalances(t, e.KVStore.BaseURL)
}
