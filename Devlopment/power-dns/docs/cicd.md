# CI/CD Pipeline

## Pipeline stages

1. format and lint the Go controller
2. validate Helm templates and YAML schemas
3. run unit and integration tests
4. build and sign container images
5. publish Helm chart artifacts
6. deploy to staging with GitOps
7. run smoke and backup-restore tests

## Release controls

- use provenance for chart and image artifacts
- promote by digest, not by mutable tags
- require signed commits or protected branches where possible
