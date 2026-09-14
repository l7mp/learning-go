# Labs

The course consists of labs that give you hands-on experience on both developing and deploying
cloud-native apps in Go. During the labs you will learn how to build, containerize, and deploy your
Go programs in the cloud. We will mostly work with Kubernetes, while last labs focus on service
meshes.

## Available labs

1. [Setting up the work environment](01-setup/)
1. [Deploying web applications into Kubernetes](02-webapps/)
1. [A silly web app: SplitDim](03-splitdim/)
1. [Cloud-native patterns: Immutability and loose coupling](04-immutability/)
1. [Cloud-native patterns: Resilience](05-resilience/)
1. [Cloud-native patterns: Manageability](06-manageability/)
1. [Service mesh](07-service-mesh/)

## Testing

Each lab README walks you through the checks by hand, which is how you are meant to work through
them.

You can use the testing harness that will execute all lab tests against your running `minikube`
instance (don't forget to run `minikube tunnel` and push your images): it will build no images and
create no cluster, and the tests use whatever `KUBECONFIG` names.

For this purpose, run `LABS_EXTERNAL_CLUSTER=1 LAB=<LAB-ID> make lab test` from the repository
root. For example, `LAB=05` means the 5th lab on resiliency; if you do not set `LAB`, it will test
*all* labs at once.

> **Warning**
> The fully automated tests described below work with docker only. In case you are running our
> prebuilt VM/cloud image, please use the method discussed above.

Once a lab is done, `LAB=<LAB-ID> make lab-test` from the repository root runs the automated tests
for that lab (so `LAB=05 make lab-test` will run only the 5th lab on resiliency): it builds your
container images, brings up a throwaway Kubernetes cluster, deploys your manifests and runs every
test against the cluster, giving you a single pass or fail. To check *all* labs at once, use `make
lab-test` (again from the repository root).

The test suite itself lives in [ci/](ci/); you are welcome to study it.

<!-- Local Variables: -->
<!-- mode: markdown; coding: utf-8 -->
<!-- auto-fill-mode: nil -->
<!-- visual-line-mode: 1 -->
<!-- markdown-enable-math: t -->
<!-- End: -->
