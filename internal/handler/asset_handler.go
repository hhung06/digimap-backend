package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/hhung06/digimap-backend/internal/dto"
	"github.com/hhung06/digimap-backend/internal/handler/middleware"
	"github.com/hhung06/digimap-backend/internal/service"
)

type assetHandler struct{ svc service.AssetService }

func newAssetHandler(svc service.AssetService) *assetHandler {
	return &assetHandler{svc: svc}
}

func (h *assetHandler) List(c *gin.Context) {
	venueID, err := parseVenueID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid venue id"))
		return
	}

	// When the editor requests specific assets by fileId, return only those.
	// Non-UUID identifiers (e.g. SHA-1 hashes) are mapped to the same deterministic
	// UUID v5 used during upload so they resolve to the correct record.
	rawFileIDs := c.QueryArray("file_ids[]")
	if len(rawFileIDs) == 0 {
		rawFileIDs = c.QueryArray("file_ids")
	}
	if len(rawFileIDs) > 0 {
		ids := make([]uuid.UUID, 0, len(rawFileIDs))
		for _, raw := range rawFileIDs {
			id, parseErr := uuid.Parse(raw)
			if parseErr != nil {
				id = uuid.NewSHA1(uuid.NameSpaceOID, []byte(raw))
			}
			ids = append(ids, id)
		}
		assets, err := h.svc.GetByIDs(c.Request.Context(), venueID, ids)
		if err != nil {
			respondError(c, err)
			return
		}
		items := make([]dto.AssetResponse, len(assets))
		for i, a := range assets {
			items[i] = dto.AssetToResponse(a, h.svc.AssetURL(a.Thumbnail), h.svc.AssetURL(a.Material))
		}
		c.JSON(http.StatusOK, dto.OK(items))
		return
	}

	p := paginationFromQuery(c)
	assetType := c.Query("asset_type")
	assets, total, err := h.svc.List(c.Request.Context(), venueID, p, assetType)
	if err != nil {
		respondError(c, err)
		return
	}
	items := make([]dto.AssetResponse, len(assets))
	for i, a := range assets {
		items[i] = dto.AssetToResponse(a, h.svc.AssetURL(a.Thumbnail), h.svc.AssetURL(a.Material))
	}
	c.JSON(http.StatusOK, dto.Paginated(items, total, p.Page, p.PageSize))
}

func (h *assetHandler) Get(c *gin.Context) {
	id, err := uuid.Parse(c.Param("assetID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid asset id"))
		return
	}
	a, err := h.svc.Get(c.Request.Context(), id)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(dto.AssetToResponse(a, h.svc.AssetURL(a.Thumbnail), h.svc.AssetURL(a.Material))))
}

// Create handles POST /venues/:id/assets.
// The editor sends a batch of base64 data-URLs: [{id, file}, ...].
func (h *assetHandler) Create(c *gin.Context) {
	venueID, err := parseVenueID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid venue id"))
		return
	}
	var items []dto.Upload2DAssetItem
	if err := c.ShouldBindJSON(&items); err != nil {
		c.JSON(http.StatusBadRequest, dto.FailMessages(dto.CodeValidationError, bindingErrors(err)))
		return
	}
	results := make([]dto.AssetResponse, 0, len(items))
	for _, item := range items {
		a, err := h.svc.Upload2D(c.Request.Context(), venueID, item.ID, item.File)
		if err != nil {
			respondError(c, err)
			return
		}
		results = append(results, dto.AssetToResponse(a, h.svc.AssetURL(a.Thumbnail), h.svc.AssetURL(a.Material)))
	}
	c.JSON(http.StatusCreated, dto.OK(results))
}

func (h *assetHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("assetID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid asset id"))
		return
	}
	var req dto.UpdateAssetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.FailMessages(dto.CodeValidationError, bindingErrors(err)))
		return
	}
	a, err := h.svc.Update(c.Request.Context(), id, req.Name)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(dto.AssetToResponse(a, h.svc.AssetURL(a.Thumbnail), h.svc.AssetURL(a.Material))))
}

func (h *assetHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("assetID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid asset id"))
		return
	}
	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

// AllAssets handles GET /all-assets/ — global list, no venue scope, optional ?asset_type filter.
// Mirrors Django's AssetAPIView.get() in indoormap-backend.
func (h *assetHandler) AllAssets(c *gin.Context) {
	assetType := c.Query("asset_type")
	assets, err := h.svc.ListAll(c.Request.Context(), assetType)
	if err != nil {
		respondError(c, err)
		return
	}
	items := make([]dto.AssetResponse, len(assets))
	for i, a := range assets {
		items[i] = dto.AssetToResponse(a, h.svc.AssetURL(a.Thumbnail), h.svc.AssetURL(a.Material))
	}
	c.JSON(http.StatusOK, items)
}

