# PowerDNS Platform

Enterprise Kubernetes operator and Helm chart for running PowerDNS Authoritative Server with a production PostgreSQL backend, designed for OpenStack Designate compatibility.

## What is included

- Go/Kubebuilder-style operator scaffold
- PowerDNSCluster CRD schema and sample resources
- production Helm chart scaffold with operational templates
- OpenStack Designate integration notes
- HA, backup, observability, security, and upgrade guidance

## Repository layout

- `cmd/manager` operator entrypoint
- `api/v1alpha1` custom resource definitions and types
- `internal/controller` reconciliation logic
- `charts/powerdns-platform` production Helm chart scaffold
- `config/crd` CRD definitions
- `config/samples` example PowerDNSCluster resource
- `manifests/production` production-ready bundle
- `docs` architecture and operational guidance
- `examples/gitops` Argo CD and Flux examples

## Architecture summary

The operator reconciles a single PowerDNSCluster custom resource into a stateless PowerDNS authoritative tier, a PostgreSQL HA backend, backup and restore workflows, DNSSEC key management, and observability/security integrations.

## Production entrypoints

- install the CRD from `config/crd/bases`
- deploy the sample or production manifest from `config/samples` or `manifests/production`
- install the Helm chart from `charts/powerdns-platform`

## Current scope

This scaffold establishes the platform control plane and deployment model. The next implementation layer is the concrete reconciliation of PowerDNS pods, PostgreSQL integration, cert-manager resources, and backup jobs.
