package upload

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/wxsimon8888/simonStu/internal/logger"
	"github.com/wxsimon8888/simonStu/internal/response"
)

const uploadDir = "storage/uploads"

func init() {
	os.MkdirAll(uploadDir, 0755)
}

// UploadResult 文件上传返回结果。
type UploadResult struct {
	FileName    string `json:"file_name"`    // 原始文件名
	SavedName   string `json:"saved_name"`   // 服务端保存的文件名（防冲突）
	Size        int64  `json:"size"`         // 文件大小（字节）
	Path        string `json:"path"`         // 存储路径
	URL         string `json:"url"`          // 可通过 HTTP 访问的 URL
	ContentType string `json:"content_type"` // 文件类型
}

// Upload 单文件上传。以 multipart/form-data 格式接收文件。
//
// 请求方式：POST，Content-Type: multipart/form-data
// 字段名：file（文件），支持常见的图片/文档/压缩包格式。
//
// curl 示例：
//
//	curl -X POST http://localhost:8080/commonApi/upload \
//	  -F "file=@/path/to/photo.jpg"
//
// 多文件上传可用多个 -F "file=@..."，接口会分别返回每个文件的信息。
func Upload(c *gin.Context) {
	// 获取上传的文件
	file, err := c.FormFile("file")
	if err != nil {
		logger.Errorf(c, "Upload 获取文件失败: %v", err)
		response.Error(c, http.StatusBadRequest, "请选择要上传的文件")
		return
	}

	// 校验文件大小（限制 10MB）
	const maxSize int64 = 10 << 20 // 10MB
	if file.Size > maxSize {
		logger.Errorf(c, "Upload 文件过大 size=%d", file.Size)
		response.Error(c, http.StatusBadRequest, fmt.Sprintf("文件大小不能超过 %dMB", maxSize>>20))
		return
	}

	// 校验文件扩展名（白名单）
	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowedExts := map[string]string{
		".jpg":  "image/jpeg",
		".jpeg": "image/jpeg",
		".png":  "image/png",
		".gif":  "image/gif",
		".webp": "image/webp",
		".svg":  "image/svg+xml",
		".pdf":  "application/pdf",
		".doc":  "application/msword",
		".docx": "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
		".xls":  "application/vnd.ms-excel",
		".xlsx": "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
		".zip":  "application/zip",
		".txt":  "text/plain",
		".csv":  "text/csv",
		".json": "application/json",
		".md":   "text/markdown",
	}
	contentType, allowed := allowedExts[ext]
	if !allowed {
		logger.Errorf(c, "Upload 不支持的文件类型: %s", ext)
		response.Error(c, http.StatusBadRequest, fmt.Sprintf("不支持的文件类型: %s", ext))
		return
	}

	// 生成唯一文件名（时间戳 + 随机后缀，防止文件名冲突）
	timestamp := time.Now().UnixNano()
	savedName := fmt.Sprintf("%d_%s", timestamp, file.Filename)

	// 确保目录存在
	os.MkdirAll(uploadDir, 0755)

	// 保存文件到磁盘
	savePath := filepath.Join(uploadDir, savedName)
	if err := c.SaveUploadedFile(file, savePath); err != nil {
		logger.Errorf(c, "Upload 保存文件失败 path=%s err=%v", savePath, err)
		response.Error(c, http.StatusInternalServerError, "文件保存失败")
		return
	}

	logger.Infof(c, "Upload 成功 file=%s size=%d path=%s", file.Filename, file.Size, savePath)

	response.Success(c, UploadResult{
		FileName:    file.Filename,
		SavedName:   savedName,
		Size:        file.Size,
		Path:        savePath,
		URL:         fmt.Sprintf("/uploads/%s", savedName),
		ContentType: contentType,
	})
}

// UploadMultiple 多文件上传。一次请求上传多个文件。
//
// curl 示例：
//
//	curl -X POST http://localhost:8080/commonApi/upload/multiple \
//	  -F "files=@/path/to/a.jpg" \
//	  -F "files=@/path/to/b.png"
func UploadMultiple(c *gin.Context) {
	// 使用同一个字段名获取多文件
	form, err := c.MultipartForm()
	if err != nil {
		logger.Errorf(c, "UploadMultiple 解析 multipart 失败: %v", err)
		response.Error(c, http.StatusBadRequest, "解析上传数据失败")
		return
	}

	files := form.File["files"]
	if len(files) == 0 {
		logger.Errorf(c, "UploadMultiple 未上传文件")
		response.Error(c, http.StatusBadRequest, "请选择要上传的文件")
		return
	}

	const maxSize int64 = 10 << 20
	var results []UploadResult
	for _, file := range files {
		if file.Size > maxSize {
			continue // 超过大小限制的跳过
		}

		timestamp := time.Now().UnixNano()
		savedName := fmt.Sprintf("%d_%s", timestamp, file.Filename)
		savePath := filepath.Join(uploadDir, savedName)

		if err := c.SaveUploadedFile(file, savePath); err != nil {
			continue
		}

		results = append(results, UploadResult{
			FileName:  file.Filename,
			SavedName: savedName,
			Size:      file.Size,
			Path:      savePath,
			URL:       fmt.Sprintf("/uploads/%s", savedName),
		})
	}

	logger.Infof(c, "UploadMultiple 成功 count=%d", len(results))

	response.Success(c, gin.H{
		"total":    len(files),
		"uploaded": len(results),
		"files":    results,
	})
}
