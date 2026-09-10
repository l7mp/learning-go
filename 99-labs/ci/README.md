# Lab tests

This module holds the automated tests for labs 02 to 06. It deploys the lab applications into a
real Kubernetes cluster from the manifests and the Dockerfiles under `99-labs/code`, and checks
that the result behaves the expected way.

## Running everything

From the repository root, with a container engine available:

``` shell
make ci-test          # test the Go exercises, then the labs
make lab-test         # test only the labs
LAB=05 make lab-test  # test only lab 05, while you are working through it
```

## Standalone test

`make lab-test` is a plain `go test` in this module. Running it directly works too, and is
the way to pass extra flags:

``` shell
cd 99-labs/ci
LABS_SOURCE=$(git rev-parse --show-toplevel) go test ./ -run TestLab05 -v -count 1 -timeout 60m
```

## Running against your own cluster

The lab tests are ordinary Go tests against whatever `KUBECONFIG` points at, so they also
run against a cluster you brought up yourself, Minikube included. Build and load the images
the manifests refer to (`localhost/helloworld`, `localhost/splitdim`, `localhost/kvstore`)
into that cluster, then:

``` shell
cd 99-labs/ci
LABS_EXTERNAL_CLUSTER=1 LABS_SOURCE=$(git rev-parse --show-toplevel) \
    go test ./ -run TestLab05 -v -count 1
```

`LABS_EXTERNAL_CLUSTER` tells the harness to stand aside: it builds no images and creates
no cluster, and the tests use whatever `KUBECONFIG` names. `LABS_SOURCE` points at the
repository root; without it the tests look for it above the working directory.

## Caveats

- The test suite expects the applications to live in `99-labs/code/{helloworld,splitdim,kvstore}`.
- The manifests should be in `deploy/kubernetes-deployment.yaml` and
  `deploy/kubernetes-service.yaml` for `helloworld`, and `deploy/kubernetes-local-db.yaml` and
  `deploy/kubernetes-kvstore.yaml` for SplitDim.
- Lab 06 reads `KVSTORE_MODE` and `KVSTORE_ADDR` from the `kvstoreMode` and `kvstoreAddr`
  entries of a ConfigMap. Naming the keys in any other way will fail the tests.
- For the tests to run on Minikube, the manifests are adjusted before deployment: `imagePullPolicy`
  is set to `IfNotPresent` for `localhost/` images.

