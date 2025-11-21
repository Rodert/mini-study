# 数据导入脚本使用说明

本目录包含用于批量导入学习内容和考试题的脚本工具。

## 目录结构

```
scripts/
├── import_content.py      # 学习内容导入脚本 (Python)
├── import_content.go      # 学习内容导入脚本 (Go)
├── import_exam.py         # 考试题导入脚本 (Python)
├── import_exam.go         # 考试题导入脚本 (Go)
├── example_contents.json  # 学习内容示例数据
├── example_exams.json     # 考试题示例数据
└── README.md             # 本说明文件
```

## 前置要求

### 使用 Python 脚本
- Python 3.6+
- 安装依赖：`pip install requests`

### 使用 Go 脚本
- Go 1.16+
- 无需额外依赖（使用标准库）

### 系统要求
- 后端服务已启动（默认: `http://localhost:8080`）
- 管理员账号已创建（默认: 工号 `admin` / 密码 `admin123456`）
- 内容分类已创建（导入学习内容前需要）

## 使用方法

### 方式一：使用 Python 脚本（推荐，更易用）

#### 1. 导入学习内容

```bash
# 基本用法
python3 import_content.py --data example_contents.json

# 指定服务器和账号
python3 import_content.py \
  --host http://localhost:8080 \
  --username admin \
  --password admin123456 \
  --data example_contents.json

# 显示可用分类列表
python3 import_content.py --data example_contents.json --categories
```

#### 2. 导入考试题

```bash
# 基本用法
python3 import_exam.py --data example_exams.json

# 指定服务器和账号
python3 import_exam.py \
  --host http://localhost:8080 \
  --username admin \
  --password admin123456 \
  --data example_exams.json
```

### 方式二：使用 Go 脚本（性能更好，与项目语言一致）

#### 1. 编译 Go 脚本（首次使用）

```bash
# 编译学习内容导入脚本
go build -o import_content import_content.go

# 编译考试题导入脚本
go build -o import_exam import_exam.go
```

#### 2. 使用编译后的脚本

```bash
# 导入学习内容
./import_content --data example_contents.json

# 导入考试题
./import_exam --data example_exams.json

# 指定服务器和账号
./import_content \
  --host http://localhost:8080 \
  --username admin \
  --password admin123456 \
  --data example_contents.json
```

#### 3. 或直接运行（不编译）

```bash
# 导入学习内容
go run import_content.go --data example_contents.json

# 导入考试题
go run import_exam.go --data example_exams.json
```

## 数据格式说明

### 学习内容格式 (contents.json)

```json
[
  {
    "title": "内容标题",              // 必填
    "type": "video",                  // 必填: "doc" 或 "video"
    "category_id": 1,                 // 必填: 分类ID
    "file_path": "/uploads/video.mp4", // 必填: 文件路径
    "cover_url": "https://...",       // 可选: 封面图片URL
    "summary": "内容摘要",            // 可选
    "visible_roles": "both",          // 可选: "employee" / "manager" / "both" (默认继承分类)
    "status": "published",            // 可选: "draft" / "published" (默认: "draft")
    "duration_seconds": 3600          // 视频类型必填: 时长（秒）
  }
]
```

**字段说明：**
- `title`: 内容标题，最大 255 字符
- `type`: 内容类型，`doc`（文档）或 `video`（视频）
- `category_id`: 分类ID，需要先在管理后台创建分类
- `file_path`: 文件存储路径（相对路径或绝对路径）
- `cover_url`: 封面图片URL（可选）
- `summary`: 内容摘要（可选）
- `visible_roles`: 可见角色范围（可选）
  - `employee`: 仅员工可见
  - `manager`: 仅店长可见
  - `both`: 全部可见（默认）
- `status`: 状态（可选）
  - `draft`: 草稿（默认）
  - `published`: 已发布
- `duration_seconds`: 视频时长（秒），视频类型必填，文档类型可以填 0

### 考试题格式 (exams.json)

```json
[
  {
    "title": "考试标题",                    // 必填
    "description": "考试描述",              // 可选
    "status": "published",                 // 可选: "draft" / "published" / "archived"
    "target_role": "employee",             // 可选: "employee" / "manager" / "all"
    "time_limit_minutes": 60,              // 可选: 时间限制（分钟），0 表示无限制
    "pass_score": 60,                      // 必填: 及格分数
    "questions": [                         // 必填: 题目列表，至少 1 题
      {
        "type": "single",                  // 必填: "single"（单选）或 "multiple"（多选）
        "stem": "题干",                    // 必填: 题目内容
        "score": 10,                       // 必填: 分值（>= 1）
        "analysis": "解析",                // 可选: 题目解析
        "options": [                       // 必填: 选项列表，至少 2 个
          {
            "label": "A",                  // 可选: 选项标签（默认 A/B/C/D...）
            "content": "选项内容",         // 必填: 选项文本
            "is_correct": true,            // 必填: 是否正确答案
            "sort_order": 0                // 可选: 排序顺序
          }
        ]
      }
    ]
  }
]
```

