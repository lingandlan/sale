# =============================================================================
# GitOps Repository Structure
# =============================================================================

# This directory contains K8s manifests for ArgoCD deployment
# Repository: github.com/your-org/marketplace-gitops

# Environments
clusters/
├── dev/                  # Development environment
│   ├── backend.yaml       # Backend K8s manifests
│   └── namespace.yaml     # Namespace definition
│
├── staging/              # Staging environment
│   ├── backend.yaml
│   └── namespace.yaml
│
└── prod/                # Production environment
    ├── backend.yaml
    └── namespace.yaml
