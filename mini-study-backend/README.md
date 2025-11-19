# Mini Study Backend

基于 Go + Gin 的后端脚手架，内置配置管理、JWT 鉴权、审计日志以及 Air / Swagger / Docker / Makefile 等常用工具，可直接用于业务开发。

## 快速开始

```bash
go mod tidy          # 首次拉依赖，需本地可访问互联网
make run             # 推荐使用 Air 热重载（需提前安装），或执行 go run ./cmd/server
```

> 默认使用本地 MySQL，请确保 `configs/config.yaml` 中的 DSN（root:root123456@tcp(127.0.0.1:3306)/mini_study）与你环境一致；系统以「工号（work_no）」作为唯一登录凭证。

## 常用脚本

- `scripts/migrate.go`：调用 GORM Migrator 自动建表
- `scripts/start.sh`：Air 启动脚本
- `scripts/test.sh`：运行单元测试

## 目录结构概览

- `cmd/server`：应用入口
- `configs/`：多环境配置（local/dev/prod）
- `internal/bootstrap`：配置、日志、数据库、路由初始化
- `internal/middleware`：JWT、CORS、请求日志、Recovery、Validator 等中间件
- `internal/model`：GORM 数据模型（用户、审计）
- `internal/repository`：数据访问层
- `internal/service`：业务逻辑与 JWT Token 服务
- `internal/handler`：HTTP 接口
- `internal/router`：路由分组与 Swagger 挂载
- `internal/utils`：Hash/JWT/文件上传/统一响应等通用工具
- `storage/uploads`：本地上传目录（使用 .gitkeep 保持空目录）

## 工具与辅助文件

- `.air.toml`：Air 热重载配置
- `Makefile`：`run`、`build`、`test`、`swagger`、`docker` 等快捷命令
- `Dockerfile` & `docker-compose.yaml`：容器化部署示例

## 环境变量

通过 `APP_ENV` 切换配置（`local` / `dev` / `prod`），具体 DSN、端口、跨域等可在 `configs/config*.yaml` 中自行调整。建议在实际部署前修改 JWT Secret、数据库账号等敏感信息。

## API 接口概览

> 所有需要登录的接口都必须在请求头中携带 `Authorization: Bearer <access_token>`。管理员接口会在服务端再次校验角色。

### 用户与认证

| 方法 | 路径 | 说明 | 鉴权 |
| --- | --- | --- | --- |
| POST | `/api/v1/users/register` | 员工自助注册，支持同时提交多个店长工号 | 否 |
| POST | `/api/v1/users/login` | 使用工号+密码登录，返回访问令牌与刷新令牌 | 否 |
| POST | `/api/v1/users/token/refresh` | 通过刷新令牌获取新的访问令牌 | 否 |
| GET | `/api/v1/users/managers` | 获取所有可选店长列表（注册前查询） | 否 |
| PATCH | `/api/v1/users/me/profile` | 修改个人姓名、手机号 | 是 |

### 管理员-用户管理

| 方法 | 路径 | 说明 |
| --- | --- | --- |
| POST | `/api/v1/admin/managers` | 创建新的店长账号 |
| POST | `/api/v1/admin/users/:id/promote-manager` | 将员工升为店长 |
| PUT | `/api/v1/admin/users/:id/managers` | 调整员工与店长的绑定关系 |

### 内容与分类

| 方法 | 路径 | 说明 | 鉴权 |
| --- | --- | --- | --- |
| GET | `/api/v1/contents/categories` | 获取当前用户可见的分类列表 | 是 |
| GET | `/api/v1/contents` | 查询已发布内容（可通过分类、类型筛选） | 是 |
| GET | `/api/v1/contents/:id` | 查看内容详情 | 是 |
| GET | `/api/v1/admin/contents` | 管理员查询内容列表（支持状态过滤） | 是（管理员） |
| POST | `/api/v1/admin/contents` | 管理员创建内容（文档/视频） | 是（管理员） |
| PUT | `/api/v1/admin/contents/:id` | 管理员更新内容（含上下架） | 是（管理员） |

> 文档/视频上传走 `/api/v1/files/upload`，返回的磁盘路径填写在内容 `file_path` 字段，视频需额外提供 `duration_seconds`。

### 学习记录

| 方法 | 路径 | 说明 |
| --- | --- | --- |
| POST | `/api/v1/learning` | 上报学习进度（视频播放位置不可回退） |
| GET | `/api/v1/learning/:content_id` | 查看某个内容的学习进度 |
| GET | `/api/v1/learning` | 查看当前用户全部学习记录 |

请求示例：

```jsonc
POST /api/v1/learning
{
  "content_id": 12,
  "video_position": 120  // 单位：秒
}
```

返回示例（部分字段）：

```jsonc
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

### 文件与系统
### 轮播图 Banner

| 方法 | 路径 | 说明 | 鉴权 |
| --- | --- | --- | --- |
| GET | `/api/v1/banners` | 登录后按角色拉取当前有效的轮播图（支持时间窗口、排序） | 是 |
| GET | `/api/v1/admin/banners` | 管理员查看全部轮播图，可按状态过滤 | 是（管理员） |
| POST | `/api/v1/admin/banners` | 管理员创建轮播图，配置图片、跳转链接、可见角色、时间窗口 | 是（管理员） |
| PUT | `/api/v1/admin/banners/:id` | 管理员更新轮播图信息或上下线 | 是（管理员） |

> 轮播图仅支持 H5 外链跳转，如需跳小程序内部页面可在 H5 中自行跳转。


| 方法 | 路径 | 说明 |
| --- | --- | --- |
| POST | `/api/v1/files/upload` | 上传文件（表单字段名 `file`，默认限制 `max_size_mb`） |
| GET | `/system/health` | 健康检查 |
| GET | `/system/version` | 版本信息 |

## Swagger 文档

## 数据初始化

运行 `go run scripts/migrate.go` 时会自动：

- 自动迁移最新的 GORM 模型（包括用户、内容、学习记录、轮播图等表）
- 确保默认管理员账号（工号 `admin` / 密码 `admin123456`）
- 初始化固定的员工/店长分类
- 写入示例轮播图（可在管理后台修改或删除）

项目已集成 `swag`，可通过 `make swagger` 生成最新文档，然后访问 `/swagger/index.html` 查看接口详情（仅在 `config.yaml` 中开启 `swagger.enabled: true` 时生效）。如需更新注释，请在 handler 层补充 Swagger 注解。
