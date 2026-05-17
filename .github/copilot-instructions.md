# PowerDNS Platform Workspace Instructions

- Use the repository as a monorepo for the operator, Helm chart, GitOps examples, and platform documentation.
- Keep changes focused on the PowerDNS platform domain: operator API, controller logic, Helm templates, and production docs.
- Prefer Go for control-plane code and Helm/Kubernetes YAML for deployment artifacts.
- Treat the operator as the source of truth for lifecycle management and keep the Helm chart declarative and secure by default.
- Preserve compatibility with OpenStack Designate, cert-manager, Prometheus, Grafana, and GitOps workflows.