# Mini Study

企业培训学习小程序系统，包含后端（Go + Gin）和微信小程序前端两部分。

[![GitHub](https://img.shields.io/badge/GitHub-mini--study-blue?logo=github)](https://github.com/Rodert/mini-study)

## 📚 项目简介

Mini Study 是一个完整的企业培训学习解决方案，支持：
- 👥 多角色用户管理（管理员、店长、员工）
- 📚 学习内容管理（文档、视频、图文）
- 📝 在线考试系统
- 📊 学习进度跟踪
- 🎯 积分管理
- 🎨 轮播图管理
- 🌱 成长圈（店长发布、管理员审核的公司级动态）
 - 🔔 系统公告管理（首页弹窗 + 后台公告维护）

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

## 🚀 快速开始

### 详细步骤

请查看 **[快速启动指南](./QUICKSTART.md)** 了解详细的初始化步骤。

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

## 📖 文档

- **[快速启动指南](./QUICKSTART.md)** - 详细的初始化步骤和配置说明
- **[后端文档](./mini-study-backend/README.md)** - 后端 API 文档和开发指南
- **[小程序文档](./mini-study-app/README.md)** - 小程序开发文档
- **[微信集成分析](./mini-study-backend/WECHAT_INTEGRATION_ANALYSIS.md)** - 微信小程序集成说明

## 🔑 默认账号

运行数据库迁移后，系统会自动创建以下账号：

| 角色 | 工号 | 密码 | 说明 |
|------|------|------|------|
| 管理员 | `admin` | `admin123456` | 系统管理员，拥有所有权限 |
| 店长 | `manager001` | `123456` | 示例店长账号 |
| 员工 | `employee001` | `123456` | 示例员工账号 |

> ⚠️ **安全提示**: 生产环境请务必修改默认密码！

## 🛠️ 技术栈

### 后端
- **语言**: Go 1.23+
- **框架**: Gin
- **ORM**: GORM
- **数据库**: MySQL
- **认证**: JWT
- **日志**: Zap
- **文档**: Swagger

### 前端
- **平台**: 微信小程序
- **语言**: JavaScript
- **框架**: 原生小程序框架

## 📝 功能特性

### 用户端
- ✅ 用户登录/注册
- ✅ 学习内容浏览（文档、视频、图文）
- ✅ 学习进度记录（文档/图文打开即完成，视频按进度计算）
- ✅ 在线考试
- ✅ 个人中心
- ✅ 成长圈浏览与复制下载（仅展示已审核通过的动态）

### 管理员端
- ✅ 用户管理（员工、店长）
- ✅ 内容管理（创建、编辑、上下架，支持文档/视频/图文三种类型）
  - 学习内容支持草稿 / 已发布 / 已下线三种状态
  - 学员端只会看到已发布内容
- ✅ 考试管理
- ✅ 轮播图管理
 - ✅ 公告管理（维护系统公告，控制首页弹窗内容）
- ✅ 积分管理
- ✅ 成长圈审核后台（按状态和关键字筛选、通过/拒绝/删除动态）

### 店长端
- ✅ 下属员工管理
- ✅ 学习进度查看
- ✅ 考试统计
- ✅ 成长圈发布与自助管理（待审核/已拒绝可删除，审核通过后仅管理员可删）

## 🔧 开发

### 后端开发

```bash
cd mini-study-backend
make run          # 启动服务（Air 热重载）
make build        # 编译项目
make test         # 运行测试
make swagger      # 生成 Swagger 文档
```

详细说明请查看 [后端文档](./mini-study-backend/README.md)

### 前端开发

1. 使用微信开发者工具打开 `mini-study-app` 目录
2. 修改 `services/api.js` 中的 API 地址
3. 编译运行

详细说明请查看 [小程序文档](./mini-study-app/README.md)

## 🚢 部署

### 后端部署

支持多种部署方式：
- Docker 部署
- Docker Compose 部署
- 直接运行二进制文件

详细说明请查看 [后端文档](./mini-study-backend/README.md)

### 前端部署

1. 在微信开发者工具中上传代码
2. 在微信公众平台提交审核
3. 发布上线

## 📄 许可证

本项目采用 MIT 许可证。

## 📜 软件著作权

本项目已申请软件著作权。

## 🔗 相关链接

- **GitHub 仓库**: [https://github.com/Rodert/mini-study](https://github.com/Rodert/mini-study)

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

- 提交 Issue: [https://github.com/Rodert/mini-study/issues](https://github.com/Rodert/mini-study/issues)
- 提交 Pull Request: [https://github.com/Rodert/mini-study/pulls](https://github.com/Rodert/mini-study/pulls)

---

**祝您使用愉快！** 🚀
