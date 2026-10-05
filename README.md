# Kubernetes Observer

A small Go service that watches Kubernetes resources and logs changes as they
happen. The project is an exploration of `client-go` shared informers and the
building blocks of deployment-health tooling.

## Current behavior

The service runs two components concurrently:

* A Kubernetes observer that watches Deployments, Pods, and Events across all
  namespaces.
* An HTTP server that listens on port `8080` and responds to `/` with
  `Hello world from server!`.

The observer logs add, update, and delete notifications for each watched
resource:

* Deployments: namespace/name and available versus desired replicas.
* Pods: namespace/name and current phase.
* Events: namespace/name, reason, and message.

Informer caches are resynchronized every 10 minutes. On `SIGINT` or `SIGTERM`,
the observer stops and the HTTP server is given up to five seconds to shut down
gracefully.

## Requirements

* Go 1.27 or newer, as declared in `go.mod`.
* Access to a Kubernetes cluster.
* Credentials allowed to `list` and `watch` Deployments, Pods, and Events in
  the namespaces the observer can access.

No Kubernetes manifests or container image configuration are included yet.

## Run locally

By default, the service reads the file named by `KUBECONFIG`. If that variable
is unset, it uses `$HOME/.kube/config`.

```sh
go run ./cmd/observer
```

In another terminal, check the HTTP server:

```sh
curl http://localhost:8080/
```

Create, update, or delete Kubernetes resources to see observer messages in the
service logs. Stop the service with `Ctrl+C`.

## Run in a cluster

When the process is running in a Pod with an appropriately configured service
account, select Kubernetes in-cluster credentials with:

```sh
go run ./cmd/observer -in-cluster
```

The service account needs permission to list and watch `deployments` in the
`apps` API group and `pods` and `events` in the core API group. Because the
current informer factory is not namespace-scoped, cluster-wide RBAC is needed
to observe every namespace.

## Project layout

* `cmd/observer/main.go` — parses flags and coordinates the observer and HTTP
  server lifecycle.
* `internal/observer/observer.go` — loads Kubernetes credentials, creates the
  client and shared informer factory, and manages cache synchronization.
* `internal/observer/deployment.go` — watches and describes Deployments.
* `internal/observer/pods.go` — watches and describes Pods.
* `internal/observer/events.go` — watches and describes Kubernetes Events.
* `internal/observer/log.go` — provides common resource-event logging.
* `internal/server/server.go` — runs the HTTP server and handles graceful
  shutdown.

## Development

Build or test all packages with:

```sh
go build ./...
go test ./...
```

## Status and next steps

The Kubernetes client, shared informer factory, resource handlers, cache
synchronization, and basic HTTP server are implemented. The observer currently
logs every informer notification independently; it does not yet correlate Pods
and Events with their owning Deployment or expose an aggregate health model.

Likely next steps include:

* Filter updates so only meaningful state changes are logged.
* Correlate Deployments with their ReplicaSets, Pods, and Events.
* Replace the placeholder HTTP response with health or deployment-status
  endpoints.
* Add deployment manifests, RBAC manifests, automated tests, and container
  packaging.
