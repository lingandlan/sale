## ADDED Requirements

### Requirement: Tag 触发发布
系统 SHALL 在推送 `v*` 格式的 git tag 时自动触发发布流程。

#### Scenario: 推送版本 tag
- **WHEN** 推送 `v1.0.0` 格式的 tag
- **THEN** 自动触发 release workflow

#### Scenario: 非 version tag 不触发
- **WHEN** 推送非 `v*` 开头的 tag
- **THEN** 不触发 release workflow

### Requirement: 自动生成 Changelog
发布流程 SHALL 自动从 git log 生成 changelog。

#### Scenario: Changelog 生成
- **WHEN** release workflow 触发
- **THEN** 收集上一个 tag 到当前 tag 之间的所有 commit，按类型分组生成 changelog

#### Scenario: 首次发布
- **WHEN** 没有上一个 tag
- **THEN** 收集所有 commit 生成 changelog

### Requirement: GitHub Release 创建
发布流程 SHALL 自动创建 GitHub Release，附带 changelog 和构建产物。

#### Scenario: Release 创建成功
- **WHEN** changelog 生成完成
- **THEN** 在 GitHub 上创建 Release，标题为 tag 名称，body 为 changelog 内容
