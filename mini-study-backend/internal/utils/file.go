package utils

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"
)

// SaveUploadedFile stores the provided file into the target directory.
func SaveUploadedFile(file *multipart.FileHeader, targetDir string) (string, error) {
	if err := os.MkdirAll(targetDir, 0o755); err != nil {
		return "", fmt.Errorf("make upload dir: %w", err)
	}

	filename := fmt.Sprintf("%d-%s", time.Now().UnixNano(), filepath.Base(file.Filename))
	dst := filepath.Join(targetDir, filename)

	srcFile, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("open upload: %w", err)
	}
	defer srcFile.Close()

	out, err := os.Create(dst)
	if err != nil {
		return "", fmt.Errorf("create upload: %w", err)
	}
	defer out.Close()

	if _, err := io.Copy(out, srcFile); err != nil {
		return "", fmt.Errorf("copy upload: %w", err)
	}

	return fmt.Sprintf("/uploads/%s", filename), nil
}
