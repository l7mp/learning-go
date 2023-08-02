# Setting up the work environment

In this lab we set up the workspace for developing and testing our cloud-native applications implemented in this course.

## Table of Contents

1. [System Requirements](#system-requirements)
1. [Installation](#installation)
1. [Test](#test)
1. [Exercises](#exercises)

## System Requirements

We provide instructions for setting up a native workspace on GNU/Linux. Windows Subsystem for Linux, Windows, or macOS might work, but we recommend a GNU/Linux virtual machine instead. A ready-to-use virtual machine will be distributed.

### System Parameters

**Operating System:** we provide instructions for Ubuntu 22.04. The lab should work on any major GNU/Linux distribution, feel free to adapt the steps.
**CPU:** 2-core x86_64 CPU should be sufficient for native installation, at least 4 cores are required for the VM.
**Memory:** 4GB for native install, 8GB is recommended for the VM.

> Warning
> Make sure your Internet connection is working, we will download software packages.

## Installation

The course requires the following software to be installed:
- [Go programming language](https://go.dev/)
- [Visual Studio Code](https://code.visualstudio.com/) or any editor
- [podman](https://podman.io/) or [docker](https://www.docker.com/)
- [minikube](https://minikube.sigs.k8s.io/docs/)
- [kubectl](https://kubernetes.io/docs/tasks/tools/)
- [jq](https://jqlang.github.io/jq/)
- `make`
- [Istio and istioctl](https://istio.io/)


> Note
> If you use the VM, jump to [Insall Istio](#install-istio)

### Install Go

The following snippet implements the [official install](https://go.dev/doc/install). Copy-paste it into a terminal window.

```shell
export GO_TAR="$(curl -s https://go.dev/VERSION?m=text).linux-amd64.tar.gz"
wget "https://go.dev/dl/$GO_TAR"
sudo rm -rf /usr/local/go && sudo tar -C /usr/local -xzf $GO_TAR
echo "export PATH=$PATH:/usr/local/go/bin" | sudo tee /etc/profile
rm -rf $GO_TAR
```

### Install VS Code

To ease writing go programs and to editing configuration files, you will need an editor. We know Emacs and Visual Studio Code, but feel free to use any editor that you like and supports Go.

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
Execute the following commands in your terminal to install and configure podman.

1. Install podman via the package manager
```shell
sudo apt-get install -y podman podman-docker
```

1. Configure podman to enable access to [Dockerhub](https://hub.docker.com/)
```shell
sudo touch /etc/containers/nodocker
echo 'unqualified-search-registries = ["docker.io"]' | sudo tee -a /etc/containers/registries.conf
```

### Install command line tools

In the course we use `jq` to observer JSON files in a human-readable form. The homeworks rely on `make`. To install these tools via the package manager, execute this command in your terminal:
```code
sudo apt-get install -y make jq
```

### Install minikube

In your terminal execute the following commands to download and install the latest minikube:
```code
curl -LO https://storage.googleapis.com/minikube/releases/latest/minikube-linux-amd64
sudo install minikube-linux-amd64 /usr/local/bin/minikube
rm minikube-linux-amd64
echo "source <(minikube completion bash)" >> ~/.bashrc
```

### Install kubectl
Execute the following commands in your terminal to install and configure kubectl.

1. Install `kubectl` via the package manager
```code
sudo curl -fsSLo /etc/apt/keyrings/kubernetes-archive-keyring.gpg https://packages.cloud.google.com/apt/doc/apt-key.gpg
echo "deb [signed-by=/etc/apt/keyrings/kubernetes-archive-keyring.gpg] https://apt.kubernetes.io/ kubernetes-xenial main" | sudo tee /etc/apt/sources.list.d/kubernetes.list
sudo apt-get install -y kubectl
```

1. Enable shell completion for bash
```code
echo "source <(kubectl completion bash)" >> ~/.bashrc
```

> Note
> If you use a different shell (e.g., zsh) , check configuration steps with `kubectl completion -h`

### Install Istio
We will often use additional functionality from the [Istio service mesh](https://istio.io). Use the below to download and enable Istio:

``` sh
curl -L https://istio.io/downloadIstio | sh -
cd istio*
kubectl kustomize "github.com/kubernetes-sigs/gateway-api/config/crd?ref=v0.6.2" | kubectl apply -f -
bin/istioctl install --set profile=minimal -y
kubectl label namespace default istio-injection=enabled --overwrite
```

## Test

### Test that tools are installed

#### minikube

Execute in terminal to check the installed minikube's version:
```shell
minikube version
```

> ✅ **Check**
>
> This command should print the minikube version.

### Check Go environment

1. Paste the Go code below to your edtor, save it as `hello.go`

``` go
package main
import "fmt"

func main() {
	fmt.Println("Hello, world")
}
```

2. Compile/run it in your terminal

```code
go run hello.go
```

> ✅ **Check**
>
> The compilation should finish without any errors, and the program prints out the good old `Hello, world` on its output.


## Exercises

### Generate homeworks

### Solve the first homework
