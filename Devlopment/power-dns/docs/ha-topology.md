# HA Topology

```mermaid
flowchart TB
  Clients[DNS Clients] --> LB[Global Load Balancer]
  LB --> PDNS1[PowerDNS Replica A]
  LB --> PDNS2[PowerDNS Replica B]
  LB --> PDNS3[PowerDNS Replica C]
  PDNS1 --> PGB[PgBouncer]
  PDNS2 --> PGB
  PDNS3 --> PGB
  PGB --> PG1[(PostgreSQL Primary)]
  PG1 --> PG2[(Synchronous Replica)]
  PG1 --> PG3[(Async Read Replica)]
  PG1 --> Backup[(Object Storage Backups)]
```

## Failure handling

- PowerDNS replica failure is handled by Kubernetes rescheduling.
- PostgreSQL primary failure is handled by the HA backend operator.
- backup restore is the last resort for data corruption or region-wide failure.
