#!/usr/bin/env python3
"""
考试题 CSV 批量导入脚本

使用方法:
    python3 import_exam_csv.py --host http://localhost:8080 --username admin --password admin123456 --data exams.csv

CSV 格式说明:
    考试信息列:
    - exam_title: 考试标题
    - exam_description: 考试描述（可选）
    - exam_status: 状态 (draft/published/archived, 可选)
    - target_role: 目标角色 (employee/manager/all, 可选)
    - time_limit_minutes: 时间限制（分钟，可选，0表示无限制）
    - pass_score: 及格分数（必填）
    
    题目信息列:
    - question_type: 题型 (single/multiple)
    - question_stem: 题干
    - question_score: 分值
    - question_analysis: 解析（可选）
    - options: 选项内容，用 | 分隔，每个选项格式为 "标签:内容:是否正确"
       例如: "A:选项A内容:true|B:选项B内容:false|C:选项C内容:true"
       或者简化为: "选项A内容:true|选项B内容:false" (标签自动生成A/B/C/D...)

示例:
    exam_title,exam_description,exam_status,target_role,time_limit_minutes,pass_score,question_type,question_stem,question_score,question_analysis,options
    产品知识考试,测试对产品知识的掌握程度,published,employee,30,60,single,产品的核心功能是什么？,10,产品的核心功能包括...,A:用户管理和数据分析:true|B:仅用户管理:false|C:仅数据分析:false|D:仅报告生成:false
"""

