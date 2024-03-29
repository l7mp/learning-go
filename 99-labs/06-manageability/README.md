# Cloud-native patterns: Manageability

In the course of this lab we will continue our journey to develop our toy web app, SplitDim, into a real cloud native microservice application that is stateless, scalable, manageable and resilient. After we achieved immutability and resilience, this time we concentrate on manageability.

![SplitDim logo, generated by logoai.com.](/99-labs/fig/splitdim-logo.png)

The below tasks guide you in making the web app (a little bit more) manageable. The tasks are followed by tests; pass each to complete the lab.

## Table of Contents

1. [Preliminaries](#preliminaries)
2. [Command line parameters](#command-line-parameters)
3. [Configuration files](#configuration-files)

## Preliminaries

Recall, SplitDim is a web app that lets groups of people keep track of who owns who. The web service implements the following API endpoints:
- `GET /`: serves a static HTML/JS artifact that you can use to interact with the app from a browser,
- `POST /api/transfer`: register a transfer between two users of a given amount,
- `GET /api/accounts`: return the list of current balances for each registered user,
- `GET /api/clear`: return the list of transfers that would allow users to clear their debts, and
- `GET /api/reset`: reset all balances to zero.

The resource state of the app, the accounts database, is hidden behind a `DataLayer` interface:
```go
// DataLayer is an API for manipulating a balance sheet.
type DataLayer interface {
    // Transfer will process a transfer.
    Transfer(t Transfer) error
    // AccountList returns the current acount of each user.
    AccountList() ([]Account, error)
    // Clear returns the list of transfers to clear all debts.
    Clear() ([]Transfer, error)
    // Reset sets all balances to zero.
    Reset() error
}
```

This allows us to support a number of different data layer implementations, each serving a different purpose:
- `local`: The local data layer that keeps the account database in memory. This is not enough to support scaling, but it is still a perfect way to test the app (say, for frontend development).
- `kvstore`: This is supposed to be the default key-value store backend.
- `resilientkvstore` (optional): If you created a separate data layer implementation for implementing the resilience patterns in the previous lab, then that may also be selectable.
- `transactionalkvstore` (optional): Likewise, if you implemented the transactional implementation for `Transfer` then that may also live in a separate data layer package.

Recall, we use the usual two environment variables, `KVSTORE_MODE` and `KVSTORE_ADDR`, to choose the key-value store datalayer implementation (say, `local`, `kvstore`, `resilientkvstore`, etc) and (optionally) specify the key-value store address on startup. As we are going to see in a minute that is the **right way** to do it in Kubernetes, but just for the sake of exercise let's add another way too: choosing the data layer via a command line argument. For instance, running the app with `splitdim -mode local` would choose the local mode while `splitdim -mode kvstore -addr localhost:8081` would use the key-value store reachable at `localhost:8081`.

Here is a sequence of steps that you can follow to achieve that:
- use the standard [`flag` package](https://pkg.go.dev/flag) for parsing command line arguments;
- add a string type command line flag called `mode` for specifying the data-layer mode, with the default being the value given in the environment variable `KVSTORE_MODE` and a fallback to the setting `local`;
- add a string type command line flag called `addr` for specifying the key-value store address, which is relevant only if the key-value store data layer is chosen. Let the default be the value given in the environment variable `KVSTORE_ADDR`, and fall back to `localhost:8081` if neither the environment variable nor the command line flag is given by the user.

Once everything is ready, you can play a bit with a local build to see if everything went fine.

> ✅ **Check**
> 
> Test your solution through the following steps:
> - start the key-value store: `cd 99-labs/code/kvstore && go run kvstore.go`;
> - build an executable from the `splitdim` app with `cd 99-labs/code/splitdim && go build -o splitdim main.go`;
> - start the app with the local data layer in the background: `./splitdim -mode local&`;
> - set the reachability info for `splitdim`: `export EXTERNAL_IP=localhost; export EXTERNAL_PORT=8080`;
> - run the tests:
>   ```go
>   go test ./... --tags=httphandler,api,localconstructor,reset,transfer,accounts,clear -v -count 1
>   ```
> - restart the app with the key-value store mode: `killall splitdim; ./splitdim -mode kvstore -addr localhost:8081&`;
> - rerun the tests:
>   ```go
>   go test ./... --tags=httphandler,api,localconstructor,reset,transfer,accounts,clear -v -count 1
>   ```
> - stop the app: `killall splitdim`.
> If all goes well, you should see all tests to PASS.

## Configuration files

Managing our application through environment variables and command line arguments is nice, but it gets complex after a certain point: what if we have dozens of important parameters, do we really want to configure each via a separate command line argument or environment variable? The solution is a configuration file of course, which can hold the entire config in a single human-readable format. We recommend the JSON or the YAML format for storing the config file: we already know how to automatically marshal/unmarshal them to/from Go structs.

Below we will implement something functionally equivalent with config files but, perhaps somewhat surprisingly, we will not have to deal with any files at all! This magic is made possible by Kubernetes ConfigMaps. In particular, we want to be able to collect all our relevant config parameters in a single Kubernetes resource, a ConfigMap. For instance, the below would mean to start `splitdim` in Kubernetes with the `local` data layer:
```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: splitdim-config
data:
  kvstoreMode: "local"
```

In contrast, the below would choose the `kvstore` data layer with the key-value store available at `kvstore.default:8081`:
```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: splitdim-config
data:
  kvstoreMode: "kvstore"
  kvstoreAddr: "kvstore.default:8081"
```

Our job is now to map the entries of the ConfigMap to the corresponding environment variables of the `splitdim` pod. 

In general, the entry called `my-configmap-key` in the `my-configmap` ConfigMap can be mapped into the environment variable called `MY_ENV` when starting a container called `my-container` as follows:
```yaml
...
spec:
  containers:
  - name: my-container
    ...
    env:
    - name: MY_ENV
      valueFrom:
        configMapKeyRef:
          name: my-configmap
          key: my-configmap-key
```

An additional `optional: true` setting makes sure that Kubernetes will not complain when the ConfigMap does not provide the requested entry.

> [!NOTE]
>
> Mapping ConfigMaps to command line arguments (as opposed to environment variables) is not so trivial. This is why we prefer environment variables over command line flags for managing the startup parameters of cloud native apps.

> [!NOTE]
>
> You can also [map the entire ConfigMap data as a single file](https://kubernetes.io/docs/tasks/configure-pod-container/configure-pod-configmap/#add-configmap-data-to-a-volume) into the filesystem of the pod and use the standard filesystem operations to read it.

Your job is now to add the necessary settings to the `splitdim` container template to map the `splitdim-config` ConfigMap into the environment variables `KVSTORE_MODE` and `KVSTORE_ADDR` and redeploy the Kubernetes manifests. Make some quick tests with `curl` to see if everything is fine

> ✅ **Check**
> 
> Test your Kubernetes deployment through the following steps:
> - choose the local data layer and restart the `splitdim` Deployment to actually use the new settings:
>   ```shell
>   kubectl apply -f - <<EOF
>   apiVersion: v1
>   kind: ConfigMap
>   metadata:
>     name: splitdim-config
>   data:
>     kvstoreMode: "local"
>   EOF
>   kubectl rollout restart deployment splitdim
>   ```
> - wait a bit until the pod restarts and run the tests:
>   ```shell
>   cd 99-labs/code
>   export EXTERNAL_IP=$(kubectl get service splitdim -o jsonpath='{.status.loadBalancer.ingress[0].ip}')
>   export EXTERNAL_PORT=80
>   go test ./... --tags=httphandler,api,localconstructor,reset,transfer,accounts,clear -v -count 1
>   ```
> - restart the app with the `kvstore` data layer:
>   ```shell
>   kubectl apply -f - <<EOF
>   apiVersion: v1
>   kind: ConfigMap
>   metadata:
>     name: splitdim-config
>   data:
>     kvstoreMode: "kvstore"
>     kvstoreAddr: "kvstore.default:8081"
>   EOF
>   kubectl rollout restart deployment splitdim
>   ```
> - again, wait until `splitdim` restarts and rerun the tests:
>   ```shell
>   go test ./... --tags=httphandler,api,localconstructor,reset,transfer,accounts,clear -v -count 1
>   ```
> If all goes well, you should see all tests to PASS.

> [!WARNING]
>
> Unfortunately, applying a new ConfigMap will not automatically restart the Kubernetes workload (e.g., the Deployment) that depends on it. This means that every time you modify the ConfigMap you have to manually restart the app with `kubectl rollout restart deployment splitdim`. It is, however, simple to add the automation to Kubernetes necessary to support this functionality: this [tutorial](https://book.kubebuilder.io/reference/watching-resources/externally-managed.html) shows how to write a Kubernetes operator that will allow the user to specify pairs of a Deployment and a corresponding ConfigMap whose update should restart the Deployment. This *controller pattern* is the magic sauce that allows Kubernetes to support so many use cases.

Unfortunately that is all that we could cover from manageability in a single lab. That doesn't mean Kubernetes does not provide lots of further options for managing your app in the Cloud. If you want to learn more, we recommend the [Kubebuilder book](https://book.kubebuilder.io) and the [Operator SDK documentation](https://sdk.operatorframework.io/docs) as absolutely invaluable resources on the subject of manageability in Kubernetes.

<!-- Local Variables: -->
<!-- mode: markdown; coding: utf-8 -->
<!-- eval: (auto-fill-mode -1) -->
<!-- visual-line-mode: 1 -->
<!-- markdown-enable-math: t -->
<!-- End: -->
