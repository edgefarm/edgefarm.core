# edgefarm.core

edgefarm.core extends an existing kubernetes cluster by secured cloud/edge computing.

edgefarm.core uses different open source tools to provide edge node integration,
secure node registration and isolation.

## Detail

### kubeedge

Using kubeedge, special kubernetes edge nodes (edge devices) can be added to the cluster.
The edge node behaves like a standard kubernetes node, with the difference that workload
can still be operated reliably even if the connection to the kubernetes cluster is temporarily unavailable.

### vault

KubeEdge implements only a rudimentary authentication procedure for edge nodes based on
self-generated tokens. By interacting with vault, a variety of authentication methods
become possible. Furthermore, features like certificate renewals and revocation lists can be implemented.

### capsule

tbd

### crossplane

tbd

## Developing edgefarm.core

To set up a local development environment there is a devspace.yaml in the `/dev` subfolder which can be used directly.
The devspace setup relies on k3d to manage local kubernetes clusters.

Dependencies:

- [devspace](https://devspace.sh/)
- [k3d](https://k3d.io/)
- kubectl
- kustomize
- helm
- [mkcert](https://github.com/FiloSottile/mkcert)

There are some predefined handy commands that simplifies the setup process.

`devspace run init`: Initialization with k3d cluster setup.

`devspace run purge`: Remove all created resources, incl. k3d cluster.

`devspace run activate`: Set the kubernetes context pointing to the cluster.

`devspace run update`: Update all dependencies.

To init and create a new environment, execute the following commands:

```sh
cd ./dev
devspace run init
devspace deploy
```

To apply your modifications, rerun `devspace deploy`.
