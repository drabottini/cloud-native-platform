# Cloud Native Platform Lab

A personal platform engineering lab demonstrating how to build and operate a small cloud-native platform using Kubernetes, GitOps, CI/CD, observability and a custom Go operator.

## Goals

- Deploy applications to Kubernetes using GitOps
- Manage environment-specific configuration with Kustomize
- Validate Kubernetes manifests through CI
- Build a custom Kubernetes operator in Go
- Add observability with Prometheus, Grafana and Loki
- Integrate secrets management with Vault
- Document platform engineering and SRE practices

## Architecture

Developer -> GitHub -> GitHub Actions -> Kubernetes manifests validation -> ArgoCD -> k3d/k3s cluster -> workloads -> Prometheus/Grafana/Loki

## Current Status

- Kubernetes base/overlay structure: done
- GitHub Actions manifest validation: done
- ArgoCD application structure: in progress
- Go Kubernetes operator: in progress
- Observability stack: planned
- Vault integration: planned
- Terraform automation: planned