# WebApp Operator

This directory contains a Go-based Kubernetes operator built with Kubebuilder and controller-runtime.

The operator introduces a custom Kubernetes resource called `WebApp`. The goal is to provide a simple platform abstraction that allows developers to deploy web applications without writing raw Kubernetes Deployment and Service manifests.

## Purpose

In a platform engineering context, teams often expose higher-level abstractions to developers in order to standardize deployments and reduce repetitive Kubernetes configuration.

The `WebApp` custom resource is a simplified example of this pattern.

Instead of asking developers to manually define a Deployment and a Service, they can create a `WebApp` object:

```yaml
apiVersion: platform.dylan.io/v1
kind: WebApp
metadata:
  name: demo-webapp
spec:
  image: nginx:1.27
  replicas: 2
  port: 80
```

The operator reconciles this desired state into native Kubernetes resources.

## What the operator manages

For each `WebApp` resource, the operator creates and maintains:

- a Kubernetes Deployment
- a Kubernetes Service

The Deployment is configured with:

- the container image defined in `spec.image`
- the replica count defined in `spec.replicas`
- the container port defined in `spec.port`

The Service exposes the application internally inside the cluster.