# Setting up the work environment

During this lab we will set up a workspace for developing and testing the Go applications we will implement. The workspace consists of a *Go development environment* (compiler, the `go` utility, standard libraries, etc.) complete with an Go-enabled editor, a local *Kubernetes cluster* that can be used to deploy and test Go code, and *Istio, a popular service mesh* distribution that makes it easier to operate our microservice applications.

## Table of Contents

1. [System Requirements](#system-requirements)
1. [Installation](#installation)
1. [Test](#test)
1. [Exercises](#exercises)

## System Requirements

We provide instructions for setting up a native workspace on GNU/Linux. Windows Subsystem for Linux, Windows, or macOS might work, but we recommend a GNU/Linux virtual machine instead. A ready-to-use [virtual machine](#vm-details) is available.

### System Parameters

**Operating System:** we provide instructions for Ubuntu 22.04. The lab should work on any major GNU/Linux distribution, feel free to adapt the steps.\
**CPU:** 2-core x86_64 CPU should be sufficient for native installation, at least 4 cores are required for the VM.\
**Memory:** 4GB for native install, 8GB is recommended for the VM.

> [!WARNING]
>
> Make sure your Internet connection is working, we will download software packages.

### VM Details

We provide a ready-to-use Ubuntu 22.04-based virtual machine. [Click here to download VM image](http://lendulet.tmit.bme.hu/~levai/files/go-vm/CloudGoVM.ova). \
To use the downloaded image, install [VirtualBox](https://www.virtualbox.org/wiki/Downloads) (tested with VirtualBox 7.0) and import the downloaded OVA file with Virtualbox.\
**Login details:**\
 username: `vagrant`\
 password: `vagrant`

You can also build the VM with [Vagrant](https://developer.hashicorp.com/vagrant/downloads) from the [Vagrantfile](/env/Vagranfile). The build takes roughly 1 hour.

## Installation

The course requires these software:
- [Go programming language](https://go.dev/)
- [Visual Studio Code](https://code.visualstudio.com/) or any editor
- [podman](https://podman.io/) or [docker](https://www.docker.com/)
- [kubectl](https://kubernetes.io/docs/tasks/tools/)
- [minikube](https://minikube.sigs.k8s.io/docs/)
- [jq](https://jqlang.github.io/jq/)
- `make`
- [Istio and istioctl](https://istio.io/)

> [!NOTE]
> If you use the VM, jump to the [Istio installation guide](#install-istio).

### Install Go

The following snippet implements the [official install](https://go.dev/doc/install). Copy-paste it into a terminal window.

```shell
export GO_TAR="$(curl -s https://go.dev/VERSION?m=text | head -n 1).linux-amd64.tar.gz"
wget "https://go.dev/dl/$GO_TAR"
sudo rm -rf /usr/local/go && sudo tar -C /usr/local -xzf $GO_TAR
echo "export PATH=$PATH:/usr/local/go/bin" | sudo tee /etc/profile
rm -rf $GO_TAR
```

### Install VS Code

To ease writing go programs and editing configuration files, you will need an editor. We use GNU Emacs and Visual Studio Code, but feel free to use any editor that you like and supports Go.

To install and setup Code, execute the following commands in your terminal:

1. Install Code with the package manager:

   ```shell
   wget -qO- https://packages.microsoft.com/keys/microsoft.asc | gpg --dearmor > packages.microsoft.gpg
   sudo install -D -o root -g root -m 644 packages.microsoft.gpg /etc/apt/keyrings/packages.microsoft.gpg
   sudo sh -c 'echo "deb [arch=amd64,arm64,armhf signed-by=/etc/apt/keyrings/packages.microsoft.gpg] https://packages.microsoft.com/repos/code stable main" > /etc/apt/sources.list.d/vscode.list'
   rm -f packages.microsoft.gpg
   sudo apt-get update
   sudo apt-get install -y code
   ```

1. Install plugins to have improve syntax highlighting, static analyzer, and additional useful features:

   ```shell
   code --install-extension golang.go redhat.vscode-yaml
   go install github.com/go-delve/delve/cmd/dlv@latest
   go install honnef.co/go/tools/cmd/staticcheck@latest
   go install golang.org/x/tools/gopls@latest
   ```

### Install Podman

Docker and podman are popular tools to build Linux container images. We will use podman since it is a bit easier to install (plus, it is written in Go!); feel free to use Docker (which is, surprise!, also written in Go) or any other tool instead. Execute the following commands in your terminal to install and configure podman.

1. Install podman via the package manager.

   ```shell
   sudo apt-get install -y podman podman-docker
   ```

1. Configure podman to enable access to [Dockerhub](https://hub.docker.com).

   ```shell
   sudo touch /etc/containers/nodocker
   echo 'unqualified-search-registries = ["docker.io"]' | sudo tee -a /etc/containers/registries.conf
   ```

### Install miscellaneous command line tools

We will often use `jq` to output JSON files in a human-readable form and the `make` utility to generate and test the homework exercises. To install these tools via the package manager, execute this command in your terminal:

```code
sudo apt-get install -y make jq
```

### Install kubectl

The `kubectl` utility is our main tool to interact with our Kubernetes clusters. Execute the following commands in your terminal to install and configure `kubectl`.

1. Install `kubectl` via the package manager

   ```code
   curl -fsSL https://pkgs.k8s.io/core:/stable:/v1.30/deb/Release.key | sudo gpg --dearmor -o /etc/apt/keyrings/kubernetes-apt-keyring.gpg
   echo "deb [signed-by=/etc/apt/keyrings/kubernetes-apt-keyring.gpg] https://pkgs.k8s.io/core:/stable:/v1.30/deb/ /" | sudo tee /etc/apt/sources.list.d/kubernetes.list
   sudo apt update
   sudo apt install -y kubectl
   ```

1. Enable shell completion for bash

   ```code
   echo "source <(kubectl completion bash)" >> ~/.bashrc
   ```

> [!NOTE]
> If you use a different shell (e.g., zsh) , check configuration steps with `kubectl completion -h`

### Install Minikube

Minikube is a lightweight Kubernetes distribution to deploy a simple local Kubernetes cluster containing only a single node (your own laptop). In your terminal execute the following commands to download and install the latest version of Minikube:

```code
curl -LO https://storage.googleapis.com/minikube/releases/latest/minikube-linux-amd64
sudo install minikube-linux-amd64 /usr/local/bin/minikube
rm minikube-linux-amd64
echo "source <(minikube completion bash)" >> ~/.bashrc
```

Once installed, we create a local Kubernetes cluster. Copy the below into your terminal:

```code
minikube start --memory=4096m --cpus=2 --driver=podman --container-runtime=cri-o
```

> [!NOTE]
> Cluster creation may take some time.

This will create a local Kubernetes cluster with 4 GB memory and 2 vCPUs using podman as the container driver, and configure `kubectl` to talk to this cluster. Feel free to customize the CPU/memory limits in the above; e.g., it is a good idea to increase the amount of CPU and memory available to your cluster to obtain a more responsive Kubernetes.

> [!NOTE]
> Once done working with Kubernetes make sure to stop it with `minikube stop`: Kubernetes may take up considerable resources and this commands frees those resources up. You can always restart your cluster with `minikube start` and continue working from where you left the last time you issued `minikube stop`.

### Install Istio

We will often use additional functionality from the [Istio service mesh](https://istio.io). After starting minikube, use the below to download and enable Istio:

``` sh
curl -L https://istio.io/downloadIstio | sh -
cd istio*
kubectl kustomize "github.com/kubernetes-sigs/gateway-api/config/crd?ref=v0.6.2" | kubectl apply -f -
bin/istioctl install --set profile=minimal -y
kubectl label namespace default istio-injection=enabled --overwrite
```

> [!WARNING]
>
> Make sure that `minikube` runs before installing Istio.

## Test

### Check your Go environment

1. Paste the Go code below to your edtor, save it as `hello.go`

``` go
package main
import "fmt"

func main() {
    fmt.Println("Hello, world")
}
```

2. Compile/run it in your terminal

```shell
go run hello.go
```

> ✅ **Check**
>
> The compilation should finish without any errors and the program should print the good old `Hello, world` greeting to the standard output.

### Test your Kubernetes cluster

Execute the below in a terminal to check  Minikube's installed version:

```shell
minikube version
```

> ✅ **Check**
>
> This command should print the minikube version.

Execute the below terminal to check whether your Kubernetes cluster is running and `kubectl` is correctly configured to talk to it.

```shell
kubectl cluster-info
```

> ✅ **Check**
>
> This command should print the running Kubernetes version with some additional information.

### Test your Istio install

In a terminal execute the below in a terminal to check that Istio is installed properly:

```shell
kubectl get all -n istio-system
```

> ✅ **Check**
>
> This command should print all the Kubernetes resources in the namespace `istio-system`. Observe the `READY` column values (should `1/1` or similar) to check all resources are up and running.

The output should be something like the following:
```
NAME                                        READY   STATUS    RESTARTS   AGE
pod/istio-ingressgateway-7f57549c9f-q9jhq   1/1     Running   0          10m
pod/istiod-ddcf4fdd9-kl5wq                  1/1     Running   0          15m

NAME                           TYPE           CLUSTER-IP      EXTERNAL-IP   PORT(S)                                      AGE
service/istio-ingressgateway   LoadBalancer   10.99.219.24    <pending>     15021:32645/TCP,80:31546/TCP,443:32140/TCP   10m
service/istiod                 ClusterIP      10.98.141.167   <none>        15010/TCP,15012/TCP,443/TCP,15014/TCP        15m

NAME                                   READY   UP-TO-DATE   AVAILABLE   AGE
deployment.apps/istio-ingressgateway   1/1     1            1           10m
deployment.apps/istiod                 1/1     1            1           15m

NAME                                              DESIRED   CURRENT   READY   AGE
replicaset.apps/istio-ingressgateway-7f57549c9f   1         1         1       10m
replicaset.apps/istiod-ddcf4fdd9                  1         1         1       15m

NAME                                                       REFERENCE                         TARGETS              MINPODS   MAXPODS   REPLICAS   AGE
horizontalpodautoscaler.autoscaling/istio-ingressgateway   Deployment/istio-ingressgateway   cpu: <unknown>/80%   1         5         1          10m
horizontalpodautoscaler.autoscaling/istiod                 Deployment/istiod                 cpu: <unknown>/80%   1         5         1          15m
bin/istioctl verify-install
```

## Exercises

The course comes with a set of exercises to practice the basics of Go programming (syntax, type system, concurrency primitives, etc.). The exercises are customized per each student; this is to increase the effort required to copy your solutions from elsewhere. Each exercise is randomly generated from a template using your student id as the random seed, which is supposed to be your student id. Once the exercises are generated, you can start to add your solutions and then run `make test` to check your solutions.

### Prerequisites

Create a local clone of this git repo:

``` shell
git clone https://github.com/l7mp/learning-go.git
cd learning-go
```
You should always add and commit your solutions to this repo (see below) to avoid losing your work.

> [!TIP]
>
> We recommend to keep a copy of your git tree somewhere safe to back up your solutions. The simplest way is to use a GitHub private fork for this purpose. We ask you to keep your GitHub repo private, in order to prevent others from copying your work.

### Generate the exercises

Change to the root of the git repo and make sure to read the instructions in the `README.md` file. The below summarizes the main steps.

``` shell
echo <MY-STUDENT-ID> > STUDENT_ID
make generate
```

> [!WARNING]
>
> You must use your own student id. We will check this, so make sure you do not mistype your id.

### Solve the first homework

Navigate into the directory `01-getting-started/01-hello-world` that contains the first exercise. You should see the following files (plus some invisible files with name starting with `.` that you can safely ignore):

- `README.md` contains the instructions for solving the exercise;
- `exercise.go`: is a placeholder for your solution;
- `exercise_test.go`: is a test file that will check if your solution is correct.

If any of these files is missing or contains a placeholder that means you forgot to generate the homeworks, so go back to the previous step.

Issue the below command to run the tests: this should fail as there is no solution yet.

``` shell
go test ./... -v -count 1
./exercise_test.go:10:18: undefined: helloWorld
FAIL	github.com/l7mp/learning-go/01-getting-started/01-hello-world [build failed]
FAIL
```

Consult the `README.md` file for how to solve the exercise and place your solution into the file `exercise.go` at the placeholder.

> [!NOTE]
> It is usually not worth copying someone else's solution: most probably your exercises will be quite different for theirs (that is what `make generate` is for).

For instance, you may be asked to write a `helloWorld` function in Go that will return the string `Hello world!` (your exercise may differ, so make sure you read the README carefully!). In this case, insert the below code into `exercise.go`:

```go
func helloWorld() string {
    return "Hello world!"
}
```

> ✅ **Check**
>
> Once correctly solved, all tests in the exercise should pass:
> ``` shell
> go test ./... -v -count 1
> === RUN   TestHelloWorld
> --- PASS: TestHelloWorld (0.00s)
> PASS
> ```

If all tests pass, then git-commit your solution: this makes sure it remains there even if you re-generate the exercises.

``` shell
git add exercise.go
git commit -m 'first exercise solved'
```

> [!NOTE]
> If you use a remote git repo to back up your work then make sure you push all your commits there using:
> ``` shell
> git push
> ```

You can test *all* your solutions from the main directory by issuing `make test`. Currently only the first test will succeed: at the end of the course you should have *all the tests* pass.

<!-- Local Variables: -->
<!-- mode: markdown; coding: utf-8 -->
<!-- eval: (auto-fill-mode -1) -->
<!-- visual-line-mode: 1 -->
<!-- markdown-enable-math: t -->
<!-- End: -->
