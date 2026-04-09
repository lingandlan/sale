# =============================================================================
# GitHub Actions Secrets 配置
# =============================================================================

# 需要在 GitHub Repository Settings → Secrets 中配置以下 Secrets

## 1. ALIYUN_REGISTRY_USERNAME
- 阿里云 RAM 用户 AccessKey ID
- 用于登录阿里云 ACR

## 2. ALIYUN_REGISTRY_PASSWORD
- 阿里云 RAM 用户 AccessKey Secret
- 用于登录阿里云 ACR

## 3. GITOPS_REPO_TOKEN
- GitHub Personal Access Token
- 用于更新 GitOps 仓库
- 需要 repo 权限

# =============================================================================
# 配置步骤
# =============================================================================

# 1. 创建阿里云 RAM 用户
#    - 登录阿里云 RAM 控制台
#    - 创建用户，授权 ACR 操作权限
#    - 获取 AccessKey ID 和 Secret

# 2. 创建 GitHub Personal Access Token
#    - Settings → Developer settings → Personal access tokens
#    - 生成新 token，勾选 repo 权限

# 3. 配置 GitHub Secrets
#    - Repository → Settings → Secrets and variables → Actions
#    - 添加上述三个 secrets
