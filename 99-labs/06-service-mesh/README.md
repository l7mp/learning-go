# Service mesh

## Preparations


### Using `istioctl`

store istio install dir into `$ISTIO_DIR`

``` bash
export ISTIO_DIR=<path-to-istio-root-dir>
```

istioctl search path

``` bash
export PATH=${ISTIO_DIR}/bin:${PATH}
```

check istioctl

``` bash
istioctl version
```

### Enable tracing

Collect traces for every 2nd call in Jaeger

``` bash
istioctl install -y -f - <<EOF
apiVersion: install.istio.io/v1alpha1
kind: IstioOperator
spec:
  meshConfig:
    enableTracing: true
    defaultConfig:
      tracing:
        sampling: 50
EOF
```


## Ingress gateway

## Traffic management

## Observability

Demonstrate Prom, Kiali and Grafana

```shell script
bin/istioctl install --set profile=minimal -y
kubectl label namespace default istio-injection=enabled --overwrite
kubectl apply -f samples/addons  # install all manifests in dir
```

``` bash
export EXTERNAL_IP=$(kubectl get service leaderboard -o jsonpath='{.status.loadBalancer.ingress[0].ip}')
watch -n .1 curl -o /dev/null -s http://${EXTERNAL_IP}:8080/api/getScores
```

``` bash
bin/istioctl dashboard kiali
```

choose "Graph", set namespace to default and choose "Workload graph" and set 

``` bash
kubectl apply -f - <<EOF
apiVersion: networking.istio.io/v1alpha3
kind: VirtualService
metadata: { name: kvstore-500 }
spec:
  hosts: [ kvstore ]
  http:
    # requests to the "/api/list" path will return 500 status every 10th time
    - match: [ uri: { exact: "/api/list" } ]
      fault: { abort: { httpStatus: 500, percentage: { value: 10 } } }
      route: [ destination: { host: kvstore } ]
    # default route: everything that is not a "put" ("list" and "get")
    - route: [ destination: { host: kvstore } ]
EOF
```




<!-- THIS DOESN'T WORK FOR SOME REASON: you can also finetune the sampling rate via Istio's Telemetry API -->
<!-- ```shell script -->
<!-- kubectl apply -f - <<EOF -->
<!-- apiVersion: telemetry.istio.io/v1alpha1 -->
<!-- kind: Telemetry -->
<!-- metadata: -->
<!--   name: mesh-default -->
<!--   namespace: istio-system -->
<!-- spec: -->
<!--   # no selector specified, applies to all workloads -->
<!--   tracing: -->
<!--   - randomSamplingPercentage: 100.0 -->
<!-- EOF -->
<!-- ``` -->

> ✅ **Check**
> Test your Kubernetes deployment. Some useful commands for testing from the shell:

<!-- Local Variables: -->
<!-- mode: markdown; coding: utf-8 -->
<!-- auto-fill-mode: nil -->
<!-- visual-line-mode: 1 -->
<!-- markdown-enable-math: t -->
<!-- End: -->
