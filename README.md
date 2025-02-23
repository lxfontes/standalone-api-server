A barebones Kubernetes API Server, for non-Kubernetes use-cases.

Similar projects:

- https://www.kcp.io/

## Why?

Kubernetes has many layers of abstraction, and is not always the best fit for every use-case.
The API Server is a powerful tool, and can be used to build custom APIs, without the need for a full Kubernetes cluster.
This allows reuse of existing Kubernetes tooling, and the ability to leverage the Kubernetes API patterns such as CRDs, RBAC, ServiceAccounts, Finalizers, etc.

## Is this a hack?

No. The API Server code used here is straight from Kubernetes.

It is also the same approach used by lightweight Kubernetes distributions like K3s, K0s.

## Example

The [example](./example) directory shows the minimal setup required to run a custom API Server, it includes:

- A Private CA structure
- A sample token file, listing users & groups
- A kubeconfig file, with a user & context pointing to the API server
- A script showing the steps to create the kubeconfig file

How to use:

Install KINE (ETCD API backed by sqlite):

```
go install github.com/k3s-io/kine@v0.13.10
```

Install Prox ( Procfile Runner ):

```
go install github.com/fgrosse/prox/cmd/prox@latest
```

Run `prox` in the [example](./example) directory.

In another shell, use:

```
export KUBECONFIG=example/kubeconfig
kubectl api-resources
kubectl auth whoami
```
