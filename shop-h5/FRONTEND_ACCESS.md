# 前端访问说明

## 开发环境访问地址

**前端服务**:
- 本地访问: http://localhost:3000
- 网络访问: http://本机IP:3000

**后端API**:
- 本地访问: http://localhost:8080
- API文档: http://localhost:8080/swagger（如果启用）

## 启动步骤

### 1. 启动后端（已完成✅）
```bash
cd /Users/zhangdaodong/code/sale/backend
go run ./cmd/server/main.go
```
后端已在端口 8080 运行

### 2. 启动前端（依赖安装中...）
```bash
cd /Users/zhangdaodong/code/sale/shop-h5
npm install          # 首次需要安装依赖
npm run dev          # 启动开发服务器
```

前端将在端口 3000 启动

## 测试账号

```
手机号: 13800000000
密码: Test123456
角色: 超级管理员
```

## 前端页面

### 已完成页面
- ✅ `/pages/login/index` - 登录页面
- ✅ `/pages/dashboard/index` - 首页仪表盘

### 页面路由
```
/pages/login/index      → 登录页
/pages/dashboard/index  → 首页仪表盘
```

## API 配置

前端已配置代理到后端：
```javascript
// vite.config.ts
proxy: {
  '/api': {
    target: 'http://localhost:8080',
    changeOrigin: true
  }
}
```

前端调用 `/api/v1/*` 会自动转发到后端

## 浏览器访问

启动完成后，在浏览器打开：
```
http://localhost:3000
```

### 开发工具
- **HBuilderX**: 推荐用于UniApp开发
- **浏览器**: Chrome/Edge 用于H5调试
- **微信开发者工具**: 用于小程序调试

## 当前状态

✅ 后端服务 - 运行中 (端口 8080)
✅ 数据库 - MySQL + Redis 运行中
⏳ 前端服务 - 依赖安装中...

**预计启动时间**: 2-3分钟（依赖安装完成）

---
**生成时间**: 2026-04-10
**项目路径**: `/Users/zhangdaodong/code/sale/shop-h5`