// ListLibraryAssets handles GET /library-assets/ — optional ?status and ?type filters.
func (h *assetHandler) ListLibraryAssets(c *gin.Context) {
	status := c.Query("status")
	assetType := c.Query("type")
	assets, err := h.svc.ListLibrary(c.Request.Context(), status, assetType)
	if err != nil {
		respondError(c, err)
		return
	}
	items := make([]dto.AssetResponse, len(assets))
	for i, a := range assets {
		items[i] = dto.AssetToResponse(a, h.svc.AssetURL(a.Thumbnail), h.svc.AssetURL(a.Material))
	}
	c.JSON(http.StatusOK, dto.OK(items))
}

// UploadLibraryAsset handles POST /library-assets — multipart upload to library S3 prefix.
// Global route: no venue scope; venueID is zero UUID.
func (h *assetHandler) UploadLibraryAsset(c *gin.Context) {
	venueID := uuid.UUID{}

	if err := c.Request.ParseMultipartForm(32 << 20); err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid multipart form"))
		return
	}

	var req dto.UploadLibraryAssetRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.FailMessages(dto.CodeValidationError, bindingErrors(err)))
		return
	}

	fh, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "file is required"))
		return
	}
	f, err := fh.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "cannot read file"))
		return
	}
	defer f.Close()

	userID := middleware.GetUserID(c)
	a, err := h.svc.UploadLibrary(c.Request.Context(), service.UploadLibraryAssetInput{
		VenueID:   &venueID,
		Status:    req.Status,
		AssetType: req.AssetType,
		CreatedBy: &userID,
		File: &service.FileUpload{
			Filename:    fh.Filename,
			ContentType: fh.Header.Get("Content-Type"),
			Size:        fh.Size,
			Reader:      f,
		},
	})
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusCreated, dto.OK(dto.AssetToResponse(a, "", "")))
}

// UpdateLibraryAsset handles PUT /library-assets/:assetID — multipart update.
func (h *assetHandler) UpdateLibraryAsset(c *gin.Context) {
	id, err := uuid.Parse(c.Param("assetID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid asset id"))
		return
	}

	if err := c.Request.ParseMultipartForm(32 << 20); err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid multipart form"))
		return
	}

	var req dto.UploadLibraryAssetRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.FailMessages(dto.CodeValidationError, bindingErrors(err)))
		return
	}

	in := service.UploadLibraryAssetInput{
		Status:    req.Status,
		AssetType: req.AssetType,
	}
	if fh, err := c.FormFile("file"); err == nil {
		f, err := fh.Open()
		if err != nil {
			c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "cannot read file"))
			return
		}
		defer f.Close()
		in.File = &service.FileUpload{
			Filename:    fh.Filename,
			ContentType: fh.Header.Get("Content-Type"),
			Size:        fh.Size,
			Reader:      f,
		}
	}

	a, err := h.svc.UpdateLibrary(c.Request.Context(), id, in)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(dto.AssetToResponse(a, "", "")))
}

func (h *assetHandler) Upload3D(c *gin.Context) {
	venueID, err := parseVenueID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid venue id"))
		return
	}

	if err := c.Request.ParseMultipartForm(32 << 20); err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid multipart form"))
		return
	}

	var req dto.Upload3DAssetRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.FailMessages(dto.CodeValidationError, bindingErrors(err)))
		return
	}

	in := service.Upload3DAssetInput{
		VenueID:     &venueID,
		ID:          req.ID,
		Name:        req.Name,
		Description: req.Description,
		FileType:    req.FileType,
		Width:       req.Width,
		Height:      req.Height,
	}
	if uid := middleware.GetUserID(c); uid != [16]byte{} {
		in.CreatedBy = &uid
	}

	if fh, err := c.FormFile("file"); err == nil {
		f, err := fh.Open()
		if err != nil {
			c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "cannot read file"))
			return
		}
		defer f.Close()
		in.File = &service.FileUpload{
			Filename:    fh.Filename,
			ContentType: fh.Header.Get("Content-Type"),
			Size:        fh.Size,
			Reader:      f,
		}
	}

	if fh, err := c.FormFile("thumbnail"); err == nil {
		f, err := fh.Open()
		if err != nil {
			c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "cannot read thumbnail"))
			return
		}
		defer f.Close()
		in.Thumbnail = &service.FileUpload{
			Filename:    fh.Filename,
			ContentType: fh.Header.Get("Content-Type"),
			Size:        fh.Size,
			Reader:      f,
		}
	}

	if fh, err := c.FormFile("material"); err == nil {
		f, err := fh.Open()
		if err != nil {
			c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "cannot read material"))
			return
		}
		defer f.Close()
		in.Material = &service.FileUpload{
			Filename:    fh.Filename,
			ContentType: fh.Header.Get("Content-Type"),
			Size:        fh.Size,
			Reader:      f,
		}
	}

	a, err := h.svc.Upload3D(c.Request.Context(), in)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusCreated, dto.OK(dto.AssetToResponse(a, h.svc.AssetURL(a.Thumbnail), h.svc.AssetURL(a.Material))))
}
