# Mini Study 微信小程序

企业培训学习系统的微信小程序前端，提供用户学习、考试、管理等功能。

[![GitHub](https://img.shields.io/badge/GitHub-mini--study-blue?logo=github)](https://github.com/Rodert/mini-study)

## 📋 目录

- [快速开始](#快速开始)
- [项目结构](#项目结构)
- [功能模块](#功能模块)
- [配置说明](#配置说明)
- [开发指南](#开发指南)
- [常见问题](#常见问题)

## 🚀 快速开始

### 环境要求

- **微信开发者工具**: 最新版本
- **Node.js** (可选): 如果使用 npm 管理依赖

### 配置步骤

#### 1. 配置 Mock 模式或真实 API

编辑 `services/config.js`，配置是否使用 Mock 数据：

```javascript
// 是否使用 Mock 数据（开发环境可以设置为 true，生产环境应设置为 false）
const USE_MOCK = false; // 设置为 true 使用 mock 数据，false 使用真实 API

// API 基础地址配置
const API_BASE_URL = 'http://localhost:8080/api/v1'; // 根据实际环境修改
```

**配置说明：**
- **USE_MOCK = true**: 使用本地 Mock 数据，无需后端服务，适合前端独立开发
- **USE_MOCK = false**: 使用真实 API，需要后端服务运行

> 💡 **提示**: 
> - 开发环境：可以设置 `USE_MOCK = true` 进行前端开发，或 `USE_MOCK = false` 连接后端
> - 生产环境：必须设置 `USE_MOCK = false` 并使用真实后端地址
> - 真机调试：将 `localhost` 改为本机 IP 地址（如 `http://192.168.1.100:8080`）

#### 2. 打开微信开发者工具

1. 打开微信开发者工具
2. 选择"导入项目"
3. 选择 `mini-study-app` 目录
4. 填写 AppID（测试号或正式 AppID）
5. 点击"编译"

#### 3. 配置网络请求域名

**开发环境：**
- 在微信开发者工具中：设置 → 项目设置 → 本地设置 → 勾选"不校验合法域名"

**生产环境：**
- 在微信公众平台 → 开发 → 开发管理 → 开发设置 → 服务器域名中配置
- 需要配置的域名：
  - request 合法域名：后端 API 地址
  - uploadFile 合法域名：后端 API 地址
  - downloadFile 合法域名：后端 API 地址

## 📁 项目结构

```
mini-study-app/
├── app.js                  # 小程序入口文件
├── app.json                # 小程序全局配置
├── app.wxss                # 小程序全局样式
├── pages/                  # 页面目录
│   ├── login/              # 登录页
│   ├── register/           # 注册页
│   ├── home/               # 首页
│   ├── profile/            # 个人中心
│   ├── learning/           # 学习模块
│   │   ├── index/          # 学习首页
│   │   ├── list/           # 学习列表
│   │   └── detail/         # 学习详情（支持文档/视频/图文三种类型）
│   ├── exams/              # 考试模块
│   │   ├── list/           # 考试列表
│   │   ├── detail/         # 考试详情
│   │   └── result/         # 考试结果
│   ├── growth/             # 成长圈模块
│   │   ├── index/          # 成长圈首页（公开动态 + 搜索 + 复制下载）
│   │   └── mine/           # 我的成长圈（店长发布/管理）
│   ├── admin/              # 管理员模块
│   │   ├── employees/      # 员工管理
│   │   ├── managers/       # 店长管理
│   │   ├── contents/       # 内容管理
│   │   ├── exams/          # 考试管理
│   │   ├── banners/        # 轮播图管理
│   │   ├── points/         # 积分管理
│   │   └── growth/         # 成长圈审核后台
│   ├── manager/            # 店长模块
│   │   ├── users/          # 员工管理
│   │   └── progress/       # 学习进度
│   └── webview/            # WebView 页面
├── services/               # 服务层
│   ├── api.js              # API 接口封装
│   └── mockService.js      # Mock 数据服务
├── assets/                 # 静态资源
│   └── icons/              # 图标资源
├── styles/                 # 样式文件
│   ├── layout.wxss         # 布局样式
│   └── theme.wxss          # 主题样式
├── mock/                   # Mock 数据
│   └── mockData.js
├── sitemap.json            # 小程序站点地图
├── project.config.json     # 项目配置
└── project.private.config.json  # 私有配置
```

## 🎯 功能模块

### 用户端功能

#### 1. 登录注册
- **登录**: 使用工号和密码登录
- **注册**: 员工自助注册，支持选择店长
- **Token 管理**: 自动保存和刷新访问令牌

#### 2. 首页
- 轮播图展示
- 学习内容推荐
- 考试入口
- 学习进度概览

#### 3. 学习模块
- **学习列表**: 查看可学习的内容（文档/视频/图文）
- **学习详情**: 
  - 文档阅读
  - 视频播放（支持进度记录）
  - 图文内容（按后端结构化块渲染文本 + 图片）
  - 学习进度显示
- **学习记录**: 查看个人学习历史

#### 4. 考试模块
- **考试列表**: 查看可参加的考试
- **考试详情**: 
  - 题目展示
  - 答题功能
  - 倒计时提醒
- **考试结果**: 查看考试成绩和答案解析

#### 5. 个人中心
- 个人信息展示
- 编辑个人资料
- 学习统计
- 退出登录

#### 6. 成长圈
- **成长圈首页**：
  - 展示所有已审核通过的成长圈动态
  - 支持按关键字搜索内容
  - 每条动态支持“一键复制文案 + 批量保存所有图片到相册”
- **我的成长圈（店长）**：
  - 店长发布成长圈动态（纯文本 + 最多 9 张图片）
  - 查看自己发布的动态列表，按状态和关键字筛选
  - 删除未通过审核的动态，重新发布

### 管理员功能

#### 1. 用户管理
- **员工管理**: 
  - 查看员工列表
  - 添加员工
  - 查看员工详情
  - 调整员工与店长的绑定关系
- **店长管理**: 
  - 添加店长账号
  - 将员工提升为店长

#### 2. 内容管理
- 创建学习内容（文档/视频）
- 编辑内容信息
- 内容上下架
- 内容分类管理

#### 3. 考试管理
- 创建考试
- 编辑考试题目
- 查看考试统计
- 考试结果管理

#### 4. 轮播图管理
- 创建轮播图
- 配置跳转链接
- 设置可见角色
- 设置时间窗口

#### 5. 积分管理
- 查看积分记录
- 按用户查询积分
- 积分统计

#### 6. 成长圈审核
- 查看所有成长圈动态列表
- 按状态（待审核/已通过/已拒绝）和关键字筛选
- 对待审核动态执行“通过/拒绝”操作
- 删除不合规动态

### 店长功能

#### 1. 员工管理
- 查看下属员工列表
- 查看员工学习进度
- 编辑员工信息

#### 2. 学习进度
- 查看下属员工学习统计
- 学习进度分析
- 考试通过率统计
- 关注成长圈发布情况（通过成长圈模块）

## ⚙️ 配置说明

### app.json 配置（节选）

主要配置项：

```json
{
  "pages": [...],           // 页面路径列表
  "window": {               // 窗口配置
    "navigationBarTitleText": "企业培训",
    "navigationBarBackgroundColor": "#ffffff"
  },
  "tabBar": {               // 底部导航栏
    "list": [
      {
        "pagePath": "pages/home/index",
        "text": "首页"
      },
      {
        "pagePath": "pages/growth/index/index",
        "text": "成长圈"
      },
      {
        "pagePath": "pages/profile/index",
        "text": "我的"
      }
    ]
  }
}
```

### API 配置

在 `services/config.js` 中配置：

```javascript
const USE_MOCK = false; // 是否使用 Mock 数据
const API_BASE_URL = 'http://localhost:8080/api/v1'; // API 地址
```

**Mock 模式说明：**
- 当 `USE_MOCK = true` 时，所有 API 调用会自动使用 `mockService.js` 中的 Mock 数据
- Mock 数据位于 `mock/mockData.js`，可以根据需要修改
- Mock 模式下的响应格式已与真实 API 保持一致

### 权限配置

小程序需要以下权限：
- 网络请求权限
- 文件上传权限
- 用户信息权限（如需要）

## 💻 开发指南

### 页面开发

每个页面包含 4 个文件：
- `index.js` - 页面逻辑
- `index.wxml` - 页面结构
- `index.wxss` - 页面样式
- `index.json` - 页面配置

### API 调用

使用封装的 API 服务：

```javascript
const api = require('../../services/api');

// 登录
api.user.login({
  work_no: 'admin',
  password: 'admin123456'
}).then(res => {
  console.log('登录成功', res);
});

// 获取内容列表
api.content.listPublished({
  category_id: 1,
  type: 'document'
}).then(res => {
  console.log('内容列表', res);
});
```

### Token 管理

Token 自动管理：
- 登录成功后自动保存 token
- 请求时自动添加 Authorization 头
- 401 错误时自动跳转登录页

### 文件上传

```javascript
const api = require('../../services/api');

wx.chooseImage({
  success: (res) => {
    const filePath = res.tempFilePaths[0];
    api.file.upload(filePath).then(res => {
      console.log('上传成功', res);
    });
  }
});
```

### 样式规范

- 使用全局样式：`styles/layout.wxss` 和 `styles/theme.wxss`
- 遵循微信小程序样式规范
- 使用 rpx 单位适配不同屏幕

## 🐛 常见问题

### 1. 网络请求失败

**问题**: 无法请求后端 API

**解决方案**:
- 检查后端服务是否启动
- 确认 `api.js` 中的 API 地址正确
- 在微信开发者工具中启用"不校验合法域名"
- 真机调试时使用本机 IP 地址

### 2. Token 过期

**问题**: 提示登录已过期

**解决方案**:
- 系统会自动跳转到登录页
- 重新登录获取新的 token
- 检查后端 JWT 配置是否正确

### 3. 文件上传失败

**问题**: 文件上传报错

**解决方案**:
- 检查文件大小是否超过限制
- 确认后端文件上传接口正常
- 检查网络连接

### 4. 真机调试问题

**问题**: 真机无法连接后端

**解决方案**:
- 确保手机和电脑在同一网络
- 将 API 地址中的 `localhost` 改为本机 IP
- 检查防火墙设置
- 确保后端服务允许外部访问

### 5. 页面跳转问题

**问题**: 页面跳转失败

**解决方案**:
- 检查页面路径是否正确
- 确认页面已在 `app.json` 中注册
- 检查页面文件是否存在

## 📚 相关文档

- [快速启动指南](../QUICKSTART.md) - 详细的初始化步骤
- [后端文档](../mini-study-backend/README.md) - 后端 API 文档
- [微信集成分析](../mini-study-backend/WECHAT_INTEGRATION_ANALYSIS.md) - 微信小程序集成说明

## 🔑 默认账号

| 角色 | 工号 | 密码 | 说明 |
|------|------|------|------|
| 管理员 | `admin` | `admin123456` | 系统管理员 |
| 店长 | `manager001` | `123456` | 示例店长 |
| 员工 | `employee001` | `123456` | 示例员工 |

> ⚠️ 生产环境请务必修改默认密码！

## 📝 开发规范

- 遵循微信小程序开发规范
- 使用统一的 API 调用方式
- 保持代码风格一致
- 添加必要的注释
- 处理异常情况

---

**祝您开发愉快！** 🚀

