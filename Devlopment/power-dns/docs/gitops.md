# GitOps Deployment Model

- store the operator chart and cluster manifests in Git
- use Argo CD or Flux to reconcile namespaces, CRDs, and the PowerDNSCluster resource
- keep environment overlays for dev, staging, and production
- separate the operator release cadence from the cluster instance cadence
