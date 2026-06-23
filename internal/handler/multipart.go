package handler

import (
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/hhung06/digimap-backend/internal/domain"
	"github.com/hhung06/digimap-backend/internal/service"
)

const maxMultipartRequestBytes int64 = 32 << 20

func bindMultipartData(c *gin.Context, dst any, maxMemory int64) error {
	if dst == nil {
		return domain.NewValidation(map[string]string{"data": "destination is required"})
	}
	if maxMemory <= 0 {
		return domain.NewValidation(map[string]string{"data": "max memory must be positive"})
	}
	if err := parseMultipart(c); err != nil {
		return err
	}

	values := c.Request.MultipartForm.Value["data"]
	if len(values) == 0 {
		return domain.NewValidation(map[string]string{"data": "required"})
	}
	if int64(len(values[0])) > maxMemory {
		return domain.NewValidation(map[string]string{"data": "too large"})
	}
	if err := json.Unmarshal([]byte(values[0]), dst); err != nil {
		return domain.NewValidation(map[string]string{"data": "invalid JSON"})
	}
	return nil
}

func multipartFile(c *gin.Context, field string) (*multipart.FileHeader, bool, error) {
	if err := parseMultipart(c); err != nil {
		return nil, false, err
	}

	files := c.Request.MultipartForm.File[field]
	if len(files) == 0 {
		return nil, false, nil
	}
	if len(files) > 1 {
		return nil, false, domain.NewValidation(map[string]string{field: "multiple files are not allowed"})
	}
	return files[0], true, nil
}

func multipartFiles(c *gin.Context, field string) ([]*multipart.FileHeader, bool, error) {
	if err := parseMultipart(c); err != nil {
		return nil, false, err
	}

	files := c.Request.MultipartForm.File[field]
	if len(files) == 0 {
		return nil, false, nil
	}
	return files, true, nil
}

func mediaUpload(file *multipart.FileHeader) (service.MediaUpload, io.Closer, error) {
	if file == nil {
		return service.MediaUpload{}, nil, domain.NewValidation(map[string]string{"file": "required"})
	}
	opened, err := file.Open()
	if err != nil {
		return service.MediaUpload{}, nil, err
	}
	return service.MediaUpload{
		Filename:    file.Filename,
		ContentType: file.Header.Get("Content-Type"),
		Size:        file.Size,
		Reader:      opened,
	}, opened, nil
}

func parseMultipart(c *gin.Context) error {
	if c.Request.MultipartForm != nil {
		return nil
	}
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxMultipartRequestBytes)
	if err := c.Request.ParseMultipartForm(maxMultipartRequestBytes); err != nil {
		return domain.NewValidation(map[string]string{"multipart": "invalid or too large"})
	}
	return nil
}
