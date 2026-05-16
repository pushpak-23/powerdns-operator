# Disaster Recovery Guide

The platform is designed around three recovery layers:

- PostgreSQL PITR restore to recover authoritative zone state
- PowerDNS pod replacement for stateless compute failure
- zone consistency validation after restore or failover

## Recovery workflows

- verify the last successful backup before each restore
- restore the database into a clean target namespace or a new cluster
- reconcile PowerDNS into the restored backend
- compare zone counts, SOA serials, and DNSSEC metadata

## Split-brain prevention

- use a single write endpoint behind PgBouncer
- keep synchronous replication for the primary failure domain when low RPO is required
- gate failover on health checks and WAL continuity
