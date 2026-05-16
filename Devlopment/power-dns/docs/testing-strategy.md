# Testing Strategy

## Unit tests

- reconcile state transitions
- status update logic
- config rendering helpers

## Integration tests

- PostgreSQL HA failover
- DNS zone create/update/delete cycles
- DNSSEC key rollover
- backup and restore verification

## End-to-end tests

- API auth
- Designate zone synchronization
- multi-AZ replica churn
- zero-downtime upgrade flow

## Chaos and resilience

- kill PowerDNS pods during query load
- fail PostgreSQL primary during write bursts
- simulate network partitions between DNS and database tiers
