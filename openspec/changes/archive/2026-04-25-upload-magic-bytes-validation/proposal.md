## Why

当前上传接口仅校验文件扩展名，未校验文件实际内容（magic bytes）。攻击者可将 webshell 伪装为 `.jpg` 上传，若服务器配置不当可能导致远程代码执行。需在服务端增加文件头签名校验，确保上传文件的真实格式与扩展名一致。

## What Changes

- 在 `UploadFile` handler 中增加 magic bytes 读取与校验逻辑
- 校验失败时拒绝上传并返回明确错误信息
- 支持 JPG、PNG、GIF、WebP、BMP 五种图片格式的签名匹配

## Capabilities

### New Capabilities
- `upload-validation`: 上传文件内容校验 — 通过 magic bytes 验证文件真实格式，防止伪装上传

### Modified Capabilities

（无 — 现有 spec 无需变更）

## Impact

- **代码**: `backend/internal/handler/upload.go` — 增加约 40 行校验逻辑
- **API**: 上传接口 `/api/upload` 响应不变，新增错误码场景（内容与扩展名不匹配）
- **依赖**: 仅使用标准库 `bytes`，无新外部依赖
