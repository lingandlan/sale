## 1. Magic Bytes 校验实现

- [x] 1.1 在 `upload.go` 中添加 `imageMagicSignatures` 签名映射表（JPG/PNG/GIF/WebP/BMP）
- [x] 1.2 实现 `validateImageMagicBytes` 函数：按扩展名匹配文件头签名，WebP 额外校验 WEBP 标记
- [x] 1.3 在 `UploadFile` handler 中扩展名校验后插入 magic bytes 读取 + 校验逻辑，校验通过后 Seek 回文件头

## 2. 验证

- [x] 2.1 编译通过 `go build ./...`
- [x] 2.2 正常图片上传功能不受影响
- [x] 2.3 伪装文件（如 `.jpg` 扩展名的文本文件）被拒绝并返回正确错误信息
