package handler

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupUploadRouter(h *RechargeHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	v1 := r.Group("/api/v1")
	testAuth := func(c *gin.Context) {
		c.Set("user_id", int64(1))
		c.Set("role", "super_admin")
		c.Next()
	}
	authed := v1.Group("")
	authed.Use(testAuth)
	{
		authed.POST("/upload", h.UploadFile)
	}
	return r
}

func makeUploadRequest(router *gin.Engine, filename string, fileContent []byte) *httptest.ResponseRecorder {
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	part, _ := writer.CreateFormFile("file", filename)
	part.Write(fileContent)
	writer.Close()

	req, _ := http.NewRequest("POST", "/api/v1/upload", &buf)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}

// ========== validateImageMagicBytes 单元测试 ==========

func TestValidateImageMagicBytes_JPG(t *testing.T) {
	jpgHeader := []byte{0xFF, 0xD8, 0xFF, 0xE0, 0x00, 0x10}
	assert.True(t, validateImageMagicBytes(jpgHeader, ".jpg"))
	assert.True(t, validateImageMagicBytes(jpgHeader, ".jpeg"))
}

func TestValidateImageMagicBytes_PNG(t *testing.T) {
	pngHeader := []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A, 0x00, 0x00}
	assert.True(t, validateImageMagicBytes(pngHeader, ".png"))
}

func TestValidateImageMagicBytes_GIF87a(t *testing.T) {
	gifHeader := []byte{0x47, 0x49, 0x46, 0x38, 0x37, 0x61, 0x00}
	assert.True(t, validateImageMagicBytes(gifHeader, ".gif"))
}

func TestValidateImageMagicBytes_GIF89a(t *testing.T) {
	gifHeader := []byte{0x47, 0x49, 0x46, 0x38, 0x39, 0x61, 0x00}
	assert.True(t, validateImageMagicBytes(gifHeader, ".gif"))
}

func TestValidateImageMagicBytes_WebP(t *testing.T) {
	// RIFF + file size (4 bytes) + WEBP
	webpHeader := []byte{0x52, 0x49, 0x46, 0x46, 0x00, 0x00, 0x00, 0x00, 0x57, 0x45, 0x42, 0x50}
	assert.True(t, validateImageMagicBytes(webpHeader, ".webp"))
}

func TestValidateImageMagicBytes_WebP_RIFFOnly(t *testing.T) {
	// RIFF header without WEBP marker (should fail)
	riffOnly := []byte{0x52, 0x49, 0x46, 0x46, 0x00, 0x00, 0x00, 0x00, 0x57, 0x41, 0x56, 0x45}
	assert.False(t, validateImageMagicBytes(riffOnly, ".webp"))
}

func TestValidateImageMagicBytes_BMP(t *testing.T) {
	bmpHeader := []byte{0x42, 0x4D, 0x00, 0x00}
	assert.True(t, validateImageMagicBytes(bmpHeader, ".bmp"))
}

func TestValidateImageMagicBytes_FakeJPG(t *testing.T) {
	fakeJPG := []byte{0x3C, 0x3F, 0x70, 0x68, 0x70} // <?php
	assert.False(t, validateImageMagicBytes(fakeJPG, ".jpg"))
}

func TestValidateImageMagicBytes_UnknownExt(t *testing.T) {
	data := []byte{0xFF, 0xD8, 0xFF}
	assert.False(t, validateImageMagicBytes(data, ".txt"))
}

// ========== UploadFile handler 集成测试 ==========

func TestRechargeHandler_UploadFile_ValidJPG(t *testing.T) {
	mockSvc := new(MockRechargeService)
	h := NewRechargeHandler(mockSvc, &MockUserRepo{})
	router := setupUploadRouter(h)

	// Real JPG header
	jpgContent := append([]byte{0xFF, 0xD8, 0xFF, 0xE0}, bytes.Repeat([]byte{0x00}, 100)...)
	w := makeUploadRequest(router, "test.jpg", jpgContent)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, float64(0), resp["code"])

	data := resp["data"].(map[string]interface{})
	url := data["url"].(string)
	assert.Contains(t, url, "/uploads/")
	assert.Contains(t, url, ".jpg")

	// Cleanup
	os.Remove(filepath.Join("uploads", filepath.Base(url)))
}

func TestRechargeHandler_UploadFile_ValidPNG(t *testing.T) {
	mockSvc := new(MockRechargeService)
	h := NewRechargeHandler(mockSvc, &MockUserRepo{})
	router := setupUploadRouter(h)

	pngContent := append([]byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}, bytes.Repeat([]byte{0x00}, 100)...)
	w := makeUploadRequest(router, "test.png", pngContent)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, float64(0), resp["code"])

	data := resp["data"].(map[string]interface{})
	os.Remove(filepath.Join("uploads", filepath.Base(data["url"].(string))))
}

func TestRechargeHandler_UploadFile_FakeJPG(t *testing.T) {
	mockSvc := new(MockRechargeService)
	h := NewRechargeHandler(mockSvc, &MockUserRepo{})
	router := setupUploadRouter(h)

	// PHP webshell disguised as JPG
	webshell := []byte("<?php system($_GET['cmd']); ?>")
	w := makeUploadRequest(router, "shell.jpg", webshell)

	assert.Equal(t, http.StatusOK, w.Code)
	assertResponseCode(t, w.Body.Bytes(), 400)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Contains(t, resp["message"], "不匹配")
}

func TestRechargeHandler_UploadFile_ContentTooShort(t *testing.T) {
	mockSvc := new(MockRechargeService)
	h := NewRechargeHandler(mockSvc, &MockUserRepo{})
	router := setupUploadRouter(h)

	w := makeUploadRequest(router, "tiny.jpg", []byte{0xFF})

	assert.Equal(t, http.StatusOK, w.Code)
	assertResponseCode(t, w.Body.Bytes(), 400)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Contains(t, resp["message"], "无效")
}

func TestRechargeHandler_UploadFile_InvalidExtension(t *testing.T) {
	mockSvc := new(MockRechargeService)
	h := NewRechargeHandler(mockSvc, &MockUserRepo{})
	router := setupUploadRouter(h)

	w := makeUploadRequest(router, "script.php", []byte{0xFF, 0xD8, 0xFF, 0xE0})

	assert.Equal(t, http.StatusOK, w.Code)
	assertResponseCode(t, w.Body.Bytes(), 400)
}

func TestRechargeHandler_UploadFile_NoFile(t *testing.T) {
	mockSvc := new(MockRechargeService)
	h := NewRechargeHandler(mockSvc, &MockUserRepo{})
	router := setupUploadRouter(h)

	req, _ := http.NewRequest("POST", "/api/v1/upload", bytes.NewBufferString(""))
	req.Header.Set("Content-Type", "multipart/form-data")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assertResponseCode(t, w.Body.Bytes(), 400)
}

func TestRechargeHandler_UploadFile_ValidWebP(t *testing.T) {
	mockSvc := new(MockRechargeService)
	h := NewRechargeHandler(mockSvc, &MockUserRepo{})
	router := setupUploadRouter(h)

	// RIFF + size + WEBP
	webpContent := append([]byte{0x52, 0x49, 0x46, 0x46, 0x00, 0x00, 0x00, 0x00, 0x57, 0x45, 0x42, 0x50}, bytes.Repeat([]byte{0x00}, 100)...)
	w := makeUploadRequest(router, "test.webp", webpContent)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, float64(0), resp["code"])

	data := resp["data"].(map[string]interface{})
	os.Remove(filepath.Join("uploads", filepath.Base(data["url"].(string))))
}
