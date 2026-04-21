package handler

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"marketplace/backend/pkg/response"
)

const (
	maxUploadSize = 5 << 20 // 5MB
	uploadDir     = "uploads"
)

var allowedImageTypes = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".gif":  true,
	".webp": true,
	".bmp":  true,
}

// UploadFile handles image file upload
func (h *RechargeHandler) UploadFile(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		response.ParamsError(c, "请选择要上传的文件")
		return
	}
	defer file.Close()

	// Check file size
	if header.Size > maxUploadSize {
		response.Error(c, response.CodeParamsError, "文件大小不能超过5MB")
		return
	}

	// Check file extension
	ext := strings.ToLower(filepath.Ext(header.Filename))
	if !allowedImageTypes[ext] {
		response.Error(c, response.CodeParamsError, "只支持图片文件（jpg/png/gif/webp）")
		return
	}

	// Generate unique filename
	filename := uuid.New().String() + ext
	savePath := filepath.Join(uploadDir, filename)

	// Ensure upload directory exists
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		response.Error(c, response.CodeInternalError, "上传失败")
		return
	}

	// Save file
	if err := c.SaveUploadedFile(header, savePath); err != nil {
		response.Error(c, response.CodeInternalError, "保存文件失败")
		return
	}

	// Return URL (relative to server root)
	url := "/uploads/" + filename
	response.Success(c, gin.H{"url": url})
}
