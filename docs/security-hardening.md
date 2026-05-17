# Security Hardening Guide

- enforce TLS for all service-to-service traffic
- enable mTLS for the DNS API and internal control-plane calls
- run all containers as non-root with read-only filesystems where possible
- apply Pod Security Standards and OPA/Gatekeeper policies
- restrict traffic using namespace and pod-level network policies
- rotate API keys, TSIG secrets, and DNSSEC keys on a schedule
- use cert-manager for short-lived certificates and Vault if centralized secret control is required
