# Steps

## Remove extra software packages
The following command removes most of the extra packages.
```console
 sudo DEBIAN_FRONTEND=noninteractive apt-get autoremove -yq \
	transmission* \
	thunderbird* \
	gimp* \
	libreoffice* \
	parole* \
	xfburn* \
	codeblocks* \
	hexchat* \
	rhythmbox*
```

> **Note**
>
> This is not a complete list, feel free to remove additional software packages

> **Warning**
>
> Do not remove `python2`, the [circle agent](https://git.ik.bme.hu/CIRCLE3/agent), which starts the VM, uses python2 (as of 2023-09).

## Install packages

> **Note**
>
> `apt dist-upgrade` might fail due to insufficient space on `/boot`. Manually delete old kernels to fix install error.

```console
 sudo sed -i "s@http://de.archive.ubuntu.com@http://hu.archive.ubuntu.com@g" /etc/apt/sources.list

 sudo DEBIAN_FRONTEND=noninteractive apt-get update
 sudo DEBIAN_FRONTEND=noninteractive apt-get dist-upgrade -yq
 sudo DEBIAN_FRONTEND=noninteractive apt-get install -yq \
	bash-completion \
	build-essential \
	tcpdump \
	git \
	curl \
	jq \
	podman \
	podman-docker \
	firefox \
	wireshark-gtk \
	fonts-noto-color-emoji

 export CNI_PLUGIN_DEB="containernetworking-plugins_1.1.1+ds1-3ubuntu0.24.04.3_amd64.deb"
 wget http://hu.archive.ubuntu.com/ubuntu/pool/universe/g/golang-github-containernetworking-plugins/$CNI_PLUGIN_DEB
 sudo dpkg -i $CNI_PLUGIN_DEB
 rm $CNI_PLUGIN_DEB
```

> **Note**
>
> Do not forget to clean apt caches.
```console
 sudo DEBIAN_FRONTEND=noninteractive apt-get autoremove -yq
 sudo DEBIAN_FRONTEND=noninteractive apt-get clean all
```

## Install Go
This snippet installs the latest stable Go version:
```console
 # export GO_TAR="$(curl -s https://go.dev/VERSION?m=text | head -n 1).linux-amd64.tar.gz"
 export GO_TAR="go1.24.6.linux-amd64.tar.gz"
 wget "https://go.dev/dl/$GO_TAR"
 sudo rm -rf /usr/local/go && sudo tar -C /usr/local -xzf $GO_TAR
 echo "export PATH=$PATH:/usr/local/go/bin" | sudo tee /etc/profile
 rm -rf $GO_TAR
```

> **Note**
>
> Use the commented-out `export GO_TAR` line to set version.

## Install kubectl
Install kubectl from the repo and enable bash completion:
```console
 export KUBEVERSION="1.30"
 echo "deb [signed-by=/etc/apt/keyrings/kubernetes-apt-keyring.gpg] https://pkgs.k8s.io/core:/stable:/v$KUBEVERSION/deb/ /" | sudo tee /etc/apt/sources.list.d/kubernetes.list
 sudo chmod 644 /etc/apt/sources.list.d/kubernetes.list
 curl -fsSL "https://pkgs.k8s.io/core:/stable:/v$KUBEVERSION/deb/Release.key" | sudo gpg --dearmor -o /etc/apt/keyrings/kubernetes-apt-keyring.gpg
 sudo chmod 644 /etc/apt/keyrings/kubernetes-apt-keyring.gpg
 sudo DEBIAN_FRONTEND=noninteractive apt-get update
 sudo DEBIAN_FRONTEND=noninteractive apt-get install -y kubectl
 echo "source <(kubectl completion bash)" >> ~/.bashrc
```

## Install VSCode
Download VSCode from the apt repo and install relevant packages:
```console
 wget -qO- https://packages.microsoft.com/keys/microsoft.asc | gpg --dearmor > packages.microsoft.gpg
 sudo install -D -o root -g root -m 644 packages.microsoft.gpg /etc/apt/keyrings/packages.microsoft.gpg
 sudo sh -c 'echo "deb [arch=amd64,arm64,armhf signed-by=/etc/apt/keyrings/packages.microsoft.gpg] https://packages.microsoft.com/repos/code stable main" > /etc/apt/sources.list.d/vscode.list'
 rm -f packages.microsoft.gpg
 sudo DEBIAN_FRONTEND=noninteractive apt-get update
 sudo DEBIAN_FRONTEND=noninteractive apt-get install -y code
 code --install-extension golang.go redhat.vscode-yaml

 go install github.com/go-delve/delve/cmd/dlv@latest
 go install honnef.co/go/tools/cmd/staticcheck@latest
 go install golang.org/x/tools/gopls@latest
```

## Install Minikube
Download, install and configure latest minikube:
```console
 # curl -LO https://storage.googleapis.com/minikube/releases/latest/minikube-linux-amd64
 curl -LO https://github.com/kubernetes/minikube/releases/download/v1.33.1/minikube-linux-amd64
 sudo install minikube-linux-amd64 /usr/local/bin/minikube
 rm minikube-linux-amd64
 echo "source <(minikube completion bash)" >> ~/.bashrc
 minikube config set WantUpdateNotification false
 minikube start --driver=podman --container-runtime=cri-o --download-only
 # minikube start --driver=podman --container-runtime=cri-o --delete-on-failure --wait-timeout 180s
 # minikube stop
 # minikube delete --all
```

## Configure podman
Surpress docker warnings and add dockerhub to registries:
```console
 sudo touch /etc/containers/nodocker
 echo 'unqualified-search-registries = ["docker.io"]' | sudo tee -a /etc/containers/registries.conf
```

## Additional (manual) configuration steps

- disable `xeyes` autostart in xfce settings / autostart
- update snaps: `sudo snap refresh`
- tweak dns settings: set systemd-resolved `StartLimitBurst=10`
