# Recommendations and Operational Notes

## Best PostgreSQL HA architecture

CloudNativePG with synchronous replication, a connection pooler, and PITR backups is the default recommendation for Kubernetes-native environments.

## Best DNS scaling pattern

Use stateless PowerDNS replicas behind service load balancing, keep database access pooled, and scale on query latency and connection pressure.

## Networking considerations

- allow UDP and TCP 53 end-to-end
- isolate API and administrative traffic from public DNS traffic
- ensure conntrack tables and ephemeral ports are sized for high QPS

## DNS caching strategies

- rely on recursive resolvers for caching where possible
- keep authoritative TTLs aligned with update frequency and failover objectives

## Bottlenecks to watch

- PostgreSQL connection saturation
- replica lag after large zone imports
- small node conntrack tables
- excessive API polling from external automation

## Operational runbooks

- backup verify before every major upgrade
- confirm zone counts after failover and restore
- rotate DNSSEC keys on a controlled schedule
- test restoration in a non-production namespace before production recovery
