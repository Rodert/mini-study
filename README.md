# mini-study

企业培训学习小程序

包含后端（Go + Gin）和微信小程序前端两部分。

## 📚 文档

- **[快速启动指南](./QUICKSTART.md)** - 快速初始化并启动项目
- [后端文档](./mini-study-backend/README.md) - 后端 API 和开发文档
- [微信集成分析](./mini-study-backend/WECHAT_INTEGRATION_ANALYSIS.md) - 微信小程序集成说明

## 🚀 快速开始

查看 [快速启动指南](./QUICKSTART.md) 了解详细的初始化步骤。

### 简要步骤

1. **准备数据库**: 创建 MySQL 数据库或使用 Docker Compose
2. **启动后端**: 
   ```bash
   cd mini-study-backend
   go mod tidy
   go run scripts/migrate.go  # 初始化数据库
   make run                    # 启动服务
   ```
3. **配置前端**: 修改 `mini-study-app/services/api.js` 中的 API 地址
4. **打开小程序**: 使用微信开发者工具打开 `mini-study-app` 目录

## 📁 项目结构

```
mini-study/
├── mini-study-backend/    # Go 后端服务
│   ├── cmd/server/        # 应用入口
│   ├── internal/          # 内部代码
│   ├── configs/           # 配置文件
│   └── scripts/           # 脚本文件
└── mini-study-app/        # 微信小程序前端
    ├── pages/             # 页面
    ├── services/          # API 服务
    └── assets/            # 静态资源
```

## 🔑 默认账号

- 管理员: `admin` / `admin123456`
- 店长: `manager001` / `manager123456`
- 员工: `employee001` / `employee123456`

> ⚠️ 生产环境请务必修改默认密码！




