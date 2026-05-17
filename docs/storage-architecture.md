# Storage Architecture

## Recommended approach

- use the PostgreSQL operator’s native replication and backup model
- prefer SSD-backed ReadWriteOnce volumes for primary and replica instances
- isolate WAL from data when the storage backend supports separate volumes
- store backups in object storage with immutability controls if available

## Recommendations

- CloudNativePG is the preferred default for Kubernetes-native HA.
- Patroni is a strong option when existing operational standards require it.
- Crunchy PostgreSQL is a good fit where enterprise support and ecosystem integration matter.

## Anti-patterns

- do not run PostgreSQL on shared NFS for high-scale DNS workloads
- do not place all replicas in one failure domain
- do not rely on manual failover without automation and status validation
