package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/hhung06/digimap-backend/internal/dto"
	"github.com/hhung06/digimap-backend/internal/service"
)

type mediaHandler struct {
	mediaSvc service.MediaService
}

func newMediaHandler(mediaSvc service.MediaService) *mediaHandler {
	return &mediaHandler{mediaSvc: mediaSvc}
}

// Upload handles multipart media uploads for any backend-owned field.
// The multipart request must include:
//   - data: JSON {"entity":"locations","record_id":"<uuid>","field":"common_logo"}
//   - file: the image file
func (h *mediaHandler) Upload(c *gin.Context) {
	var req dto.MediaUploadRequest
	if err := bindMultipartData(c, &req, 1<<20); err != nil {
		respondError(c, err)
		return
	}

	file, found, err := multipartFile(c, "file")
	if err != nil {
		respondError(c, err)
		return
	}
	if !found {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "file is required"))
		return
	}

	media, mediaCloser, err := mediaUpload(file)
	if err != nil {
		respondError(c, err)
		return
	}
	defer func() { _ = mediaCloser.Close() }()

	ctx := c.Request.Context()
	target := service.MediaTarget{Entity: req.Entity, RecordID: req.RecordID, Field: req.Field}

	key, err := h.mediaSvc.Upload(ctx, target, media)
	if err != nil {
		respondError(c, err)
		return
	}

	url := h.mediaSvc.URL(ctx, target, key)
	c.JSON(http.StatusOK, dto.OK(dto.MediaUploadResponse{Key: key, URL: url}))
}
