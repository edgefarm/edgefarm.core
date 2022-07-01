# edgefarm.core

edgefarm.core extends an existing kubernetes cluster by secured cloud/edge computing.

edgefarm.core uses different open source tools to provide edge node integration,
secure node registration and isolation.

## Detail

### kubeedge

Using kubeedge, special kubernetes edge nodes (edge devices) can be added to the cluster.
The edge node behaves like a standard kubernetes node, with the difference that workload
can still be operated reliably even if the connection to the kubernetes cluster is temporarily unavailable.

### metrics-server

Metrics Server is a scalable, efficient source of container resource metrics for Kubernetes built-in autoscaling pipelines.

Metrics Server collects resource metrics from Kubelets and exposes them in Kubernetes apiserver through Metrics API for use by Horizontal Pod Autoscaler and Vertical Pod Autoscaler. Metrics API can also be accessed by kubectl top, making it easier to debug autoscaling pipelines.

## Developing edgefarm.core

To set up a local development environment there is a devspace.yaml in the `/dev` subfolder which can be used directly.
The devspace setup relies on k3d to manage local kubernetes clusters.

Dependencies:

- [devspace](https://devspace.sh/)
- [kind](https://kind.sigs.k8s.io)
- kubectl
- kustomize
- helm
- [mkcert](https://github.com/FiloSottile/mkcert)

There are some predefined handy commands that simplifies the setup process.

`devspace run init`: Initialization with cluster setup.

`devspace run purge`: Remove all created resources, incl. cluster.

`devspace run activate`: Set the kubernetes context pointing to the cluster.

`devspace run update`: Update all dependencies.

To init and create a new environment, execute the following commands:

```sh
cd ./dev
devspace run init
devspace deploy
```

To apply your modifications, rerun `devspace deploy`.

To handle virtual kubeedge nodes use the following commands:

```sh
cd ./dev
# Start virtual nodes. See ./dev/hack/kubeedge-node/manifests/deployment.yaml for settings
# replicas: the number of virtual devices wanted
# env CLOUDCORE_ADDRESS: the hostname of cloudcore. This might be in the format of 192-168-1-42.nip.io
devspace run instantiate-nodes

# Delete virtual nodes and all workload they execute
# Note: deleting the workload pods from the cluster takes a while until the scheduler detects 
# that the nodes are missing
devspace run purge-nodes
```

To cleanup the environment, execute the following command:

```sh
cd ./dev
devspace run purge
```