# Contributing

## Local Development

Start Rancher

```
scripts/start_rancher.sh
```

Login using the details in the output of the script.

Build the driver

```
go build -o kontainer-engine-driver-kamatera
```

Start a local http server to server the driver to Rancher

```
scripts/serve_driver.sh
```

See the script output for the URL to use in Rancher.

See Kamatera/ui-cluster-driver-kamatera for the UI component.

Login to Rancher -> Cluster Management -> Drivers -> Cluster Drivers -> Create

* Download URL: http://127.0.0.1:8944/kontainer-engine-driver-kamatera
* Custom UI URL: http://127.0.0.1:3000/component.js
* Whitelist Domains: 127.0.0.1

### Testing with local cloudcli-server k8s api

Start a local server according to the instructions in kamatera/cloudcli-server-kubernetes

Optionally, set DRY_RUN=true in the cloudcli-server-kubernetes environment - this allows to quickly test the driver without making changes.

In case you use DRY_RUN, start a local Kubernetes cluster with [kind](https://kind.sigs.k8s.io/) and place it in the tests directory of cloudcli-server-kubernetes:

```
kind create cluster
kind get kubeconfig > ../cloudcli-server-kubernetes/tests/.kubeconfig
```

Build the driver with the local server URL and debugging enabled to use from Rancher

```
go build -o kontainer-engine-driver-kamatera -ldflags "-X 'main.cloudcliBaseUrl=http://127.0.0.1:8000' -X 'main.cloudcliDebug=true'"
```

### Testing the GRPC server directly

Start the GRPC server

```
go run -ldflags "-X 'main.cloudcliBaseUrl=http://127.0.0.1:8000' -X 'main.cloudcliDebug=true'" github.com/Kamatera/kontainer-engine-driver-kamatera 8888
```

Make a request

```
tests/grpc_create.sh
```