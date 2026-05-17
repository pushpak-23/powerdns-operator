# Performance Tuning Recommendations

## PowerDNS

- keep authoritative pods stateless and horizontally scaled
- separate API traffic from query traffic where possible
- raise file descriptor and conntrack limits on worker nodes
- tune cache and thread settings based on QPS and latency targets

## PostgreSQL

- prefer synchronous replication for low-RPO primary clusters
- use PgBouncer to protect PostgreSQL from connection spikes
- isolate WAL and data volumes when the storage backend supports it
- set autovacuum and checkpoint parameters to match zone churn

## Kubernetes

- use pod anti-affinity and topology spread constraints
- reserve CPU for query-serving pods to reduce latency jitter
- avoid overly aggressive HPA scale-down for DNS workloads
