#!/usr/bin/env python3
"""
考试题批量导入脚本

使用方法:
    python3 import_exam.py --host http://localhost:8080 --username admin --password admin123456 --data exams.json

示例数据格式 (exams.json):
    [
        {
            "title": "产品知识考试",
            "description": "测试对产品知识的掌握程度",
            "status": "published",  # draft / published / archived
            "target_role": "employee",  # employee / manager / all
            "time_limit_minutes": 60,  # 时间限制（分钟），0 表示无限制
            "pass_score": 60,  # 及格分数（必填）
            "questions": [
                {
                    "type": "single",  # single（单选）或 multiple（多选）
                    "stem": "产品的核心功能是什么？",
                    "score": 10,  # 分值
                    "analysis": "产品核心功能包括...",
                    "options": [
                        {
                            "label": "A",  # 选项标签，可选，默认按 A/B/C/D 顺序
                            "content": "功能A",
                            "is_correct": true,  # 是否正确答案
                            "sort_order": 0  # 排序顺序，可选
                        },
                        {
                            "label": "B",
                            "content": "功能B",
                            "is_correct": false
                        }
                    ]
                }
            ]
        }
    ]
"""

import argparse
import json
import requests
import sys
from typing import Dict, List, Optional


class ExamImporter:
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

    def create_exam(self, exam_data: Dict) -> bool:
        """创建单个考试"""
        url = f"{self.base_url}/api/v1/admin/exams"
        headers = {
            "Authorization": f"Bearer {self.token}",
            "Content-Type": "application/json"
        }
        try:
            response = requests.post(url, json=exam_data, headers=headers)
            response.raise_for_status()
            data = response.json()
            if data.get("code") == 200:
                exam_id = data.get("data", {}).get("id", "?")
                question_count = data.get("data", {}).get("question_count", "?")
                total_score = data.get("data", {}).get("total_score", "?")
                print(f"  ✓ 创建成功: {exam_data.get('title')} (ID: {exam_id}, 题目: {question_count}, 总分: {total_score})")
                return True
            else:
                print(f"  ✗ 创建失败: {exam_data.get('title')} - {data.get('message', '未知错误')}")
                return False
        except Exception as e:
            print(f"  ✗ 创建失败: {exam_data.get('title')} - {e}")
            return False

    def import_from_file(self, file_path: str):
        """从 JSON 文件导入考试"""
        if not self.token:
            print("✗ 请先登录")
            return

        # 读取并导入数据
        try:
            with open(file_path, 'r', encoding='utf-8') as f:
                exams = json.load(f)
        except Exception as e:
            print(f"✗ 读取文件失败: {e}")
            return

        if not isinstance(exams, list):
            print("✗ 数据格式错误: 根元素必须是数组")
            return

        print(f"\n开始导入 {len(exams)} 个考试...\n")
        success_count = 0
        for idx, exam in enumerate(exams, 1):
            question_count = len(exam.get('questions', []))
            print(f"[{idx}/{len(exams)}] {exam.get('title', '无标题')} ({question_count} 题)")
            if self.create_exam(exam):
                success_count += 1

        print(f"\n导入完成: 成功 {success_count}/{len(exams)}")


def main():
    parser = argparse.ArgumentParser(description='批量导入考试题')
    parser = argparse.ArgumentParser(description='批量导入考试题')
    parser.add_argument('--host', default='http://localhost:8080', help='API 服务器地址')
    parser.add_argument('--username', default='admin', help='管理员用户名（工号）')
    parser.add_argument('--password', default='admin123456', help='管理员密码')
    parser.add_argument('--data', required=True, help='JSON 数据文件路径')

    args = parser.parse_args()

    importer = ExamImporter(args.host, args.username, args.password)

    if not importer.login():
        sys.exit(1)

    importer.import_from_file(args.data)


if __name__ == '__main__':
    main()

