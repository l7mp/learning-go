package labenv

import (
	"fmt"
	"testing"
	"time"

	"labci/pkg/kube"
	"labci/pkg/kvstore"
	"labci/pkg/splitdim"
)

// kvstoreLBManifest exposes the key-value store on a LoadBalancer so that the tests can
// read the store directly and drive its fault injection API. The labs deliberately keep
// the store cluster-internal, so this Service belongs to the CI and not to the student.
const kvstoreLBManifest = `
apiVersion: v1
kind: Service
metadata:
  name: kvstore-ci
  labels:
    app: kvstore
    labci/owned: "true"
spec:
  selector:
    app: kvstore
  type: LoadBalancer
  ports:
  - name: http-kvstore
    port: 8081
    targetPort: 8081
    protocol: TCP
`

// SetupHelloWorld deploys the lab 02 helloworld app from the student's manifests.
func SetupHelloWorld(t *testing.T) *Env {
	t.Helper()

	e := New(t)

	e.apply(e.Path("99-labs", "code", "helloworld", "deploy", "kubernetes-deployment.yaml"))
	e.apply(e.Path("99-labs", "code", "helloworld", "deploy", "kubernetes-service.yaml"))

	e.HelloIP, e.HelloPort = e.loadBalancer("helloworld")

	return e
}

// SetupSplitDimLocal deploys the lab 03 SplitDim app with its in-memory data layer.
func SetupSplitDimLocal(t *testing.T) *Env {
	t.Helper()

	e := New(t)

	e.apply(e.Path("99-labs", "code", "splitdim", "deploy", "kubernetes-local-db.yaml"))
	e.finishSplitDim()

	return e
}

// SetupSplitDimKVStore deploys the key-value store and then the SplitDim app configured
// to use it, which is the target state of labs 04 to 06. The store is started with fault
// injection enabled so that the resilience tests can make its calls fail on demand.
func SetupSplitDimKVStore(t *testing.T) *Env {
	t.Helper()

	e := New(t)
	e.deployKVStore()

	e.apply(e.Path("99-labs", "code", "splitdim", "deploy", "kubernetes-kvstore.yaml"))
	e.discoverKVStoreMode()
	e.finishSplitDim()

	return e
}

// discoverKVStoreMode learns which KVSTORE_MODE value the student uses to select the
// key-value store data layer, falling back to the name the lab text uses.
func (e *Env) discoverKVStoreMode() {
	e.KVStoreMode = "kvstore"

	d, err := e.Kube.Deployment(e.Ctx, "splitdim")
	if err != nil {
		return
	}

	ctr, err := kube.Container(&d.Spec.Template.Spec, "splitdim")
	if err != nil {
		return
	}

	if ref := kube.ConfigMapRefFor(ctr, "KVSTORE_MODE"); ref != nil {
		v, ok, err := e.Kube.ConfigMapValue(e.Ctx, ref.Name, ref.Key)
		if err == nil && ok && v != "" && v != "local" {
			e.KVStoreMode = v
			e.T.Logf("the key-value store data layer is selected by KVSTORE_MODE=%q", v)
			return
		}
	}

	for _, v := range ctr.Env {
		if v.Name == "KVSTORE_MODE" && v.Value != "" && v.Value != "local" {
			e.KVStoreMode = v.Value
			e.T.Logf("the key-value store data layer is selected by KVSTORE_MODE=%q", v.Value)
			return
		}
	}
}

// deployKVStore brings up the key-value store from the manifest shipped with the labs,
// with fault injection switched on and an extra LoadBalancer for the tests.
func (e *Env) deployKVStore() {
	e.T.Helper()

	path := e.Path("99-labs", "code", "kvstore", "deploy", "kubernetes-statefulset.yaml")

	objs, err := kube.ParseFile(path)
	if err != nil {
		e.T.Fatalf("%s", err)
	}

	for _, obj := range objs {
		if obj.GetKind() != "StatefulSet" {
			continue
		}
		if err := kube.SetContainerEnv(obj, "kvstore", "KVSTORE_FAULT_INJECTION", "1"); err != nil {
			e.T.Fatalf("enabling fault injection on the key-value store: %s", err)
		}
	}

	e.applyObjects(objs, e.relative(path))
	e.applyYAML(kvstoreLBManifest)

	ip, port := e.loadBalancer("kvstore-ci")
	e.KVStore = kvstore.New(fmt.Sprintf("http://%s:%d", ip, port))

	if err := e.KVStore.WaitReady(2 * time.Minute); err != nil {
		e.T.Fatalf("%s", err)
	}

	if !e.KVStore.FaultInjectionAvailable() {
		e.T.Fatalf("the key-value store does not serve the fault injection API even though "+
			"it was started with KVSTORE_FAULT_INJECTION set: has %s been modified?",
			"99-labs/code/kvstore/pkg/server")
	}
}

// finishSplitDim discovers the SplitDim LoadBalancer and waits for the app to answer.
func (e *Env) finishSplitDim() {
	e.T.Helper()

	e.SplitDimIP, e.SplitDimPort = e.loadBalancer("splitdim")
	e.SplitDim = splitdim.New(fmt.Sprintf("http://%s:%d", e.SplitDimIP, e.SplitDimPort))

	if err := e.SplitDim.WaitReady(2 * time.Minute); err != nil {
		e.T.Fatalf("%s", err)
	}
}

// loadBalancer returns the external address and port of a LoadBalancer Service, failing
// the test with an explanation if the Service never gets one.
func (e *Env) loadBalancer(name string) (string, int32) {
	e.T.Helper()

	ip, err := e.Kube.LoadBalancerIP(e.Ctx, name, 3*time.Minute)
	if err != nil {
		e.T.Fatalf("Service %q never got an external address: %s\n"+
			"the lab asks for a Service of type LoadBalancer", name, err)
	}

	port, err := e.Kube.ServicePort(e.Ctx, name)
	if err != nil {
		e.T.Fatalf("%s", err)
	}

	e.T.Logf("Service %q is reachable at %s:%d", name, ip, port)

	return ip, port
}

// WaitDownstreamReady blocks until SplitDim can reach the key-value store, by making it
// perform an operation that goes all the way through to the store.
//
// Probing the store directly is not enough. The labs expose it on a headless Service, so
// when it has no ready endpoints its DNS name resolves to NXDOMAIN, and CoreDNS caches
// that answer for a while. A store that has just been restarted therefore answers our own
// probes over its ClusterIP some seconds before the app can resolve it at all.
func (e *Env) WaitDownstreamReady(timeout time.Duration) {
	e.T.Helper()

	deadline := time.Now().Add(timeout)

	var last error
	for time.Now().Before(deadline) {
		if err := e.SplitDim.Reset(); err == nil {
			return
		} else {
			last = err
		}

		time.Sleep(time.Second)
	}

	e.T.Fatalf("SplitDim could not reach the key-value store within %s: %s", timeout, last)
}

// RestartSplitDim rolls the SplitDim Deployment and waits for it to serve again. Tests use
// it after changing a ConfigMap, which Kubernetes does not propagate on its own.
func (e *Env) RestartSplitDim(replicas int) {
	e.T.Helper()

	if err := e.Kube.RolloutRestart(e.Ctx, "splitdim"); err != nil {
		e.T.Fatalf("%s", err)
	}

	if err := e.Kube.WaitRollout(e.Ctx, "splitdim", readyTimeout); err != nil {
		e.T.Fatalf("SplitDim did not come back after a restart: %s", err)
	}

	if err := e.SplitDim.WaitReady(2 * time.Minute); err != nil {
		e.T.Fatalf("%s", err)
	}
}
