# 数据导入脚本使用说明

本目录包含用于批量导入学习内容和考试题的脚本工具。

## 目录结构

```
scripts/
├── import_content.py      # 学习内容导入脚本 (Python, JSON格式)
├── import_content.go      # 学习内容导入脚本 (Go, JSON格式)
├── import_content_csv.py  # 学习内容导入脚本 (Python, CSV格式)
├── import_exam.py         # 考试题导入脚本 (Python, JSON格式)
├── import_exam.go         # 考试题导入脚本 (Go, JSON格式)
├── import_exam_csv.py     # 考试题导入脚本 (Python, CSV格式)
├── example_contents.json  # 学习内容示例数据 (JSON)
├── example_contents.csv   # 学习内容示例数据 (CSV)
├── example_exams.json     # 考试题示例数据 (JSON)
├── example_exams.csv      # 考试题示例数据 (CSV)
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

### 方式一：使用 JSON 格式导入（推荐，结构化数据）

#### 1. 导入学习内容 (JSON)

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

#### 2. 导入考试题 (JSON)

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

### 方式二：使用 CSV 格式导入（推荐，适合 Excel 编辑）

#### 1. 导入学习内容 (CSV)

```bash
# 基本用法
python3 import_content_csv.py --data example_contents.csv

# 指定服务器和账号
python3 import_content_csv.py \
  --host http://localhost:8080 \
  --username admin \
  --password admin123456 \
  --data example_contents.csv

# 显示可用分类列表
python3 import_content_csv.py --data example_contents.csv --categories
```

**CSV 格式要求：**
- 必须包含以下列（列名不区分大小写）：
  - `title`: 内容标题（必填）
  - `type`: 内容类型 `doc` 或 `video`（必填）
  - `category_id`: 分类ID（必填）
  - `file_path`: 文件路径（必填）
  - `cover_url`: 封面图路径（可选，支持本地文件自动上传）
  - `summary`: 内容摘要（可选）
  - `visible_roles`: 可见角色 `employee`/`manager`/`both`（可选）
  - `status`: 状态 `draft`/`published`（可选，默认 `draft`）
  - `duration_seconds`: 视频时长（秒，视频类型必填）

#### 2. 导入考试题 (CSV)

```bash
# 基本用法
python3 import_exam_csv.py --data example_exams.csv

# 指定服务器和账号
python3 import_exam_csv.py \
  --host http://localhost:8080 \
  --username admin \
  --password admin123456 \
  --data example_exams.csv
```

**CSV 格式要求：**
- 考试信息列：
  - `exam_title`: 考试标题（必填）
  - `exam_description`: 考试描述（可选）
  - `exam_status`: 状态 `draft`/`published`/`archived`（可选）
  - `target_role`: 目标角色 `employee`/`manager`/`all`（可选）
  - `time_limit_minutes`: 时间限制（分钟，可选，0表示无限制）
  - `pass_score`: 及格分数（必填）
  
- 题目信息列：
  - `question_type`: 题型 `single`（单选）或 `multiple`（多选）（必填）
  - `question_stem`: 题干（必填）
  - `question_score`: 分值（必填，默认1）
  - `question_analysis`: 解析（可选）
  - `options`: 选项内容，用 `|` 分隔，格式为 `标签:内容:是否正确`
    - 示例1（带标签）：`A:选项A内容:true|B:选项B内容:false`
    - 示例2（无标签，自动生成A/B/C/D...）：`选项A内容:true|选项B内容:false`

**CSV 题目说明：**
- 同一个考试的多道题目可以在多行中，使用相同的 `exam_title`
- 脚本会自动将同一考试的所有题目合并为一个考试

### 方式三：使用 Go 脚本（性能更好，与项目语言一致）

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

### 学习内容格式

#### JSON 格式 (contents.json)

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
- `file_path`: 文件路径（必填），规范格式：
  - **视频/文档文件**：使用 `video/文件名.mp4` 格式（如 `video/basic.mp4`）
  - 支持本地文件自动上传，脚本会自动上传到服务器
  - 也支持已上传路径（如 `/uploads/xxx.mp4`）或外部URL
- `cover_url`: 封面图片路径（可选），规范格式：
  - **封面图片**：使用 `img/文件名.jpg` 格式（如 `img/basic.jpg`）
  - 支持本地文件自动上传，脚本会自动上传到服务器
  - 也支持已上传路径（如 `/uploads/xxx.jpg`）或外部URL
- `summary`: 内容摘要（可选）
- `visible_roles`: 可见角色范围（可选）
  - `employee`: 仅员工可见
  - `manager`: 仅店长可见
  - `both`: 全部可见（默认）
- `status`: 状态（可选）
  - `draft`: 草稿（默认）
  - `published`: 已发布
- `duration_seconds`: 视频时长（秒），视频类型必填，文档类型可以填 0

#### CSV 格式 (contents.csv)

```csv
title,type,category_id,file_path,cover_url,summary,visible_roles,status,duration_seconds
产品培训视频,video,1,video/basic.mp4,img/basic.jpg,本视频介绍产品核心功能,both,published,3600
操作手册,doc,2,video/manual.pdf,img/manual.jpg,详细的产品操作手册,both,published,0
```

**必填列：**
- `title`: 内容标题
- `type`: 内容类型（`doc` 或 `video`）
- `category_id`: 分类ID
- `file_path`: 文件路径（规范：`video/文件名.mp4`，支持本地文件自动上传）

**可选列：**
- `cover_url`: 封面图路径（规范：`img/文件名.jpg`，支持本地文件自动上传）
- `summary`: 内容摘要
- `visible_roles`: 可见角色（`employee`/`manager`/`both`）
- `status`: 状态（`draft`/`published`）
- `duration_seconds`: 视频时长（秒）

**注意事项：**
- 列名不区分大小写，支持空格和下划线
- 自动检测分隔符（逗号或分号）
- 支持 UTF-8 编码，自动处理 BOM

### 考试题格式

#### JSON 格式 (exams.json)

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

**JSON 格式验证：**

```bash
# Python 验证
python3 -m json.tool example_contents.json > /dev/null && echo "格式正确"

