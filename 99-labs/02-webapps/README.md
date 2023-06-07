# Deploying a web app into Kubernetes

In this lab we build a simple web app in Go and learn how to deploy it into Kubernetes. Each section is closed with an exercise and a test that you can run to check whether you successfully completed the exercise.

## Table of Contents
1. [A basic web service in Go](#a-basic-web-service-in-go)
1. [Building a container image](#building-a-container-aimage)
1. [Deploy into Kubernetes](#deploy-into-kubernetes)
1. [Cleanup](#cleanup)

## A basic web service in Go

One of the greatest strengths of Go lies in the its suitability for developing web applications. It has a HTTP server as part of the standard library, offers great performance, and it is easy to deploy as a container into Kubernetes. This exercise will walk you through building a basic "hello-world" web application with Go, which we will deploying into Kubernetes later. 

### A Go web server

The below simple Go program will open a HTTP port at 8080, listen to HTTP queries on the path `/` and respond with the standard greeting `Hello, world!`. 

``` go
package main

import (
	"fmt"
	"net/http"
)

func HelloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello, world!")
}

func main() {
	http.HandleFunc("/", HelloHandler)
	http.ListenAndServe(":8080", nil)
}
```

The first line, `package main` declares that the code in the `main.go` file belongs to the `main` package. In the next few lines, the `net/http` and the `fmt` packages are imported into the file. The former provides the HTTP server implementation we use in our app, while the latter helps access operating system functionality. The `HelloHandler` function is a standard HTTP request handler with the signature `func(w http.ResponseWriter, r *http.Request)`, where `r` can be used to access the details of the HTTP request and `w` can be used to write the response. We are now using `fmt.Printf` to write the string `Hello, world!` to the HTTP response. When not specifying otherwise, Go will automatically set the HTTP status 200 (`OK`) in the response. 

The `main` function first assigns the `HelloHandler` request handler to the HTTP path `/`, meaning that whenever the server is called with an empty path (say, using the URL `http://example.com/`) it will automatically call our handler that will respond with the `Hello, world!` greeting. Finally, `http.ListenAndServe` spawns the HTTP server on port 8080. This function is blocking, so the program won't exit unless it encounters en error (or it is explicitly killed)

### Test

Copy the code into a new file named `main.go` and execute it with `go run main.go` (you can also use `go build` that will automatically build the `main` package into an executable). If all goes well, the prompt should disappear: the HTTP server silently starts and awaits requests. Send one using the omnipotent `curl` tool:

``` sh
curl http://localhost:8080/
Hello, world!
```

> **Note**: The address `localhost` is the short name for the loopback interface address, which defaults to `127.0.0.1`. Also note that the trailing slash is optional, so the following request would be equivalent: `curl http://127.0.0.1:8080`

Congratulations, you have build and run your first Go web app! Now fire up your favorite browser and direct it to the URL `localhost:8080` and see your first webpage rendered in full glory.

### Exercise

Modify the program to return the hostname of the server it is running at:
- first, import the `os` package from the Go standard lib: add `"os"` to the `import` list in parentheses;
- declare a global variable called `hostname` of type string that we will use to store the hostname: `var hostname string`;
- before starting the web server in the `main` function, store the hostname in the global variable;
- first query the hostname: `h, err := os.Hostname()`, where `h` is a temporary variable for the return value and `err` will contain the error if something go wrong;
- make sure to check the error before proceeding: 
  ```go
  	if err != nil {
		panic(err)
	}
  ```
- store the returned hostname in the global variable so that the HTTP handler, which runs in a separate function and so cannot reach `h`, will access it: `hostname = h`
- finally, modify the HTTP handler `HelloHandler` to add value of the global `hostname` variable to the response: `fmt.Fprintf(w, "Hello world from %s!", hostname)`

> ✅ Check
>
> Run the below first test to make sure that you have successfully completed the first exercise
> ``` sh
> go test -run TestHelloWorldLocal
> PASS
> ```

## Building a container image

### The Dockerfile

### Docker/podman workflow

### Building images in Minikube

## Deploy into Kubernetes

Kubernetes is undoubtedly the present and the future of managing containerized microservice applications.

A Kubernetes *cluster* consists of a set of worker machines (*nodes*) that run containerized applications, plus a *control plane* that manages the worker nodes and the containers, or the containers co-packaged into a Kubernetes *pod*, running on them. The control plane components make global decisions about the cluster and maintain cluster state. In production settings some nodes are dedicated to run the control plane separately from the worker nodes that execute the *workload* (the pods), but control plane components can be run on any node and can share the same node as the workers. This allows to run Kubernetes on a single node (e.g., Minikube).

<img src="fig/kubernetes_arch.png" alt="Kubernetes architecture", longdesc="https://sysdig.com/blog/monitor-kubernetes-api-server">

The most important component of the control plane is the `kube-apiserver` that exposes the Kubernetes API, acting as a frontend for the Kubernetes control plane. The cluster state is maintained in `etcd`, a consistent and highly-available key value store used as Kubernetes backing store for all cluster data. The `kube-scheduler` watches for newly created pods with no assigned node and selects a node for them to run on, the *kube-controller-manager* implements the common Kubernetes control loops, and `kube-dns` implements a cluster-wide DNS service which allows access to ephemeral pod and service IP addresses using short DNS domain names. 

Each worker node runs `kubelet`, which manages the containers/pods are running on the node, ~kube-proxy~, an implementation of the Kubernetes Service concept (see later), and a *container runtime* that is the software that is responsible for running the containers themselves.

### Using `kubectl`

The simplest way to interact with a Kubernetes cluster is via the `kubectl` command line interface (CLI). The CLI exposes the most common operations as simple *commands*, translates the commands into REST API requests, sends the requests along to the Kubernetes API server (the `kube-apiserver`) for execution and reports the results back to the user in a human readable for.

Before reaching 

- Set up autocomplete in bash into the current shell.
  ```sh
  source <(kubectl completion bash)
  ```
- Add autocomplete permanently to your bash shell.
  ```sh
  echo "source <(kubectl completion bash)" >> ~/.bashrc # add 
  ```
- Type `kubectl clu<TAB>` to check if autocompletion works: the shell should automatically complete your command to `kubectl cluster-info`.

> **Note**: Never use CLIs without autocompletion and always reach out for `<TAB>` while typing. If a shell does not provide autocompletion, don't use it.

Some useful `kubectl` commands:
- `kubectl config view`: view current client config, like API server URL and access tokens;
- `kubectl apply -f <file>`: bring the cluster to a desired state specified in the YAML `<file>`;
- `kubectl get <resources>`: query the current config of Kubernetes resources, where `<resources>` can be a resource name in plural or singular (`deployment`, `service`, `pod`), a shortname (`deploy`, `svc`, `pod`) or the full name (e.g., `deployment.apps`);
- `kubectl get <resource> <name>`: query the current state of the Kubernetes resources called `<name>` of type `<resource>`
- `kubectl delete <resource> <name>`: delete the resource called `<name>` of type `<resource>`;
- `kubectl describe <resource> <name>`: obtain a longer description of the `<resource>` called `<name>`;
- `kubectl explain <resource>`: get a human-readable description of a Kubernetes resource;
- `kubectl edit <resource> <name>`: online edit the YAML configuration of the `<resource>` called `<name>`;
- `kubectl create deployment <name> --image=<image-name>`: create a Deployment called `<name>` using the image `<image-name>`;
- `kubectl scale deployment <name> --replicas=<count>`: scale the deployment called `<name>` to `<count>` pods;
- `kubectl expose <name> --port=<port>`: make the deployment called `<name>` publicly available on port `<port>`;
- `kubectl label <resource> <name> <label-name>=<label-value>`: set the label `<label-name>` to `<label-value>` on the `<resource>` called `<name>`;
- `kubectl logs <name>`: display the logs from the pod called `<name>`;
- `kubectl attach <name> -i`: attach to the console of the pod called  `<name>` (type `Ctrl-C` to exit);
- `kubectl exec <name> -it  -- <command>`: execute the shell `<command>` on the pod called `<name>` (works only if the container image has a shell, containers built from Go programs usually don't have one unless we explicitly add it during the build);
- `kubectl cp`: copy files to/from pods;
- `kubectl port-forward service <name> <port>`: listen on local port 5000 and forward to port 5000 on the service running in Kubernetes.

### Kubernetes manifests

### Deployments

### Services

test Using the API server proxy

### Expose

#### NodePort

#### LoadBalancer

### Scale

## Cleanup
