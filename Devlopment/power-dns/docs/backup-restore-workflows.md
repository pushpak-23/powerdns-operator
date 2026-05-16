# Backup and Restore Workflows

## Backup

1. Snapshot or dump the PostgreSQL backend on the configured schedule.
2. Archive WAL continuously if PITR is enabled.
3. Verify backup integrity automatically.
4. Record the last verified backup in the cluster status.

## Restore

1. Pause automated writes.
2. Restore the database into a new or recovered cluster.
3. Reconcile PowerDNS against the restored backend.
4. Validate zone counts, SOA serials, and DNSSEC signatures.

## Verification

- run restore drills regularly
- keep a separate validation namespace for backup tests
