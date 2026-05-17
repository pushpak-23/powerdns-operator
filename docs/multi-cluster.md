# Multi-Cluster Deployment Design

- run one operator-managed DNS control plane per region or failure domain
- keep GitOps overlays per cluster and per environment
- use global traffic steering only for client discovery endpoints, not for write paths
- replicate zone data through controlled sync and restore workflows rather than ad hoc writes

## Recommended pattern

Use a primary regional writer with read-only or delayed replicas in other regions, then expose region-local authoritative services for latency-sensitive traffic.
