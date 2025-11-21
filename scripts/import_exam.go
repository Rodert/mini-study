package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

// ExamOption 考试选项
type ExamOption struct {
	Label     string `json:"label"`
	Content   string `json:"content"`
	IsCorrect bool   `json:"is_correct"`
	SortOrder int    `json:"sort_order,omitempty"`
}

// ExamQuestion 考试题目
type ExamQuestion struct {
	Type     string       `json:"type"`     // single 或 multiple
	Stem     string       `json:"stem"`     // 题干
	Score    int          `json:"score"`    // 分值
	Analysis string       `json:"analysis,omitempty"`
	Options  []ExamOption `json:"options"`
}

// ExamData 考试数据结构
type ExamData struct {
	Title            string         `json:"title"`
	Description      string         `json:"description,omitempty"`
	Status           string         `json:"status,omitempty"`           // draft/published/archived
	TargetRole       string         `json:"target_role,omitempty"`      // employee/manager/all
	TimeLimitMinutes int            `json:"time_limit_minutes,omitempty"`
	PassScore        int            `json:"pass_score"`
	Questions        []ExamQuestion `json:"questions"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"message"`
	Data struct {
		Token string `json:"token"`
	} `json:"data"`
}

// APIResponse 通用 API 响应
type APIResponse struct {
	Code int         `json:"code"`
	Msg  string      `json:"message"`
	Data interface{} `json:"data"`
}

type ExamImporter struct {
	baseURL  string
	username string
	password string
	token    string
	client   *http.Client
}

func NewExamImporter(baseURL, username, password string) *ExamImporter {
	return &ExamImporter{
		baseURL:  strings.TrimSuffix(baseURL, "/"),
		username: username,
		password: password,
		client:   &http.Client{},
	}
}

func (e *ExamImporter) Login() error {
	url := fmt.Sprintf("%s/api/v1/users/login", e.baseURL)
	payload := map[string]string{
		"work_no":  e.username,
		"password": e.password,
	}
	jsonData, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("创建请求失败: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := e.client.Do(req)
	if err != nil {
		return fmt.Errorf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	var loginResp LoginResponse
	if err := json.NewDecoder(resp.Body).Decode(&loginResp); err != nil {
		return fmt.Errorf("解析响应失败: %v", err)
	}

	if loginResp.Code != 200 || loginResp.Data.Token == "" {
		return fmt.Errorf("登录失败: %s", loginResp.Msg)
	}

	e.token = loginResp.Data.Token
	fmt.Printf("✓ 登录成功: %s\n", e.username)
	return nil
}

func (e *ExamImporter) CreateExam(exam ExamData) error {
	url := fmt.Sprintf("%s/api/v1/admin/exams", e.baseURL)
	jsonData, _ := json.Marshal(exam)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", e.token))
	req.Header.Set("Content-Type", "application/json")

	resp, err := e.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var apiResp APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return err
	}

	if apiResp.Code != 200 {
		return fmt.Errorf("%s", apiResp.Msg)
	}

	dataMap, ok := apiResp.Data.(map[string]interface{})
	if ok {
		id := "?"
		questionCount := "?"
		totalScore := "?"
		if val, exists := dataMap["id"]; exists {
			id = fmt.Sprintf("%.0f", val)
		}
		if val, exists := dataMap["question_count"]; exists {
			questionCount = fmt.Sprintf("%.0f", val)
		}
		if val, exists := dataMap["total_score"]; exists {
			totalScore = fmt.Sprintf("%.0f", val)
		}
		fmt.Printf("  ✓ 创建成功: %s (ID: %s, 题目: %s, 总分: %s)\n",
			exam.Title, id, questionCount, totalScore)
	} else {
		fmt.Printf("  ✓ 创建成功: %s\n", exam.Title)
	}

	return nil
}

func (e *ExamImporter) ImportFromFile(filePath string) error {
	if e.token == "" {
		return fmt.Errorf("请先登录")
	}

	// 读取并导入数据
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("读取文件失败: %v", err)
	}

	var exams []ExamData
	if err := json.Unmarshal(data, &exams); err != nil {
		return fmt.Errorf("解析 JSON 失败: %v", err)
	}

	fmt.Printf("\n开始导入 %d 个考试...\n\n", len(exams))
	successCount := 0
	for idx, exam := range exams {
		questionCount := len(exam.Questions)
		fmt.Printf("[%d/%d] %s (%d 题)\n", idx+1, len(exams), exam.Title, questionCount)
		if err := e.CreateExam(exam); err != nil {
			fmt.Printf("  ✗ 创建失败: %v\n", err)
		} else {
			successCount++
		}
	}

	fmt.Printf("\n导入完成: 成功 %d/%d\n", successCount, len(exams))
	return nil
}

func main() {
	var (
		host     = flag.String("host", "http://localhost:8080", "API 服务器地址")
		username = flag.String("username", "admin", "管理员用户名（工号）")
		password = flag.String("password", "admin123456", "管理员密码")
		dataFile = flag.String("data", "", "JSON 数据文件路径（必填）")
	)
	flag.Parse()

	if *dataFile == "" {
		fmt.Fprintf(os.Stderr, "错误: 必须指定数据文件路径 (--data)\n")
		flag.Usage()
		os.Exit(1)
	}

	importer := NewExamImporter(*host, *username, *password)

	if err := importer.Login(); err != nil {
		fmt.Fprintf(os.Stderr, "✗ %v\n", err)
		os.Exit(1)
	}

	if err := importer.ImportFromFile(*dataFile); err != nil {
		fmt.Fprintf(os.Stderr, "✗ %v\n", err)
		os.Exit(1)
	}
}

