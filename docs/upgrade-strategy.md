# Upgrade Strategy

Use version-aware rolling upgrades in this order:

1. PostgreSQL operator or HA backend if a backend schema migration is required.
2. PgBouncer and connection policy changes.
3. PowerDNS authoritative pods.
4. Operator version.

## Recommended practices

- block upgrades until backup freshness and restore checks pass
- run canary replicas in a limited zone set before broad rollout
- keep API version compatibility gates in the CRD status
- use GitOps drift detection so manual edits are corrected automatically
