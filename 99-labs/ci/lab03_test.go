package ci

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"labci/pkg/labenv"
	"labci/pkg/splitdim"
)

// Lab 03 asks the student to build the SplitDim web app over an in-memory data layer and
// deploy it to Kubernetes behind a LoadBalancer.

// TestLab03 deploys SplitDim with its local data layer and runs the lab's own API suite
// against it, plus an end-to-end check of the worked example from the README.
func TestLab03(t *testing.T) {
	e := labenv.SetupSplitDimLocal(t)
	defer e.Teardown()

	t.Run("APISuite", func(t *testing.T) {
		e.RunSuite(t, "splitdim", labenv.SuiteTags)
	})

	t.Run("WorkedExample", func(t *testing.T) {
		checkWorkedExample(t, e.SplitDim)
	})
}

// checkWorkedExample replays the transfer sequence the lab READMEs walk through: a owes
// nothing, b is square, and c ends up owing a a single unit.
func checkWorkedExample(t *testing.T, c *splitdim.Client) {
	t.Helper()

	if err := c.Reset(); err != nil {
		t.Fatalf("%s", err)
	}

	for _, tr := range []splitdim.Transfer{
		{Sender: "a", Receiver: "b", Amount: 1},
		{Sender: "b", Receiver: "c", Amount: 1},
	} {
		status, _, err := c.Transfer(tr)
		if err != nil {
			t.Fatalf("%s", err)
		}
		assert.Equal(t, http.StatusOK, status, "transfer %s -> %s", tr.Sender, tr.Receiver)
	}

	accounts, err := c.Accounts()
	if err != nil {
		t.Fatalf("%s", err)
	}
	assert.Equal(t, []splitdim.Account{
		{Holder: "a", Balance: 1},
		{Holder: "b", Balance: 0},
		{Holder: "c", Balance: -1},
	}, accounts, "balances after the two transfers")

	transfers, err := c.Clear()
	if err != nil {
		t.Fatalf("%s", err)
	}
	assert.Equal(t, []splitdim.Transfer{{Sender: "c", Receiver: "a", Amount: 1}}, transfers,
		"the cheapest way to settle up is for c to pay a")

	// Clearing is hypothetical: it must not move any money.
	accounts, err = c.Accounts()
	if err != nil {
		t.Fatalf("%s", err)
	}
	assert.Equal(t, []splitdim.Account{
		{Holder: "a", Balance: 1},
		{Holder: "b", Balance: 0},
		{Holder: "c", Balance: -1},
	}, accounts, "clear should report transfers, not perform them")
}
