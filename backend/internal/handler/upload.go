package handler

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"marketplace/backend/pkg/response"
)

const (
	maxUploadSize  = 5 << 20 // 5MB
	uploadDir      = "uploads"
	magicBytesSize = 12      // enough to detect all supported image types
)

var allowedImageTypes = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".gif":  true,
	".webp": true,
	".bmp":  true,
}

// imageMagicSignatures maps extensions to their file header magic bytes
var imageMagicSignatures = map[string][][]byte{
	".jpg":  {{0xFF, 0xD8, 0xFF}},
	".jpeg": {{0xFF, 0xD8, 0xFF}},
	".png":  {{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}},
	".gif":  {{0x47, 0x49, 0x46, 0x38, 0x37, 0x61}, {0x47, 0x49, 0x46, 0x38, 0x39, 0x61}}, // GIF87a, GIF89a
	".webp": {{0x52, 0x49, 0x46, 0x46}}, // RIFF...WEBP
	".bmp":  {{0x42, 0x4D}},              // BM
}

func validateImageMagicBytes(data []byte, ext string) bool {
	signatures, ok := imageMagicSignatures[ext]
	if !ok {
		return false
	}
	for _, sig := range signatures {
		if bytes.HasPrefix(data, sig) {
			// webp: RIFF....WEBP — also check WEBP marker at offset 8
			if ext == ".webp" && len(data) >= 12 {
				return string(data[8:12]) == "WEBP"
			}
			return true
		}
	}
	return false
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

	// Validate file magic bytes
	buf := make([]byte, magicBytesSize)
	n, err := file.Read(buf)
	if err != nil || n < 4 {
		response.Error(c, response.CodeParamsError, "文件内容无效")
		return
	}
	if !validateImageMagicBytes(buf[:n], ext) {
		response.Error(c, response.CodeParamsError, "文件内容与扩展名不匹配")
		return
	}
	if _, err := file.Seek(0, 0); err != nil {
		response.Error(c, response.CodeInternalError, "上传失败")
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
