# Kubernetes Deployment Observer

A small Go service that watches Kubernetes Deployments and reports changes in their rollout state. I’m building it to learn how `client-go` informers work and to explore how a developer platform could surface useful deployment health information.

## First milestone

Watch Deployments in a local kind cluster and log their namespace, name, desired replicas, and available replicas when they are added or updated.

## How it works

The observer connects to Kubernetes using a local kubeconfig when run from a terminal, or service account credentials when run inside a Pod. A shared informer watches Deployment objects. Its handlers pass each Deployment to a function that describes its current state.

## Project layout

* `cmd/observer/main.go` — parses startup flags and runs the observer.
* `internal/observer/observer.go` — configures the Kubernetes client and informer.
* `internal/observer/deployment.go` — describes Deployment state.

## Planned next steps

* Log only meaningful rollout changes.
* Watch Pods and Kubernetes Events to explain why a Deployment is unhealthy.
* Expose the resulting assessment to other tools through an API or event stream.

**Status:** The client and informer factory setup is in progress; Deployment watch handlers are next.
