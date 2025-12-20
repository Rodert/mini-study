# Docker Compose 启动与数据库更新说明

本项目后端通常通过 `docker compose` 启动，数据库使用容器 + 数据卷的方式持久化。

本文档说明两种常见场景：

- 如何在**保留已有数据**的前提下，让新表 / 新字段生效（增量迁移）。
- 如何**重置数据库**，从头重新初始化（仅限本地开发环境）。

> 说明：
> - 建表与字段更新由 `internal/bootstrap/db.go` 中的 GORM `AutoMigrate` 完成。
> - 每次后端容器启动时，都会执行 `AutoMigrate`，自动创建缺失的表/字段，不会删除已有表和数据。

---

## 一、增量更新：已有数据，只新增表/字段

**使用场景**：

- 之前已经用 `docker compose` 启动过后端，数据库里有正式/测试数据；
- 现在后端代码新增了 model（比如 `Notice` 表）或给已有表加了字段；
- 希望**保留旧数据，只让新表/字段自动创建出来**。

### 步骤 1：拉取最新代码

在宿主机项目根目录执行：

```bash
git pull
```

### 步骤 2：重新构建后端镜像（保留数据库容器和数据卷）

```bash
# 在项目根目录（有 docker-compose.yaml 的位置）
docker compose build          # 或 docker-compose build
# 如只需构建后端服务，可以：
# docker compose build backend
```

### 步骤 3：重新启动服务

```bash
docker compose up -d          # 或 docker-compose up -d
# 如只想重启后端：
# docker compose up -d backend
```

> 只要**不删除数据库容器和数据卷**，原有数据都会保留；
> 新增的表和字段会在后端容器启动时通过 `AutoMigrate` 自动创建。

### 步骤 4：（可选）检查日志确认迁移成功

```bash
docker compose logs backend
```

正常情况下，你会看到类似 `auto migrate` / `database connected` 的日志输出。

### 步骤 5：（可选）重新执行初始化脚本

`scripts/migrate.go` 负责创建默认账号、默认分类、默认轮播图等，它是**幂等**的，多次执行是安全的：

```bash
docker compose exec backend go run scripts/migrate.go
```

如果数据库里已经有默认数据，脚本会跳过创建操作。

---

## 二、重建数据库：清空数据重新初始化

**⚠️ 高危操作，仅建议在本地开发 / 测试环境使用。**

**使用场景**：

- 本地测试数据已经不可用，想完全清空数据库，重新初始化；
- 不再需要当前数据库里的任何数据。

### 步骤 1：停止并删除容器和数据卷

```bash
docker compose down -v        # 或 docker-compose down -v
```

- `down` 会停止并删除 compose 管理的容器；
- `-v` 会一并删除相关数据卷，相当于**清空数据库**。

### 步骤 2：重新启动服务（会自动建表）

```bash
docker compose up -d          # 或 docker-compose up -d
```

- 后端启动时会再次执行 `AutoMigrate`，按最新的 model 定义重新建表。

### 步骤 3：执行迁移脚本初始化默认数据

```bash
docker compose exec backend go run scripts/migrate.go
```

脚本会自动：

- 创建默认管理员 / 店长 / 员工账号；
- 初始化默认学习分类；
- 初始化默认轮播图等示例数据。

---

## 三、常见问题

- **Q：增量更新时会不会影响现有数据？**  
  A：不会。`AutoMigrate` 只会创建缺失的表、字段和索引，不会删除表、字段或已有记录。

- **Q：生产环境能否使用 `down -v`？**  
  A：严格禁止。生产环境请仅使用“增量更新”的方式，通过重新构建镜像 + 重启服务让新结构生效，不要删除数据库数据卷。

- **Q：需要手动执行 SQL 吗？**  
  A：正常情况下不需要，建表/加字段统一由 `AutoMigrate` 负责。如果确有特殊结构变更（如字段重命名、数据迁移），请额外编写专用脚本并在运维流程中单独执行。