**字段说明：**

**考试级别：**
- `title`: 考试标题，必填
- `description`: 考试描述（可选）
- `status`: 状态（可选）
  - `draft`: 草稿（默认）
  - `published`: 已发布
  - `archived`: 已归档
- `target_role`: 目标角色（可选）
  - `employee`: 员工（默认）
  - `manager`: 店长
  - `all`: 全部
- `time_limit_minutes`: 时间限制（分钟），0 表示无限制（可选）
- `pass_score`: 及格分数，必填，不能高于总分

**题目级别：**
- `type`: 题型，必填
  - `single`: 单选题（只能选一个答案）
  - `multiple`: 多选题（可以选多个答案）
- `stem`: 题干内容，必填
- `score`: 分值，必填，最小值为 1
- `analysis`: 题目解析（可选）
- `options`: 选项列表，必填，至少 2 个选项

**选项级别：**
- `label`: 选项标签（A/B/C/D...），可选，不提供时自动生成
- `content`: 选项内容，必填
- `is_correct`: 是否正确答案，必填
- `sort_order`: 排序顺序，可选，默认按数组顺序

**注意事项：**
1. 每道题至少需要一个正确答案
2. 单选题必须设置唯一正确答案
3. 多选题可以设置多个正确答案
4. 每题至少需要 2 个选项
5. 及格分数不能高于总分（系统会自动计算总分）

## 准备数据文件

### 1. 准备分类信息

在导入学习内容之前，需要先在管理后台创建内容分类，或查看现有分类：

```bash
python3 import_content.py --data example_contents.json --categories
```

### 2. 准备文件路径

确保文件路径正确：
- 如果是相对路径（如 `/uploads/video.mp4`），文件需要已上传到服务器的 `storage/uploads` 目录
- 如果是绝对路径，确保路径可访问
- 封面图片可以使用外部 URL

### 3. 验证数据格式

导入前建议先检查 JSON 格式是否正确：

```bash
# Python 验证
python3 -m json.tool example_contents.json > /dev/null && echo "格式正确"

# 或使用 jq（如果已安装）
jq . example_contents.json > /dev/null && echo "格式正确"
```

## 错误处理

### 常见错误及解决方案

1. **登录失败**
   - 检查用户名和密码是否正确
   - 确认服务器地址是否正确
   - 确认服务器是否已启动

2. **分类不存在**
   - 使用 `--categories` 参数查看可用分类
   - 在管理后台先创建分类
   - 确保 `category_id` 正确

3. **文件路径错误**
   - 检查 `file_path` 是否正确
   - 确保文件已上传到服务器
   - 检查文件权限

4. **数据格式错误**
   - 检查 JSON 格式是否正确
   - 确认必填字段是否都已填写
   - 检查字段值是否符合要求（如 `type` 只能是 "doc" 或 "video"）

5. **视频时长未设置**
   - 视频类型的 `duration_seconds` 必须大于 0

## 批量导入技巧

### 1. 分批导入

如果数据量很大，建议分批导入：

```bash
# 将大文件拆分成多个小文件
split -l 10 large_contents.json content_batch_

# 逐个导入
for file in content_batch_*; do
  python3 import_content.py --data "$file"
done
```

### 2. 使用脚本自动化

创建自动化脚本：

```bash
#!/bin/bash
# import_all.sh

BASE_URL="http://localhost:8080"
USERNAME="admin"
PASSWORD="admin123456"

echo "导入学习内容..."
python3 import_content.py \
  --host "$BASE_URL" \
  --username "$USERNAME" \
  --password "$PASSWORD" \
  --data contents.json

echo "导入考试题..."
python3 import_exam.py \
  --host "$BASE_URL" \
  --username "$USERNAME" \
  --password "$PASSWORD" \
  --data exams.json

echo "全部完成！"
```

## 其他导入方式

### 使用 Go 脚本（适合 Go 开发者）

如果需要 Go 版本的导入脚本，可以参考 API 文档自行实现，或使用 `curl` 命令：

```bash
# 登录获取 token
TOKEN=$(curl -s -X POST http://localhost:8080/api/v1/users/login \
  -H "Content-Type: application/json" \
  -d '{"work_no":"admin","password":"admin123456"}' \
  | jq -r '.data.token')

# 导入内容
curl -X POST http://localhost:8080/api/v1/admin/contents \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d @example_contents.json
```

### 使用 Postman/Insomnia

1. 导入 API 文档（Swagger JSON）
2. 创建 Collection
3. 使用 Collection Runner 批量导入

## 更多帮助

- API 文档: `http://localhost:8080/swagger/index.html`
- 后端 README: `mini-study-backend/README.md`

