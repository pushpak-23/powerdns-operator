# OpenStack Designate Integration Guide

PowerDNS is a strong backend option for Designate when the platform needs DNSSEC, AXFR/IXFR, and API-driven zone automation.

## Integration pattern

- use PowerDNS authoritative API as the write path for zone and record changes
- synchronize Designate zones into PowerDNS through a backend driver or service adapter
- use TSIG for transfer trust where secondary DNS or external replication is involved
- map Keystone authentication to tenant-scoped zone permissions

## Production notes

- isolate tenants by zone ownership and API token boundaries
- enable reverse DNS for PTR automation
- validate DNSSEC signing and key rollovers during zone sync
- use region-aware endpoints for multi-region control planes

## Operational flow

1. Designate accepts a zone mutation request.
2. The backend adapter writes the change to PowerDNS.
3. The operator validates zone state and surfaces status back through conditions.
4. DNSSEC metadata and transfer ACLs are reconciled as part of the zone spec.