import argparse
import csv
import json
import requests
import sys
from typing import Dict, List, Optional, Tuple


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

    def parse_options(self, options_str: str) -> List[Dict]:
        """解析选项字符串
        格式: "A:内容:true|B:内容:false" 或 "内容:true|内容:false"
        """
        if not options_str or not options_str.strip():
            return []
        
        options = []
        parts = options_str.split('|')
        
        for idx, part in enumerate(parts):
            part = part.strip()
            if not part:
                continue
            
            # 解析选项：可能有标签或没有标签
            # 格式1: "A:内容:true" (有标签)
            # 格式2: "内容:true" (无标签，自动生成)
            segments = part.split(':')
            
            if len(segments) == 3:
                # 有标签格式
                label, content, is_correct_str = segments
                label = label.strip().upper()
                content = content.strip()
                is_correct = is_correct_str.strip().lower() in ('true', '1', 'yes', 'correct', '正确')
            elif len(segments) == 2:
                # 无标签格式，自动生成 A/B/C/D...
                label = chr(65 + idx)  # A, B, C, D...
                content, is_correct_str = segments
                content = content.strip()
                is_correct = is_correct_str.strip().lower() in ('true', '1', 'yes', 'correct', '正确')
            else:
                print(f"  ⚠ 选项格式错误，跳过: {part}")
                continue
            
            if not content:
                continue
            
            options.append({
                "label": label,
                "content": content,
                "is_correct": is_correct,
                "sort_order": idx
            })
        
        return options

    def parse_question(self, row: Dict, headers: List[str]) -> Optional[Dict]:
        """解析题目数据"""
        normalized_headers = {h.strip().lower().replace(' ', '_'): h for h in headers}
        
        question_type = row.get(normalized_headers.get('question_type', 'question_type'), '').strip().lower()
        if question_type not in ['single', 'multiple']:
            return None
        
        question_stem = row.get(normalized_headers.get('question_stem', 'question_stem'), '').strip()
        if not question_stem:
            return None
        
        question_score = row.get(normalized_headers.get('question_score', 'question_score'), '').strip()
        try:
            score = int(question_score) if question_score else 1
        except ValueError:
            score = 1
        
        question_analysis = row.get(normalized_headers.get('question_analysis', 'question_analysis'), '').strip()
        
        options_str = row.get(normalized_headers.get('options', 'options'), '').strip()
        options = self.parse_options(options_str)
        
        if len(options) < 2:
            print(f"  ⚠ 题目选项不足（至少需要2个）: {question_stem[:30]}...")
            return None
        
        # 检查是否有正确答案
        has_correct = any(opt['is_correct'] for opt in options)
        if not has_correct:
            print(f"  ⚠ 题目没有正确答案: {question_stem[:30]}...")
            return None
        
        question = {
            "type": question_type,
            "stem": question_stem,
            "score": score,
            "options": options
        }
        
        if question_analysis:
            question["analysis"] = question_analysis
        
        return question

    def parse_exam_info(self, row: Dict, headers: List[str]) -> Optional[Dict]:
        """解析考试信息"""
        normalized_headers = {h.strip().lower().replace(' ', '_'): h for h in headers}
        
        exam_title = row.get(normalized_headers.get('exam_title', 'exam_title'), '').strip()
        if not exam_title:
            return None
        
        exam_description = row.get(normalized_headers.get('exam_description', 'exam_description'), '').strip()
        exam_status = row.get(normalized_headers.get('exam_status', 'exam_status'), '').strip().lower()
        target_role = row.get(normalized_headers.get('target_role', 'target_role'), '').strip().lower()
        
        time_limit_minutes = row.get(normalized_headers.get('time_limit_minutes', 'time_limit_minutes'), '').strip()
        try:
            time_limit = int(time_limit_minutes) if time_limit_minutes else 0
        except ValueError:
            time_limit = 0
        
        pass_score = row.get(normalized_headers.get('pass_score', 'pass_score'), '').strip()
        try:
            pass_score_int = int(pass_score) if pass_score else 60
        except ValueError:
            print(f"  ✗ pass_score 格式错误: {pass_score}")
            return None
        
        exam_info = {
            "title": exam_title,
            "pass_score": pass_score_int,
            "time_limit_minutes": time_limit
        }
        
        if exam_description:
            exam_info["description"] = exam_description
        
        if exam_status:
            exam_info["status"] = exam_status
        
        if target_role:
            exam_info["target_role"] = target_role
        
        return exam_info

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

    def import_from_csv(self, csv_path: str):
        """从 CSV 文件导入考试"""
        if not self.token:
            print("✗ 请先登录")
            return

        try:
            with open(csv_path, 'r', encoding='utf-8-sig') as f:
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
                
                # 按考试分组题目
                exams = {}
                current_exam_key = None
                row_num = 1
                
                for row in reader:
                    row_num += 1
                    
                    # 解析考试信息（每行可能包含或更新考试信息）
                    exam_info = self.parse_exam_info(row, headers)
                    if exam_info:
                        current_exam_key = exam_info['title']
                        if current_exam_key not in exams:
                            exams[current_exam_key] = {
                                **exam_info,
                                "questions": []
                            }
                    
                    # 解析题目
                    question = self.parse_question(row, headers)
                    if question and current_exam_key:
                        exams[current_exam_key]["questions"].append(question)
                    elif question:
                        print(f"  ⚠ 第 {row_num} 行：题目缺少考试信息，跳过")
                
                if not exams:
                    print("✗ 没有有效的考试数据")
                    return
                
                print(f"\n开始导入 {len(exams)} 个考试...\n")
                success_count = 0
                total_questions = 0
                
                for exam_title, exam_data in exams.items():
                    question_count = len(exam_data.get("questions", []))
                    total_questions += question_count
                    
                    if question_count == 0:
                        print(f"[跳过] {exam_title} (没有题目)")
                        continue
                    
                    print(f"[{success_count + 1}/{len(exams)}] {exam_title} ({question_count} 题)")
                    if self.create_exam(exam_data):
                        success_count += 1
                
                print(f"\n导入完成: 成功 {success_count}/{len(exams)}，共 {total_questions} 道题目")
                
        except FileNotFoundError:
            print(f"✗ 文件不存在: {csv_path}")
        except Exception as e:
            print(f"✗ 读取 CSV 文件失败: {e}")


def main():
    parser = argparse.ArgumentParser(description='批量导入考试题（CSV 格式）')
    parser.add_argument('--host', default='http://localhost:8080', help='API 服务器地址')
    parser.add_argument('--username', default='admin', help='管理员用户名（工号）')
    parser.add_argument('--password', default='admin123456', help='管理员密码')
    parser.add_argument('--data', required=True, help='CSV 数据文件路径')

    args = parser.parse_args()

    importer = ExamImporter(args.host, args.username, args.password)

    if not importer.login():
        sys.exit(1)

    importer.import_from_csv(args.data)


if __name__ == '__main__':
    main()

