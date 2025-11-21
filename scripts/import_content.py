#!/usr/bin/env python3
"""
学习内容批量导入脚本

使用方法:
    python3 import_content.py --host http://localhost:8080 --username admin --password admin123456 --data contents.json

示例数据格式 (contents.json):
    [
        {
            "title": "产品培训视频",
            "type": "video",  # doc 或 video
            "category_id": 1,  # 分类ID，需要先在系统中创建分类
            "file_path": "/uploads/video.mp4",
            "cover_url": "./covers/cover.jpg",  # 支持本地文件路径（自动上传）或 URL 或已上传路径
            "summary": "本视频介绍产品核心功能",
            "visible_roles": "both",  # employee / manager / both
            "status": "published",  # draft / published
            "duration_seconds": 3600  # 视频时长（秒），视频类型必填
        }
    ]
    
封面图支持三种格式:
    1. 本地文件路径: "./covers/cover.jpg" - 脚本会自动上传
    2. 已上传路径: "/uploads/xxx.jpg" - 直接使用
    3. 外部URL: "https://example.com/cover.jpg" - 直接使用
"""

import argparse
import json
import os
import requests
import sys
from typing import Dict, List, Optional


class ContentImporter:
    def __init__(self, base_url: str, username: str, password: str):
        self.base_url = base_url.rstrip('/')
        self.username = username
        self.password = password
        self.token: Optional[str] = None

    def login(self) -> bool:
        """登录并获取 token"""
        url = f"{self.base_url}/api/v1/users/login"
        payload = {
            "work_no": self.username,
            "password": self.password
        }
        try:
            response = requests.post(url, json=payload)
            response.raise_for_status()
            data = response.json()
            if data.get("code") == 200 and data.get("data"):
                self.token = data["data"].get("token")
                print(f"✓ 登录成功: {self.username}")
                return True
            else:
                print(f"✗ 登录失败: {data.get('message', '未知错误')}")
                return False
        except Exception as e:
            print(f"✗ 登录失败: {e}")
            return False

    def list_categories(self) -> List[Dict]:
        """获取所有分类"""
        url = f"{self.base_url}/api/v1/contents/categories"
        headers = {"Authorization": f"Bearer {self.token}"}
        try:
            response = requests.get(url, headers=headers)
            response.raise_for_status()
            data = response.json()
            if data.get("code") == 200:
                return data.get("data", [])
            return []
        except Exception as e:
            print(f"✗ 获取分类列表失败: {e}")
            return []

    def upload_file(self, file_path: str) -> Optional[str]:
        """上传文件并返回存储路径"""
        if not os.path.exists(file_path):
            print(f"  ⚠ 文件不存在: {file_path}，跳过上传")
            return None
        
        url = f"{self.base_url}/api/v1/files/upload"
        headers = {"Authorization": f"Bearer {self.token}"}
        
        try:
            with open(file_path, 'rb') as f:
                files = {'file': (os.path.basename(file_path), f)}
                response = requests.post(url, headers=headers, files=files)
                response.raise_for_status()
                data = response.json()
                if data.get("code") == 200:
                    upload_path = data.get("data", {}).get("path", "")
                    print(f"  ✓ 上传成功: {os.path.basename(file_path)} -> {upload_path}")
                    return upload_path
                else:
                    print(f"  ✗ 上传失败: {data.get('message', '未知错误')}")
                    return None
        except Exception as e:
            print(f"  ✗ 上传失败: {e}")
            return None

    def create_content(self, content_data: Dict) -> bool:
        """创建单个学习内容"""
        # 处理封面图上传
        cover_url = content_data.get("cover_url", "")
        if cover_url:
            # 如果 cover_url 是本地文件路径（不是 http/https 或 /uploads/ 开头），则先上传
            if not cover_url.startswith(("http://", "https://", "/uploads/")):
                upload_path = self.upload_file(cover_url)
                if upload_path:
                    content_data["cover_url"] = upload_path
                else:
                    # 上传失败，可以选择移除 cover_url 或使用原值
                    print(f"  ⚠ 封面图上传失败，将跳过封面图")
                    content_data["cover_url"] = ""
        
        url = f"{self.base_url}/api/v1/admin/contents"
        headers = {
            "Authorization": f"Bearer {self.token}",
            "Content-Type": "application/json"
        }
        try:
            response = requests.post(url, json=content_data, headers=headers)
            response.raise_for_status()
            data = response.json()
            if data.get("code") == 200:
                content_id = data.get("data", {}).get("id", "?")
                print(f"  ✓ 创建成功: {content_data.get('title')} (ID: {content_id})")
                return True
            else:
                print(f"  ✗ 创建失败: {content_data.get('title')} - {data.get('message', '未知错误')}")
                return False
        except Exception as e:
            print(f"  ✗ 创建失败: {content_data.get('title')} - {e}")
            return False

    def import_from_file(self, file_path: str, show_categories: bool = False):
        """从 JSON 文件导入内容"""
        if not self.token:
            print("✗ 请先登录")
            return

        # 显示分类列表（可选）
        if show_categories:
            print("\n可用分类:")
            categories = self.list_categories()
            for cat in categories:
                print(f"  - ID: {cat.get('id')}, 名称: {cat.get('name')}, 角色范围: {cat.get('role_scope')}")

        # 读取并导入数据
        try:
            with open(file_path, 'r', encoding='utf-8') as f:
                contents = json.load(f)
        except Exception as e:
            print(f"✗ 读取文件失败: {e}")
            return

        if not isinstance(contents, list):
            print("✗ 数据格式错误: 根元素必须是数组")
            return

        print(f"\n开始导入 {len(contents)} 条学习内容...\n")
        success_count = 0
        for idx, content in enumerate(contents, 1):
            print(f"[{idx}/{len(contents)}] {content.get('title', '无标题')}")
            if self.create_content(content):
                success_count += 1

        print(f"\n导入完成: 成功 {success_count}/{len(contents)}")


def main():
    parser = argparse.ArgumentParser(description='批量导入学习内容')
    parser.add_argument('--host', default='http://localhost:8080', help='API 服务器地址')
    parser.add_argument('--username', default='admin', help='管理员用户名（工号）')
    parser.add_argument('--password', default='admin123456', help='管理员密码')
    parser.add_argument('--data', required=True, help='JSON 数据文件路径')
    parser.add_argument('--categories', action='store_true', help='显示可用分类列表')

    args = parser.parse_args()

    importer = ContentImporter(args.host, args.username, args.password)

    if not importer.login():
        sys.exit(1)

    importer.import_from_file(args.data, show_categories=args.categories)


if __name__ == '__main__':
    main()

