package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// ContentData 学习内容数据结构
type ContentData struct {
	Title           string `json:"title"`
	Type            string `json:"type"` // doc 或 video
	CategoryID      uint   `json:"category_id"`
	FilePath        string `json:"file_path"`
	CoverURL        string `json:"cover_url,omitempty"`
	Summary         string `json:"summary,omitempty"`
	VisibleRoles    string `json:"visible_roles,omitempty"` // employee/manager/both
	Status          string `json:"status,omitempty"`        // draft/published
	DurationSeconds int64  `json:"duration_seconds"`
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

type ContentImporter struct {
	baseURL  string
	username string
	password string
	token    string
	client   *http.Client
}

func NewContentImporter(baseURL, username, password string) *ContentImporter {
	return &ContentImporter{
		baseURL:  strings.TrimSuffix(baseURL, "/"),
		username: username,
		password: password,
		client:   &http.Client{},
	}
}

func (c *ContentImporter) Login() error {
	url := fmt.Sprintf("%s/api/v1/users/login", c.baseURL)
	payload := map[string]string{
		"work_no":  c.username,
		"password": c.password,
	}
	jsonData, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("创建请求失败: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
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

	c.token = loginResp.Data.Token
	fmt.Printf("✓ 登录成功: %s\n", c.username)
	return nil
}

func (c *ContentImporter) ListCategories() ([]map[string]interface{}, error) {
	url := fmt.Sprintf("%s/api/v1/contents/categories", c.baseURL)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var apiResp APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, err
	}

	if apiResp.Code != 200 {
		return nil, fmt.Errorf("获取分类失败: %s", apiResp.Msg)
	}

	categories, ok := apiResp.Data.([]interface{})
	if !ok {
		return nil, fmt.Errorf("分类数据格式错误")
	}

	result := make([]map[string]interface{}, 0, len(categories))
	for _, cat := range categories {
		if catMap, ok := cat.(map[string]interface{}); ok {
			result = append(result, catMap)
		}
	}

	return result, nil
}

func (c *ContentImporter) UploadFile(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("打开文件失败: %v", err)
	}
	defer file.Close()

	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	part, err := writer.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
		return "", fmt.Errorf("创建表单字段失败: %v", err)
	}

	_, err = file.Seek(0, 0)
	if err != nil {
		return "", err
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return "", fmt.Errorf("复制文件内容失败: %v", err)
	}

	writer.Close()

	url := fmt.Sprintf("%s/api/v1/files/upload", c.baseURL)
	req, err := http.NewRequest("POST", url, &requestBody)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := c.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var apiResp APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return "", err
	}

	if apiResp.Code != 200 {
		return "", fmt.Errorf("%s", apiResp.Msg)
	}

	dataMap, ok := apiResp.Data.(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("响应数据格式错误")
	}

	path, ok := dataMap["path"].(string)
	if !ok {
		return "", fmt.Errorf("路径字段缺失")
	}

	fmt.Printf("  ✓ 上传成功: %s -> %s\n", filepath.Base(filePath), path)
	return path, nil
}

func (c *ContentImporter) CreateContent(content ContentData) error {
	// 处理封面图上传
	if content.CoverURL != "" {
		// 如果 cover_url 是本地文件路径（不是 http/https 或 /uploads/ 开头），则先上传
		if !strings.HasPrefix(content.CoverURL, "http://") &&
			!strings.HasPrefix(content.CoverURL, "https://") &&
			!strings.HasPrefix(content.CoverURL, "/uploads/") {
			uploadPath, err := c.UploadFile(content.CoverURL)
			if err != nil {
				fmt.Printf("  ⚠ 封面图上传失败: %v，将跳过封面图\n", err)
				content.CoverURL = ""
			} else {
				content.CoverURL = uploadPath
			}
		}
	}

	url := fmt.Sprintf("%s/api/v1/admin/contents", c.baseURL)
	jsonData, _ := json.Marshal(content)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
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
		if id, exists := dataMap["id"]; exists {
			fmt.Printf("  ✓ 创建成功: %s (ID: %.0f)\n", content.Title, id)
		} else {
			fmt.Printf("  ✓ 创建成功: %s\n", content.Title)
		}
	} else {
		fmt.Printf("  ✓ 创建成功: %s\n", content.Title)
	}

	return nil
}

func (c *ContentImporter) ImportFromFile(filePath string, showCategories bool) error {
	if c.token == "" {
		return fmt.Errorf("请先登录")
	}

	// 显示分类列表（可选）
	if showCategories {
		fmt.Println("\n可用分类:")
		categories, err := c.ListCategories()
		if err != nil {
			return err
		}
		for _, cat := range categories {
			id := cat["id"]
			name := cat["name"]
			roleScope := cat["role_scope"]
			fmt.Printf("  - ID: %v, 名称: %v, 角色范围: %v\n", id, name, roleScope)
		}
	}

	// 读取并导入数据
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("读取文件失败: %v", err)
	}

	var contents []ContentData
	if err := json.Unmarshal(data, &contents); err != nil {
		return fmt.Errorf("解析 JSON 失败: %v", err)
	}

	fmt.Printf("\n开始导入 %d 条学习内容...\n\n", len(contents))
	successCount := 0
	for idx, content := range contents {
		fmt.Printf("[%d/%d] %s\n", idx+1, len(contents), content.Title)
		if err := c.CreateContent(content); err != nil {
			fmt.Printf("  ✗ 创建失败: %v\n", err)
		} else {
			successCount++
		}
	}

	fmt.Printf("\n导入完成: 成功 %d/%d\n", successCount, len(contents))
	return nil
}

func main() {
	var (
		host     = flag.String("host", "http://localhost:8080", "API 服务器地址")
		username = flag.String("username", "admin", "管理员用户名（工号）")
		password = flag.String("password", "admin123456", "管理员密码")
		dataFile = flag.String("data", "", "JSON 数据文件路径（必填）")
		showCats = flag.Bool("categories", false, "显示可用分类列表")
	)
	flag.Parse()

	if *dataFile == "" {
		fmt.Fprintf(os.Stderr, "错误: 必须指定数据文件路径 (--data)\n")
		flag.Usage()
		os.Exit(1)
	}

	importer := NewContentImporter(*host, *username, *password)

	if err := importer.Login(); err != nil {
		fmt.Fprintf(os.Stderr, "✗ %v\n", err)
		os.Exit(1)
	}

	if err := importer.ImportFromFile(*dataFile, *showCats); err != nil {
		fmt.Fprintf(os.Stderr, "✗ %v\n", err)
		os.Exit(1)
	}
}
