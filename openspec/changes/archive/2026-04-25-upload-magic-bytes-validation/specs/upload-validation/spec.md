## ADDED Requirements

### Requirement: 上传文件内容校验
系统 SHALL 在保存上传文件前，读取文件头部 magic bytes 并与扩展名对应的签名进行匹配校验。校验不通过时 MUST 拒绝上传并返回错误提示。

#### Scenario: 正常 JPG 图片上传
- **WHEN** 用户上传一个扩展名为 `.jpg` 且文件头以 `FF D8 FF` 开头的文件
- **THEN** 系统 ALLOW 上传，返回图片 URL

#### Scenario: 正常 PNG 图片上传
- **WHEN** 用户上传一个扩展名为 `.png` 且文件头以 `89 50 4E 47 0D 0A 1A 0A` 开头的文件
- **THEN** 系统 ALLOW 上传，返回图片 URL

#### Scenario: WebP 文件双重校验
- **WHEN** 用户上传一个扩展名为 `.webp` 且文件头为 `RIFF....WEBP` 的文件
- **THEN** 系统 ALLOW 上传，返回图片 URL

#### Scenario: 伪装文件被拒绝
- **WHEN** 用户上传一个扩展名为 `.jpg` 但文件头非 `FF D8 FF` 开头的文件（如 PHP webshell）
- **THEN** 系统 MUST 拒绝上传，返回错误信息"文件内容与扩展名不匹配"

#### Scenario: 文件内容过短被拒绝
- **WHEN** 上传文件内容不足 4 字节
- **THEN** 系统 MUST 拒绝上传，返回错误信息"文件内容无效"

### Requirement: 支持的图片格式签名
系统 SHALL 支持以下图片格式的 magic bytes 签名校验：
- JPG: `FF D8 FF`
- PNG: `89 50 4E 47 0D 0A 1A 0A`
- GIF87a: `47 49 46 38 37 61`
- GIF89a: `47 49 46 38 39 61`
- WebP: RIFF + WEBP（offset 8-12）
- BMP: `42 4D`

#### Scenario: 所有支持的格式均可通过校验
- **WHEN** 用户分别上传 `.jpg`、`.jpeg`、`.png`、`.gif`、`.webp`、`.bmp` 格式的合法图片
- **THEN** 系统均 ALLOW 上传并返回图片 URL
