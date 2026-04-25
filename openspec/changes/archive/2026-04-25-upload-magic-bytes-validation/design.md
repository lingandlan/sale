## Context

当前 `UploadFile` handler 仅校验文件扩展名（`.jpg`/`.png` 等），未检查文件实际内容。文件保存后由 Gin `r.Static("/uploads", "./uploads")` 直接作为静态资源返回，无执行风险，但在反向代理配置错误（如 Nginx 将 uploads 目录下文件转发给 PHP-FPM）时存在安全隐患。

## Goals / Non-Goals

**Goals:**
- 在服务端校验上传文件的真实格式，拒绝内容与扩展名不匹配的文件
- 覆盖所有已允许的图片类型：JPG、PNG、GIF、WebP、BMP

**Non-Goals:**
- 不做图片内容深度检测（如恶意像素 payload）
- 不增加图片压缩/重编码
- 不改变文件存储路径或命名策略

## Decisions

### 1. Magic bytes 校验 vs 文件 MIME type 检测

**选择**: Magic bytes（文件头签名）校验

**理由**: `http.DetectContentType` 只检测前 512 字节且粒度较粗（无法区分 GIF87a/GIF89a），magic bytes 更精确可控。且不引入新依赖，纯标准库实现。

### 2. 校验时机

**选择**: 在扩展名校验之后、保存文件之前读取前 12 字节校验

**理由**: 先用扩展名做快速过滤，再通过 magic bytes 做内容验证。校验后 `Seek(0,0)` 重置读取位置，确保 `SaveUploadedFile` 写入完整内容。

### 3. WebP 双重校验

**选择**: RIFF 头 + WEBP 标记（offset 8-12）

**理由**: RIFF 是通用容器格式（也用于 WAV/AVI），必须额外校验 WEBP 标记避免误判。

## Risks / Trade-offs

- **误拒风险** → 极低。标准图片文件的 magic bytes 是固定的，不会出现合法图片被拒的情况。
- **性能影响** → 可忽略。只读取 12 字节 + 一次 Seek，无磁盘 IO。
- **绕过可能** → magic bytes 可伪造，但需同时匹配扩展名。对当前威胁模型（防止低级伪装上传）已足够。