# 或使用 jq（如果已安装）
jq . example_contents.json > /dev/null && echo "格式正确"
```

**CSV 格式验证：**

CSV 格式会在导入时自动验证。建议：
- 使用 Excel 或 WPS 编辑 CSV 文件
- 确保列名正确（可参考示例文件）
- 选项字段使用英文引号包裹（如果包含逗号）

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

### 1. 选择合适的数据格式

- **JSON 格式**：适合程序生成、嵌套数据复杂的情况
- **CSV 格式**：适合 Excel 编辑、数据量大、简单表格的情况

### 2. 分批导入

如果数据量很大，建议分批导入：

```bash
# JSON 格式：将大文件拆分成多个小文件
split -l 10 large_contents.json content_batch_

# CSV 格式：可以使用 Excel 分页或手动拆分
# 使用 Python 脚本拆分 CSV
python3 -c "
import csv
with open('large_contents.csv', 'r') as f:
    reader = csv.reader(f)
    header = next(reader)
    for i, row in enumerate(reader):
        if i % 10 == 0:
            out = open(f'content_batch_{i//10}.csv', 'w')
            writer = csv.writer(out)
            writer.writerow(header)
        writer.writerow(row)
"

# 逐个导入
for file in content_batch_*; do
  python3 import_content_csv.py --data "$file"
done
```

### 3. 使用脚本自动化

创建自动化脚本：

```bash
#!/bin/bash
# import_all.sh

BASE_URL="http://localhost:8080"
USERNAME="admin"
PASSWORD="admin123456"

echo "导入学习内容（JSON）..."
python3 import_content.py \
  --host "$BASE_URL" \
  --username "$USERNAME" \
  --password "$PASSWORD" \
  --data contents.json

echo "导入学习内容（CSV）..."
python3 import_content_csv.py \
  --host "$BASE_URL" \
  --username "$USERNAME" \
  --password "$PASSWORD" \
  --data contents.csv

echo "导入考试题（JSON）..."
python3 import_exam.py \
  --host "$BASE_URL" \
  --username "$USERNAME" \
  --password "$PASSWORD" \
  --data exams.json

echo "导入考试题（CSV）..."
python3 import_exam_csv.py \
  --host "$BASE_URL" \
  --username "$USERNAME" \
  --password "$PASSWORD" \
  --data exams.csv

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

## CSV 与 JSON 格式对比

| 特性 | JSON 格式 | CSV 格式 |
|------|----------|---------|
| 可读性 | 中等（结构化） | 高（表格形式） |
| 编辑工具 | 文本编辑器、JSON 工具 | Excel、WPS、文本编辑器 |
| 嵌套数据 | 支持复杂嵌套 | 简单格式（选项用分隔符） |
| 数据量 | 适合任意大小 | 适合大量数据 |
| 程序生成 | 容易 | 容易 |
| 人工编辑 | 较难 | 容易 |
| 导入速度 | 快 | 快 |
| 学习曲线 | 低 | 极低 |

**建议：**
- 少量数据、复杂结构 → 使用 JSON
- 大量数据、表格形式 → 使用 CSV
- 需要 Excel 编辑 → 使用 CSV
- 需要程序批量生成 → 两种都可以

## 更多帮助

- API 文档: `http://localhost:8080/swagger/index.html`
- 后端 README: `mini-study-backend/README.md`
- 示例文件: `scripts/example_*.json` 和 `scripts/example_*.csv`

