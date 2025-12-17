# Mini Study Backend

基于 Go + Gin 的企业培训学习系统后端服务，提供完整的用户管理、内容管理、学习记录、考试系统等功能。

[![GitHub](https://img.shields.io/badge/GitHub-mini--study-blue?logo=github)](https://github.com/Rodert/mini-study)

## 📋 目录

- [快速开始](#快速开始)
- [项目结构](#项目结构)
- [环境配置](#环境配置)
- [API 接口](#api-接口)
- [开发工具](#开发工具)
- [部署说明](#部署说明)

## 🚀 快速开始

### 环境要求

- **Go**: 1.23.0 或更高版本
- **MySQL**: 5.7+ 或 8.0+
- **Air** (可选): 用于热重载开发

### 安装依赖

```bash
go mod tidy
```

### 配置数据库

编辑 `configs/config.yaml`，修改数据库连接信息：

```yaml
database:
  driver: mysql
  dsn: root:root123456@tcp(127.0.0.1:3306)/mini_study?charset=utf8mb4&parseTime=True&loc=Local
```

> ⚠️ 请根据实际环境修改 DSN 中的用户名、密码和端口

### 初始化数据库

运行迁移脚本，自动创建表结构和初始化数据：

```bash
go run scripts/migrate.go
```

这将自动：
- ✅ 创建所有数据表
- ✅ 创建默认管理员账号（工号: `admin` / 密码: `admin123456`）
- ✅ 初始化默认分类和示例数据

### 启动服务

**方式一：使用 Air 热重载（推荐开发环境）**

```bash
make run
# 或
air
```

**方式二：直接运行**

```bash
go run ./cmd/server/main.go
```

**方式三：使用 Docker Compose**

```bash
make compose
# 或
docker-compose up -d
```

服务默认运行在：`http://localhost:8080`

### 验证启动

- 健康检查: `http://localhost:8080/system/health`
- Swagger 文档: `http://localhost:8080/swagger/index.html` (需在配置中启用)

## 📁 项目结构

```
mini-study-backend/
├── cmd/server/              # 应用入口
│   └── main.go
├── configs/                 # 配置文件
│   ├── config.yaml          # 默认配置
│   ├── config.dev.yaml      # 开发环境
│   └── config.prod.yaml     # 生产环境
├── internal/
│   ├── bootstrap/           # 初始化模块
│   │   ├── config.go        # 配置加载
│   │   ├── db.go            # 数据库初始化
│   │   ├── logger.go        # 日志初始化
│   │   ├── middleware.go    # 中间件注册
│   │   └── router.go        # 路由初始化
│   ├── handler/             # HTTP 接口层
│   │   ├── user_handler.go
│   │   ├── content_handler.go
│   │   ├── exam_handler.go
│   │   ├── learning_handler.go
│   │   ├── banner_handler.go
│   │   └── ...
│   ├── service/             # 业务逻辑层
│   │   ├── user_service.go
│   │   ├── content_service.go
│   │   ├── exam_service.go
│   │   └── ...
│   ├── repository/          # 数据访问层
│   │   ├── user_repo.go
│   │   ├── content_repo.go
│   │   └── ...
│   ├── model/               # 数据模型
│   │   ├── user.go
│   │   ├── content.go
│   │   ├── exam.go
│   │   └── ...
│   ├── dto/                 # 数据传输对象
│   │   ├── user_dto.go
│   │   ├── content_dto.go
│   │   └── ...
│   ├── middleware/          # 中间件
│   │   ├── jwt.go           # JWT 鉴权
│   │   ├── cors.go          # 跨域处理
│   │   ├── logger.go        # 请求日志
│   │   ├── recovery.go      # 错误恢复
│   │   └── validator.go     # 参数验证
│   ├── router/              # 路由定义
│   │   ├── api.go           # API 路由
│   │   ├── system.go        # 系统路由
│   │   └── swagger.go       # Swagger 路由
│   ├── utils/               # 工具函数
│   │   ├── jwt.go           # JWT 工具
│   │   ├── hash.go          # 密码加密
│   │   ├── file.go          # 文件处理
│   │   └── response.go      # 统一响应
│   └── docs/                # Swagger 文档
├── scripts/                 # 脚本文件
│   ├── migrate.go           # 数据库迁移
│   ├── start.sh             # 启动脚本
│   └── test.sh              # 测试脚本
├── storage/                 # 存储目录
│   └── uploads/             # 上传文件
├── test/                    # 测试文件
├── Dockerfile               # Docker 镜像
├── docker-compose.yaml      # Docker Compose 配置
├── Makefile                 # 构建脚本
├── .air.toml                # Air 配置
└── go.mod                   # Go 模块定义
```

## ⚙️ 环境配置

### 环境变量

通过 `APP_ENV` 环境变量切换配置：

```bash
# 使用开发环境配置
export APP_ENV=dev
make run

# 使用生产环境配置
export APP_ENV=prod
make run
```

### 配置文件说明

配置文件支持多环境：
- `configs/config.yaml` - 默认配置（本地开发）
- `configs/config.dev.yaml` - 开发环境
- `configs/config.prod.yaml` - 生产环境

主要配置项：

```yaml
server:
  port: 8080                    # 服务端口
  mode: debug                    # 运行模式: debug/release

database:
  driver: mysql                  # 数据库驱动
  dsn: root:password@tcp(...)   # 数据库连接字符串

jwt:
  secret: your-secret-key        # JWT 密钥（生产环境请修改）
  access_token_ttl: 24h         # 访问令牌有效期
  refresh_token_ttl: 168h       # 刷新令牌有效期

cors:
  allowed_origins: ["*"]         # 允许的跨域来源
  allowed_methods: ["*"]        # 允许的 HTTP 方法

swagger:
  enabled: true                  # 是否启用 Swagger
```

> ⚠️ **安全提示**: 生产环境请务必修改 JWT Secret、数据库密码等敏感信息！

## 📡 API 接口

> 所有需要登录的接口都必须在请求头中携带 `Authorization: Bearer <access_token>`

### 用户与认证

| 方法 | 路径 | 说明 | 鉴权 |
| --- | --- | --- | --- |
| POST | `/api/v1/users/register` | 员工自助注册，支持同时提交多个店长工号 | 否 |
| POST | `/api/v1/users/login` | 使用工号+密码登录，返回访问令牌与刷新令牌 | 否 |
| POST | `/api/v1/users/token/refresh` | 通过刷新令牌获取新的访问令牌 | 否 |
| GET | `/api/v1/users/me` | 获取当前用户信息 | 是 |
| GET | `/api/v1/users/managers` | 获取所有可选店长列表（注册前查询） | 否 |
| PATCH | `/api/v1/users/me/profile` | 修改个人姓名、手机号 | 是 |

### 管理员-用户管理

| 方法 | 路径 | 说明 | 鉴权 |
| --- | --- | --- | --- |
| GET | `/api/v1/admin/users` | 查询用户列表 | 管理员 |
| GET | `/api/v1/admin/users/:id` | 查询单个用户详情 | 管理员 |
| POST | `/api/v1/admin/managers` | 创建新的店长账号 | 管理员 |
| POST | `/api/v1/admin/employees` | 创建新的员工账号 | 管理员 |
| POST | `/api/v1/admin/users/:id/promote-manager` | 将员工升为店长 | 管理员 |
| PUT | `/api/v1/admin/users/:id/managers` | 调整员工与店长的绑定关系 | 管理员 |

### 内容与分类

| 方法 | 路径 | 说明 | 鉴权 |
| --- | --- | --- | --- |
| GET | `/api/v1/contents/categories` | 获取当前用户可见的分类列表 | 是 |
| GET | `/api/v1/contents` | 查询已发布内容（可通过分类、类型筛选，类型支持 doc/video/article） | 是 |
| GET | `/api/v1/contents/:id` | 查看内容详情 | 是 |
| GET | `/api/v1/admin/contents` | 管理员查询内容列表（支持状态过滤） | 管理员 |
| POST | `/api/v1/admin/contents` | 管理员创建内容（文档/视频/图文） | 管理员 |
| PUT | `/api/v1/admin/contents/:id` | 管理员更新内容（含上下架） | 管理员 |

> 内容类型 `type` 支持：`doc`(文档) / `video`(视频) / `article`(图文)。
> - 文档/视频上传走 `/api/v1/files/upload`，返回的磁盘路径填写在内容 `file_path` 字段，视频需额外提供 `duration_seconds`。
> - 图文内容通过请求体中的 `article_blocks` 字段传输结构化文本/图片块，后端以 JSON 串存储在 `BodyBlocksJSON` 字段。

### 学习记录

| 方法 | 路径 | 说明 | 鉴权 |
| --- | --- | --- | --- |
| POST | `/api/v1/learning` | 上报学习进度（视频播放位置不可回退） | 是 |
| GET | `/api/v1/learning/:content_id` | 查看某个内容的学习进度 | 是 |
| GET | `/api/v1/learning` | 查看当前用户全部学习记录 | 是 |

**请求示例：**

```json
POST /api/v1/learning
{
  "content_id": 12,
  "video_position": 120  // 单位：秒
}
```

**返回示例：**

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "content_id": 12,
    "video_position": 180,
    "duration_seconds": 3600,
    "progress": 5,
    "status": "in_progress"
  }
}
```

### 考试系统

| 方法 | 路径 | 说明 | 鉴权 |
| --- | --- | --- | --- |
| GET | `/api/v1/exams` | 获取可用的考试列表 | 是 |
| GET | `/api/v1/exams/:id` | 获取考试详情 | 是 |
| POST | `/api/v1/exams/:id/submit` | 提交考试答案 | 是 |
| GET | `/api/v1/exams/my/results` | 获取我的考试结果 | 是 |
| GET | `/api/v1/admin/exams` | 管理员查询考试列表 | 管理员 |
| POST | `/api/v1/admin/exams` | 管理员创建考试 | 管理员 |
| PUT | `/api/v1/admin/exams/:id` | 管理员更新考试 | 管理员 |

### 轮播图 Banner

| 方法 | 路径 | 说明 | 鉴权 |
| --- | --- | --- | --- |
| GET | `/api/v1/banners` | 登录后按角色拉取当前有效的轮播图（支持时间窗口、排序） | 是 |
| GET | `/api/v1/admin/banners` | 管理员查看全部轮播图，可按状态过滤 | 管理员 |
| POST | `/api/v1/admin/banners` | 管理员创建轮播图，配置图片、跳转链接、可见角色、时间窗口 | 管理员 |
| PUT | `/api/v1/admin/banners/:id` | 管理员更新轮播图信息或上下线 | 管理员 |

> 轮播图仅支持 H5 外链跳转，如需跳小程序内部页面可在 H5 中自行跳转。

### 成长圈（Growth Circle）

成长圈是公司级的动态流功能，由店长发布、管理员审核，通过后所有角色可见。

**用户端接口：**

| 方法 | 路径 | 说明 | 鉴权 |
| --- | --- | --- | --- |
| GET | `/api/v1/growth` | 查询已审核通过的成长圈动态列表，可按关键字搜索内容 | 是 |
| GET | `/api/v1/growth/mine` | 查询当前登录用户发布的成长圈动态，可按状态/关键字筛选 | 是 |
| POST | `/api/v1/growth` | 店长发布成长圈动态（纯文本 + 多图） | 店长 |
| DELETE | `/api/v1/growth/:id` | 删除成长圈动态：店长可删自己未通过的动态，管理员可删任意动态 | 店长/管理员 |

**管理员审核接口：**

| 方法 | 路径 | 说明 | 鉴权 |
| --- | --- | --- | --- |
| GET | `/api/v1/admin/growth` | 管理员查询成长圈动态列表，可按状态/关键字筛选 | 管理员 |
| POST | `/api/v1/admin/growth/:id/approve` | 管理员审核通过指定动态（状态置为 approved） | 管理员 |
| POST | `/api/v1/admin/growth/:id/reject` | 管理员拒绝指定动态（状态置为 rejected） | 管理员 |

### 积分管理

| 方法 | 路径 | 说明 | 鉴权 |
| --- | --- | --- | --- |
| GET | `/api/v1/admin/points` | 管理员查询积分记录列表 | 管理员 |
| GET | `/api/v1/admin/users/:id/points` | 管理员查询指定用户的积分记录 | 管理员 |

### 文件与系统

| 方法 | 路径 | 说明 | 鉴权 |
| --- | --- | --- | --- |
| POST | `/api/v1/files/upload` | 上传文件（表单字段名 `file`，默认限制 `max_size_mb`） | 是 |
| GET | `/system/health` | 健康检查 | 否 |
| GET | `/system/version` | 版本信息 | 否 |

## 🛠️ 开发工具

### Makefile 命令

```bash
make run          # 启动服务（使用 Air 热重载）
make build        # 编译项目
make test         # 运行单元测试
make swagger      # 生成 Swagger 文档
make docker       # 构建 Docker 镜像
make compose      # 使用 Docker Compose 启动
```

### 数据库迁移

```bash
# 运行迁移脚本
go run scripts/migrate.go
```

### Swagger 文档

项目已集成 `swag`，可通过以下命令生成最新文档：

```bash
make swagger
```

然后访问 `http://localhost:8080/swagger/index.html` 查看接口详情。

> 仅在 `config.yaml` 中开启 `swagger.enabled: true` 时生效

### Air 热重载

安装 Air：

```bash
go install github.com/air-verse/air@latest
```

使用 Air 启动服务（自动热重载）：

```bash
make run
# 或
air
```

配置文件：`.air.toml`

## 🚢 部署说明

### Docker 部署

**构建镜像：**

```bash
make docker
# 或
docker build -t mini-study-backend .
```

**使用 Docker Compose：**

```bash
make compose
# 或
docker-compose up -d
```

### 生产环境配置

1. 修改 `configs/config.prod.yaml` 中的配置
2. 设置环境变量 `APP_ENV=prod`
3. 确保数据库连接信息正确
4. 修改 JWT Secret 和数据库密码
5. 配置反向代理（如 Nginx）
6. 配置 HTTPS

### 默认账号

运行数据库迁移后，系统会自动创建以下账号：

| 角色 | 工号 | 密码 | 说明 |
|------|------|------|------|
| 管理员 | `admin` | `admin123456` | 系统管理员，拥有所有权限 |
| 店长 | `manager001` | `123456` | 示例店长账号 |
| 员工 | `employee001` | `123456` | 示例员工账号 |

> ⚠️ **安全提示**: 生产环境请务必修改默认密码！

## 📚 相关文档

- [快速启动指南](../QUICKSTART.md) - 详细的初始化步骤
- [微信集成分析](./WECHAT_INTEGRATION_ANALYSIS.md) - 微信小程序集成说明
- [项目总览](../README.md) - 项目整体介绍

## 🔧 技术栈

- **框架**: Gin
- **ORM**: GORM
- **数据库**: MySQL
- **认证**: JWT
- **日志**: Zap
- **配置**: Viper
- **文档**: Swagger
- **验证**: go-playground/validator

## 📝 开发规范

- 遵循 Go 代码规范
- 使用统一的错误处理
- 使用统一的响应格式
- 所有接口需要添加 Swagger 注释
- 数据库操作统一使用 Repository 层
- 业务逻辑统一放在 Service 层

---

**祝您开发愉快！** 🚀
