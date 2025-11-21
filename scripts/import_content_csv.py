#!/usr/bin/env python3
"""
学习内容 CSV 批量导入脚本

使用方法:
    python3 import_content_csv.py --host http://localhost:8080 --username admin --password admin123456 --data contents.csv

CSV 格式说明:
    必须包含以下列（列名不区分大小写）:
    - title: 内容标题
    - type: 内容类型 (doc/video)
    - category_id: 分类ID
    - file_path: 文件路径
    - cover_url: 封面图路径（可选，支持本地文件自动上传）
    - summary: 内容摘要（可选）
    - visible_roles: 可见角色 (employee/manager/both, 可选)
    - status: 状态 (draft/published, 可选，默认draft)
    - duration_seconds: 视频时长（秒，视频类型必填）

示例 CSV 文件:
    title,type,category_id,file_path,cover_url,summary,visible_roles,status,duration_seconds
    产品培训视频,video,1,/uploads/video.mp4,./covers/cover.jpg,本视频介绍产品核心功能,both,published,3600
"""

import argparse
import csv
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

    def normalize_path(self, path: str) -> str:
        """规范化路径格式，统一去除 ./ 和开头空格"""
        if not path:
            return ""
        # 去除前后空格
        path = path.strip()
        # 去除开头的 ./
        if path.startswith("./"):
            path = path[2:]
        # 去除开头的 /
        if path.startswith("/") and not path.startswith("/uploads/"):
            path = path[1:]
        return path
    
    def is_local_file(self, path: str) -> bool:
        """判断是否为本地文件路径"""
        if not path:
            return False
        path = path.strip()
        # 如果是 http/https 或已上传路径，不是本地文件
        if path.startswith(("http://", "https://", "/uploads/")):
            return False
        # 否则可能是本地文件
        return True
    
    def upload_file(self, local_file_path: str, target_dir: str = "") -> Optional[str]:
        """上传文件并返回存储路径
        
        Args:
            local_file_path: 本地文件路径
            target_dir: 目标目录（video 或 img），用于提示
        """
        # 规范化路径
        local_file_path = self.normalize_path(local_file_path)
        
        if not os.path.exists(local_file_path):
            print(f"  ⚠ 文件不存在: {local_file_path}，跳过上传")
            return None
        
        url = f"{self.base_url}/api/v1/files/upload"
        headers = {"Authorization": f"Bearer {self.token}"}
        
        try:
            with open(local_file_path, 'rb') as f:
                files = {'file': (os.path.basename(local_file_path), f)}
                response = requests.post(url, headers=headers, files=files)
                response.raise_for_status()
                data = response.json()
                if data.get("code") == 200:
                    upload_path = data.get("data", {}).get("path", "")
                    dir_name = target_dir if target_dir else "文件"
                    print(f"  ✓ 上传成功 [{dir_name}]: {os.path.basename(local_file_path)} -> {upload_path}")
                    return upload_path
                else:
                    print(f"  ✗ 上传失败: {data.get('message', '未知错误')}")
                    return None
        except Exception as e:
            print(f"  ✗ 上传失败: {e}")
            return None

    def create_content(self, content_data: Dict) -> bool:
        """创建单个学习内容"""
        # 处理文件路径规范化
        file_path = content_data.get("file_path", "")
        if file_path:
            file_path = self.normalize_path(file_path)
            # 如果 file_path 是本地文件，则上传到 video 目录
            if self.is_local_file(file_path):
                upload_path = self.upload_file(file_path, "video")
                if upload_path:
                    content_data["file_path"] = upload_path
                else:
                    print(f"  ✗ 文件上传失败，无法创建内容")
                    return False
            else:
                # 规范化已上传路径或外部URL
                content_data["file_path"] = file_path
        
        # 处理封面图上传
        cover_url = content_data.get("cover_url", "")
        if cover_url:
            cover_url = self.normalize_path(cover_url)
            # 如果 cover_url 是本地文件，则上传到 img 目录
            if self.is_local_file(cover_url):
                upload_path = self.upload_file(cover_url, "img")
                if upload_path:
                    content_data["cover_url"] = upload_path
                else:
                    print(f"  ⚠ 封面图上传失败，将跳过封面图")
                    content_data["cover_url"] = ""
            else:
                # 规范化已上传路径或外部URL
                content_data["cover_url"] = cover_url

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

    def normalize_column_name(self, col_name: str) -> str:
        """标准化列名（去除空格，转小写）"""
        return col_name.strip().lower().replace(' ', '_')

    def parse_csv_row(self, row: Dict, headers: List[str]) -> Optional[Dict]:
        """解析 CSV 行数据为内容对象"""
        normalized_headers = {self.normalize_column_name(h): h for h in headers}
        
        content = {}
        
        # 必填字段
        title = row.get(normalized_headers.get('title', 'title'))
        content_type = row.get(normalized_headers.get('type', 'type'))
        category_id = row.get(normalized_headers.get('category_id', 'category_id'))
        file_path = row.get(normalized_headers.get('file_path', 'file_path'))
        
        if not all([title, content_type, category_id, file_path]):
            print(f"  ✗ 缺少必填字段: title={title}, type={content_type}, category_id={category_id}, file_path={file_path}")
            return None
        
        content['title'] = title.strip()
        content['type'] = content_type.strip().lower()
        if content['type'] not in ['doc', 'video']:
            print(f"  ✗ 类型错误: {content['type']} (必须是 doc 或 video)")
            return None
        
        try:
            content['category_id'] = int(category_id)
        except ValueError:
            print(f"  ✗ category_id 格式错误: {category_id}")
            return None
        
        content['file_path'] = self.normalize_path(file_path)
        
        # 可选字段
        cover_url = row.get(normalized_headers.get('cover_url', 'cover_url'), '').strip()
        if cover_url:
            content['cover_url'] = self.normalize_path(cover_url)
        
        summary = row.get(normalized_headers.get('summary', 'summary'), '').strip()
        if summary:
            content['summary'] = summary
        
        visible_roles = row.get(normalized_headers.get('visible_roles', 'visible_roles'), '').strip().lower()
        if visible_roles:
            content['visible_roles'] = visible_roles
        
        status = row.get(normalized_headers.get('status', 'status'), '').strip().lower()
        if status:
            content['status'] = status
        
        duration_seconds = row.get(normalized_headers.get('duration_seconds', 'duration_seconds'), '').strip()
        if duration_seconds:
            try:
                content['duration_seconds'] = int(duration_seconds)
            except ValueError:
                print(f"  ⚠ duration_seconds 格式错误: {duration_seconds}，将使用 0")
                content['duration_seconds'] = 0
        else:
            # 如果是视频类型但没有提供时长，默认 0
            content['duration_seconds'] = 0
        
        return content

    def import_from_csv(self, csv_path: str, show_categories: bool = False):
        """从 CSV 文件导入内容"""
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
            with open(csv_path, 'r', encoding='utf-8-sig') as f:  # utf-8-sig 处理 BOM
                # 检测分隔符
                sample = f.read(1024)
                f.seek(0)
                sniffer = csv.Sniffer()
                delimiter = sniffer.sniff(sample).delimiter
                
                reader = csv.DictReader(f, delimiter=delimiter)
                headers = reader.fieldnames or []
                
                if not headers:
                    print("✗ CSV 文件没有表头")
                    return
                
                print(f"\n检测到列: {', '.join(headers)}\n")
                
                contents = []
                for idx, row in enumerate(reader, 1):
                    content = self.parse_csv_row(row, headers)
                    if content:
                        contents.append(content)
                    else:
                        print(f"  ⚠ 跳过第 {idx + 1} 行（表头为第 1 行）")
                
                if not contents:
                    print("✗ 没有有效的数据行")
                    return
                
                print(f"\n开始导入 {len(contents)} 条学习内容...\n")
                success_count = 0
                for idx, content in enumerate(contents, 1):
                    print(f"[{idx}/{len(contents)}] {content.get('title', '无标题')}")
                    if self.create_content(content):
                        success_count += 1
                
                print(f"\n导入完成: 成功 {success_count}/{len(contents)}")
                
        except FileNotFoundError:
            print(f"✗ 文件不存在: {csv_path}")
        except Exception as e:
            print(f"✗ 读取 CSV 文件失败: {e}")


def main():
    parser = argparse.ArgumentParser(description='批量导入学习内容（CSV 格式）')
    parser.add_argument('--host', default='http://localhost:8080', help='API 服务器地址')
    parser.add_argument('--username', default='admin', help='管理员用户名（工号）')
    parser.add_argument('--password', default='admin123456', help='管理员密码')
    parser.add_argument('--data', required=True, help='CSV 数据文件路径')
    parser.add_argument('--categories', action='store_true', help='显示可用分类列表')

    args = parser.parse_args()

    importer = ContentImporter(args.host, args.username, args.password)

    if not importer.login():
        sys.exit(1)

    importer.import_from_csv(args.data, show_categories=args.categories)


if __name__ == '__main__':
    main()

