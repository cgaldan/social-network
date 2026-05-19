package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"social-network/internal/config"
	"social-network/internal/domain"
	"social-network/internal/service"
	"social-network/packages/logger"
	"strings"
	"time"
)

const UploadURLPrefix = "/uploads/"

var allowedContentTypes = map[string]string{
	"image/jpeg": ".jpg",
	"image/png":  ".png",
	"image/gif":  ".gif",
	"image/webp": ".webp",
}

type UploadHandler struct {
	authService service.AuthServiceInterface
	logger      *logger.Logger
	config      *config.Config
}

func NewUploadHandler(authService service.AuthServiceInterface, logger *logger.Logger, config *config.Config) *UploadHandler {
	return &UploadHandler{
		authService: authService,
		logger:      logger,
		config:      config,
	}
}

func (h *UploadHandler) UploadImage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	token := r.Header.Get("Authorization")
	if _, err := h.authService.ValidateSession(token); err != nil {
		json.NewEncoder(w).Encode(domain.UploadResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, h.config.Upload.MaxFileSize)

	if err := r.ParseMultipartForm(h.config.Upload.MaxFileSize); err != nil {
		json.NewEncoder(w).Encode(domain.UploadResponse{
			Success: false,
			Message: "File too large or invalid form (max 5MB)",
		})
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		json.NewEncoder(w).Encode(domain.UploadResponse{
			Success: false,
			Message: "Missing file field",
		})
		return
	}
	defer file.Close()

	sniff := make([]byte, 512)
	n, err := file.Read(sniff)
	if err != nil && err != io.EOF {
		json.NewEncoder(w).Encode(domain.UploadResponse{
			Success: false,
			Message: "Failed to read file",
		})
		return
	}
	contentType := http.DetectContentType(sniff[:n])
	ext, ok := allowedContentTypes[contentType]
	if !ok {
		json.NewEncoder(w).Encode(domain.UploadResponse{
			Success: false,
			Message: "Only JPEG, PNG, GIF, or WEBP images are allowed",
		})
		return
	}

	if _, err := file.Seek(0, io.SeekStart); err != nil {
		json.NewEncoder(w).Encode(domain.UploadResponse{
			Success: false,
			Message: "Failed to rewind file",
		})
		return
	}

	if err := os.MkdirAll(h.config.Upload.UploadPath, 0o755); err != nil {
		h.logger.Error("Failed to create upload dir", "error", err)
		json.NewEncoder(w).Encode(domain.UploadResponse{
			Success: false,
			Message: "Server upload error",
		})
		return
	}

	name, err := uniqueFilename(ext)
	if err != nil {
		h.logger.Error("Failed to generate filename", "error", err)
		json.NewEncoder(w).Encode(domain.UploadResponse{
			Success: false,
			Message: "Server upload error",
		})
		return
	}

	dstPath := filepath.Join(h.config.Upload.UploadPath, name)
	dst, err := os.Create(dstPath)
	if err != nil {
		h.logger.Error("Failed to create upload file", "error", err, "path", dstPath)
		json.NewEncoder(w).Encode(domain.UploadResponse{
			Success: false,
			Message: "Server upload error",
		})
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		h.logger.Error("Failed to write upload", "error", err, "path", dstPath)
		os.Remove(dstPath)
		json.NewEncoder(w).Encode(domain.UploadResponse{
			Success: false,
			Message: "Failed to save file",
		})
		return
	}

	h.logger.Info("File uploaded", "name", name, "size", header.Size, "contentType", contentType)

	json.NewEncoder(w).Encode(domain.UploadResponse{
		Success: true,
		Message: "Upload successful",
		URL:     UploadURLPrefix + name,
	})
}

func uniqueFilename(ext string) (string, error) {
	buf := make([]byte, 8)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	if !strings.HasPrefix(ext, ".") {
		ext = "." + ext
	}
	return fmt.Sprintf("%d-%s%s", time.Now().UnixNano(), hex.EncodeToString(buf), ext), nil
}
