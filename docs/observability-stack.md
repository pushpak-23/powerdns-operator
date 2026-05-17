# Observability Stack

## Metrics

- PowerDNS query rate, latency, response codes, and cache hit ratio
- PostgreSQL replication lag, WAL generation, checkpoint pressure, and connection saturation
- operator reconcile duration, failures, and queue depth

## Logging

- structured JSON logs from PowerDNS and the operator
- centralized collection with Loki or equivalent log backends

## Tracing

- OpenTelemetry traces for API request handling and reconcile workflows

## Alerting

- zone mismatch
- DNSSEC key rollover failure
- replication lag beyond threshold
- backup freshness violation
- API latency SLO breach
