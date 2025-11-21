# 快速启动指南

本文档将帮助您快速初始化并启动 mini-study 项目。

## 📋 环境要求

### 后端环境
- **Go**: 1.23.0 或更高版本
- **MySQL**: 5.7+ 或 8.0+
- **Air** (可选): 用于热重载开发，推荐安装
  ```bash
  go install github.com/air-verse/air@latest
  ```

### 前端环境
- **微信开发者工具**: 最新版本
- **Node.js** (可选): 如果使用 npm 管理依赖

## 🚀 快速启动步骤

### 1. 克隆项目（如果还没有）

```bash
git clone <repository-url>
cd mini-study
```

### 2. 数据库准备

#### 方式一：使用 Docker Compose（推荐）

```bash
cd mini-study-backend
docker-compose up -d mysql
```

这将启动 MySQL 容器，默认配置：
- 端口: `3306`
- 用户名: `root`
- 密码: `root`
- 数据库: `mini_study`

#### 方式二：使用本地 MySQL

1. 创建数据库：
```sql
CREATE DATABASE mini_study CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

2. 确保 MySQL 服务已启动

### 3. 后端配置

#### 3.1 修改配置文件

编辑 `mini-study-backend/configs/config.yaml`，确保数据库连接信息正确：

```yaml
database:
  driver: mysql
  dsn: root:root123456@tcp(127.0.0.1:3306)/mini_study?charset=utf8mb4&parseTime=True&loc=Local
```

> ⚠️ **重要**: 请根据您的实际 MySQL 配置修改 `dsn` 中的用户名、密码和端口。

#### 3.2 安装依赖

```bash
cd mini-study-backend
go mod tidy
```

#### 3.3 初始化数据库

运行迁移脚本，自动创建表结构和初始化数据：

```bash
go run scripts/migrate.go
```

这将：
- ✅ 自动创建所有数据表
- ✅ 创建默认管理员账号（工号: `admin` / 密码: `admin123456`）
- ✅ 初始化默认分类和示例数据

#### 3.4 启动后端服务

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

**方式三：使用 Docker**

```bash
make compose
# 或
docker-compose up -d
```

后端服务默认运行在：`http://localhost:8080`

#### 3.5 验证后端启动

- 健康检查: `http://localhost:8080/system/health`
- Swagger 文档: `http://localhost:8080/swagger/index.html` (如果已启用)

### 4. 前端配置

#### 4.1 修改 API 地址

编辑 `mini-study-app/services/api.js`，修改 API 基础地址：

```javascript
const API_BASE_URL = 'http://localhost:8080/api/v1'; // 根据实际环境修改
```

> 💡 **提示**: 
> - 开发环境：使用 `http://localhost:8080`
> - 生产环境：修改为实际的后端服务器地址
> - 如果使用真机调试，需要将 `localhost` 改为本机 IP 地址（如 `http://192.168.1.100:8080`）

#### 4.2 打开微信开发者工具

1. 打开微信开发者工具
2. 选择"导入项目"
3. 选择 `mini-study-app` 目录
4. 填写 AppID（测试号或正式 AppID）
5. 点击"编译"

#### 4.3 配置网络请求域名

在微信公众平台配置服务器域名：
- 开发环境：在微信开发者工具中，设置 → 项目设置 → 本地设置 → 勾选"不校验合法域名"
- 生产环境：在微信公众平台 → 开发 → 开发管理 → 开发设置 → 服务器域名中配置

## 📝 默认账号信息

运行数据库迁移后，系统会自动创建以下账号：

| 角色 | 工号 | 密码 | 说明 |
|------|------|------|------|
| 管理员 | `admin` | `admin123456` | 系统管理员，拥有所有权限 |
| 店长 | `manager001` | `manager123456` | 示例店长账号 |
| 员工 | `employee001` | `employee123456` | 示例员工账号 |

> ⚠️ **安全提示**: 生产环境请务必修改默认密码！

## 🔧 常用命令

### 后端

```bash
# 启动服务（Air 热重载）
make run

# 编译项目
make build

# 运行测试
make test

# 生成 Swagger 文档
make swagger

# Docker 构建
make docker

# Docker Compose 启动
make compose
```

### 数据库迁移

```bash
# 运行迁移脚本
go run scripts/migrate.go
```

## 🌍 环境变量

通过设置 `APP_ENV` 环境变量切换配置：

```bash
# 使用开发环境配置
export APP_ENV=dev
make run

# 使用生产环境配置
export APP_ENV=prod
make run
```

配置文件位置：
- `configs/config.yaml` - 默认配置
- `configs/config.dev.yaml` - 开发环境
- `configs/config.prod.yaml` - 生产环境

## 🐛 常见问题

### 1. 数据库连接失败

**问题**: `dial tcp 127.0.0.1:3306: connect: connection refused`

**解决方案**:
- 检查 MySQL 服务是否启动
- 确认 `config.yaml` 中的数据库连接信息正确
- 检查防火墙设置

### 2. 端口被占用

**问题**: `bind: address already in use`

**解决方案**:
- 修改 `configs/config.yaml` 中的 `server.port`
- 或关闭占用端口的进程

### 3. 微信小程序无法请求后端

**问题**: 网络请求失败

**解决方案**:
- 确保后端服务已启动
- 检查 `api.js` 中的 API 地址是否正确
- 在微信开发者工具中启用"不校验合法域名"选项
- 真机调试时，确保手机和电脑在同一网络，并使用本机 IP 地址

### 4. 依赖安装失败

**问题**: `go mod tidy` 失败

**解决方案**:
- 检查网络连接
- 设置 Go 代理（如需要）:
  ```bash
  go env -w GOPROXY=https://goproxy.cn,direct
  ```

### 5. Air 命令未找到

**问题**: `air: command not found`

**解决方案**:
```bash
go install github.com/air-verse/air@latest
# 确保 $GOPATH/bin 或 $GOBIN 在 PATH 中
```

## 📚 更多信息

- 后端详细文档: `mini-study-backend/README.md`
- API 文档: 启动后端后访问 `http://localhost:8080/swagger/index.html`
- 微信集成分析: `mini-study-backend/WECHAT_INTEGRATION_ANALYSIS.md`

## 🎯 下一步

1. ✅ 后端服务已启动
2. ✅ 前端小程序已打开
3. ✅ 使用默认管理员账号登录测试
4. 📖 查看 Swagger API 文档了解接口详情
5. 🔧 根据实际需求修改配置和代码

---

**祝您开发愉快！** 🚀

