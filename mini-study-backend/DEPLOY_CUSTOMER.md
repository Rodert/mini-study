# 客户部署说明（mini-study 后端）

本文档说明客户如何使用 `docker-compose.customer.yaml` 启动系统，无需源码。

---

## 一、运行环境要求

- 已安装：
  - Docker 20+ 版本
  - Docker Compose（或内置 `docker compose` 命令）
- 服务器可以访问公网（从阿里云 ACR 拉取镜像）。
- 当前目录下包含：
  - `docker-compose.customer.yaml`

> 若需要完全离线部署，可参考文末“离线部署（可选）”。

---

## 二、首次部署步骤（在线）

1. **准备目录**（在部署目录下执行）：

   ```bash
   mkdir -p storage/uploads
   mkdir -p mysql_data
   mkdir -p redis_data
   ```

2. **启动服务**：

   ```bash
   docker compose -f docker-compose.customer.yaml up -d
   ```

   - 首次启动会自动从阿里云拉取镜像：  
     `crpi-4otucz63tm2q5dhq.cn-beijing.personal.cr.aliyuncs.com/project-shiyu/mini-study-backend:v1.0.0`

3. **查看状态**：

   ```bash
   docker compose -f docker-compose.customer.yaml ps
   ```

4. **查看日志（如需排查问题）**：

   ```bash
   docker compose -f docker-compose.customer.yaml logs -f app
   ```

---

## 三、数据存储与目录说明

- `./storage/uploads` 映射到容器内 `/app/storage/uploads`，用于存放上传文件。  
  删除该目录会导致上传文件丢失。
- `./mysql_data` 映射到 MySQL 数据目录 `/var/lib/mysql`，用于持久化数据库。  
  删除该目录会导致数据库数据丢失。
- `./redis_data` 映射到 Redis 数据目录 `/data`。

> 建议将部署目录及上述数据目录纳入备份策略。

---

## 四、升级版本（在线）

当提供了新版本镜像（例如 `v1.0.1`）后：

1. 编辑 `docker-compose.customer.yaml`，将 `image` 中的标签改为新版本，例如：

   ```yaml
   image: crpi-4otucz63tm2q5dhq.cn-beijing.personal.cr.aliyuncs.com/project-shiyu/mini-study-backend:v1.0.1
   ```

2. 执行：

   ```bash
   docker compose -f docker-compose.customer.yaml pull app migrate
   docker compose -f docker-compose.customer.yaml up -d
   ```

   这样会拉取新镜像并以新版本重启相关服务。

---

## 五、停止与重启

- **停止所有服务**：

  ```bash
  docker compose -f docker-compose.customer.yaml down
  ```

- **重启服务（保持数据不变）**：

  ```bash
  docker compose -f docker-compose.customer.yaml up -d
  ```

---

## 六、数据备份与恢复

> 以下命令示例基于当前 `docker-compose.customer.yaml` 中的默认配置：
> - MySQL 容器服务名为 `mysql`
> - 数据库名称为 `mini_study`
> - MySQL root 密码为 `root`

### 1. 备份 MySQL 数据库

1. 在部署目录下创建备份目录：

   ```bash
   mkdir -p backup
   ```

2. 执行备份命令（在部署目录下）：

   ```bash
   docker compose -f docker-compose.customer.yaml exec -T mysql \
     mysqldump -uroot -proot --databases mini_study > backup/mini_study_$(date +%F).sql
   ```

   - 会在 `backup` 目录下生成一个形如 `mini_study_2025-12-23.sql` 的备份文件。
   - 若实际的数据库名称、账号或密码有调整，请相应修改上述命令中的参数。

### 2. 备份上传文件

上传的文件存放在部署目录下的 `storage/uploads` 目录（已通过 Volume 映射到容器内）。

在部署目录下执行：

```bash
mkdir -p backup

tar czf backup/uploads_$(date +%F).tar.gz storage/uploads
```

将生成形如 `uploads_2025-12-23.tar.gz` 的文件，可用于迁移或归档。

### 3. 恢复数据（仅供参考）

#### 3.1 恢复数据库

在目标环境中，确保已部署好本系统且 MySQL 已可用，然后在部署目录下执行：

```bash
cat backup/mini_study_2025-12-23.sql | \
  docker compose -f docker-compose.customer.yaml exec -T mysql \
  mysql -uroot -proot
```

请根据实际备份文件名和数据库账号信息调整命令。

#### 3.2 恢复上传文件

在目标环境的部署目录下，将备份文件解压到当前目录：

```bash
tar xzf backup/uploads_2025-12-23.tar.gz -C .
```

解压后会恢复 `storage/uploads` 目录中的文件结构。

---

## 七、离线部署（可选说明，供交付方使用）

> 本节主要给系统提供方参考，客户如需完全离线部署，可由提供方先打包镜像，再在客户现场导入。

### 1. 提供方本地打包镜像

在 `mini-study-backend` 目录本地执行：

```bash
docker build -t crpi-4otucz63tm2q5dhq.cn-beijing.personal.cr.aliyuncs.com/project-shiyu/mini-study-backend:v1.0.0 .

docker save -o backend-v1.0.0.tar \
  crpi-4otucz63tm2q5dhq.cn-beijing.personal.cr.aliyuncs.com/project-shiyu/mini-study-backend:v1.0.0
```

如客户无法访问 Docker Hub，也可以一并打包基础镜像：

```bash
docker pull mysql:8 redis:7

docker save -o base-images.tar mysql:8 redis:7
```

将以下文件拷贝到客户服务器：

- `docker-compose.customer.yaml`
- `backend-v1.0.0.tar`
- （可选）`base-images.tar`

### 2. 客户现场导入镜像并启动

在客户服务器上执行：

```bash
# 导入业务镜像
docker load -i backend-v1.0.0.tar

# 如有基础镜像包
# docker load -i base-images.tar

# 启动服务
docker compose -f docker-compose.customer.yaml up -d
```

由于 `docker-compose.customer.yaml` 中的 `image` 名称与打包导入时保持一致，Docker 会直接使用本地镜像，不会访问外网。
