package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/hhung06/digimap-backend/internal/domain"
	"github.com/hhung06/digimap-backend/internal/dto"
	"github.com/hhung06/digimap-backend/internal/enricher"
	"github.com/hhung06/digimap-backend/internal/handler/middleware"
	"github.com/hhung06/digimap-backend/internal/service"
)

type surveyHandler struct {
	svc      service.SurveyService
	enrichers *enricher.Registry
}

func newSurveyHandler(svc service.SurveyService, enrichers *enricher.Registry) *surveyHandler {
	return &surveyHandler{svc: svc, enrichers: enrichers}
}

func (h *surveyHandler) List(c *gin.Context) {
	venueID, err := parseVenueID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid venue id"))
		return
	}
	p := paginationFromQuery(c)
	surveys, total, err := h.svc.List(c.Request.Context(), venueID, p)
	if err != nil {
		respondError(c, err)
		return
	}
	extras, _ := h.enrichers.EnrichForVenue(c.Request.Context(), venueID, enricher.ResourceSurvey)
	items := make([]any, len(surveys))
	for i, s := range surveys {
		items[i] = enricher.MergeInto(dto.SurveyToResponse(s), extras)
	}
	c.JSON(http.StatusOK, dto.Paginated(items, total, p.Page, p.PageSize))
}

func (h *surveyHandler) Get(c *gin.Context) {
	id, err := uuid.Parse(c.Param("surveyID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid survey id"))
		return
	}
	venueID, err := parseVenueID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid venue id"))
		return
	}
	s, err := h.svc.Get(c.Request.Context(), id)
	if err != nil {
		respondError(c, err)
		return
	}
	extras, _ := h.enrichers.EnrichForVenue(c.Request.Context(), venueID, enricher.ResourceSurvey)
	c.JSON(http.StatusOK, dto.OK(enricher.MergeInto(dto.SurveyToResponse(s), extras)))
}

func (h *surveyHandler) Create(c *gin.Context) {
	venueID, err := parseVenueID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid venue id"))
		return
	}
	var req dto.CreateSurveyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.FailMessages(dto.CodeValidationError, bindingErrors(err)))
		return
	}
	userID := middleware.GetUserID(c)
	s := &domain.Survey{
		VenueID: &venueID, ExternalID: req.ExternalID,
		Title: req.Title, Content: req.Content,
		StartDate: req.StartDate, EndDate: req.EndDate,
		Status: req.Status, PublishType: req.PublishType,
		IsForced: req.IsForced, App: req.App,
		SegmentFilters: req.SegmentFilters, CreatedBy: &userID,
	}
	if err := h.svc.Create(c.Request.Context(), s); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusCreated, dto.OK(dto.SurveyToResponse(s)))
}

func (h *surveyHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("surveyID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid survey id"))
		return
	}
	var req dto.UpdateSurveyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.FailMessages(dto.CodeValidationError, bindingErrors(err)))
		return
	}
	s := &domain.Survey{
		ID: id, ExternalID: req.ExternalID,
		Title: req.Title, Content: req.Content,
		StartDate: req.StartDate, EndDate: req.EndDate,
		Status: req.Status, PublishType: req.PublishType,
		IsForced: req.IsForced, App: req.App,
		SegmentFilters: req.SegmentFilters,
	}
	if err := h.svc.Update(c.Request.Context(), s); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(dto.SurveyToResponse(s)))
}

func (h *surveyHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("surveyID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid survey id"))
		return
	}
	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

// ── Questions ─────────────────────────────────────────────────────────────────

func (h *surveyHandler) CreateQuestion(c *gin.Context) {
	surveyID, err := uuid.Parse(c.Param("surveyID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid survey id"))
		return
	}
	var req dto.QuestionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.FailMessages(dto.CodeValidationError, bindingErrors(err)))
		return
	}
	q := &domain.Question{
		SurveyID: surveyID, QuestionNumber: req.QuestionNumber,
		QuestionType: req.QuestionType, QuestionText: req.QuestionText,
		IsRequired: req.IsRequired, IsOther: req.IsOther,
	}
	if err := h.svc.CreateQuestion(c.Request.Context(), q); err != nil {
		respondError(c, err)
		return
	}
	// Create options if provided
	for _, oReq := range req.Options {
		opt := &domain.Option{
			QuestionID: q.ID, OptionNumber: oReq.OptionNumber, OptionText: oReq.OptionText,
		}
		if err := h.svc.CreateOption(c.Request.Context(), opt); err != nil {
			respondError(c, err)
			return
		}
		q.Options = append(q.Options, opt)
	}
	qr := dto.QuestionResponse{
		ID: q.ID, SurveyID: q.SurveyID,
		QuestionNumber: q.QuestionNumber, QuestionType: q.QuestionType,
		QuestionText: q.QuestionText, IsRequired: q.IsRequired, IsOther: q.IsOther,
		CreatedAt: q.CreatedAt, UpdatedAt: q.UpdatedAt,
	}
	for _, o := range q.Options {
		qr.Options = append(qr.Options, dto.OptionResponse{
			ID: o.ID, QuestionID: o.QuestionID,
			OptionNumber: o.OptionNumber, OptionText: o.OptionText, CreatedAt: o.CreatedAt,
		})
	}
	c.JSON(http.StatusCreated, dto.OK(qr))
}

func (h *surveyHandler) UpdateQuestion(c *gin.Context) {
	questionID, err := uuid.Parse(c.Param("questionID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid question id"))
		return
	}
	surveyID, err := uuid.Parse(c.Param("surveyID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid survey id"))
		return
	}
	var req dto.QuestionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.FailMessages(dto.CodeValidationError, bindingErrors(err)))
		return
	}
	q := &domain.Question{
		ID: questionID, SurveyID: surveyID,
		QuestionNumber: req.QuestionNumber, QuestionType: req.QuestionType,
		QuestionText: req.QuestionText, IsRequired: req.IsRequired, IsOther: req.IsOther,
	}
	if err := h.svc.UpdateQuestion(c.Request.Context(), q); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(dto.QuestionResponse{
		ID: q.ID, SurveyID: q.SurveyID,
		QuestionNumber: q.QuestionNumber, QuestionType: q.QuestionType,
		QuestionText: q.QuestionText, IsRequired: q.IsRequired, IsOther: q.IsOther,
		UpdatedAt: q.UpdatedAt,
	}))
}

func (h *surveyHandler) DeleteQuestion(c *gin.Context) {
	id, err := uuid.Parse(c.Param("questionID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid question id"))
		return
	}
	if err := h.svc.DeleteQuestion(c.Request.Context(), id); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

// ── Options ───────────────────────────────────────────────────────────────────

func (h *surveyHandler) CreateOption(c *gin.Context) {
	questionID, err := uuid.Parse(c.Param("questionID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid question id"))
		return
	}
	var req dto.OptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.FailMessages(dto.CodeValidationError, bindingErrors(err)))
		return
	}
	o := &domain.Option{QuestionID: questionID, OptionNumber: req.OptionNumber, OptionText: req.OptionText}
	if err := h.svc.CreateOption(c.Request.Context(), o); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusCreated, dto.OK(dto.OptionResponse{
		ID: o.ID, QuestionID: o.QuestionID, OptionNumber: o.OptionNumber, OptionText: o.OptionText, CreatedAt: o.CreatedAt,
	}))
}

func (h *surveyHandler) UpdateOption(c *gin.Context) {
	id, err := uuid.Parse(c.Param("optionID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid option id"))
		return
	}
	questionID, err := uuid.Parse(c.Param("questionID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid question id"))
		return
	}
	var req dto.OptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.FailMessages(dto.CodeValidationError, bindingErrors(err)))
		return
	}
	o := &domain.Option{ID: id, QuestionID: questionID, OptionNumber: req.OptionNumber, OptionText: req.OptionText}
	if err := h.svc.UpdateOption(c.Request.Context(), o); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(dto.OptionResponse{
		ID: o.ID, QuestionID: o.QuestionID, OptionNumber: o.OptionNumber, OptionText: o.OptionText,
	}))
}

func (h *surveyHandler) DeleteOption(c *gin.Context) {
	id, err := uuid.Parse(c.Param("optionID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid option id"))
		return
	}
	if err := h.svc.DeleteOption(c.Request.Context(), id); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

// ── Responses ─────────────────────────────────────────────────────────────────

func (h *surveyHandler) ListResponses(c *gin.Context) {
	surveyID, err := uuid.Parse(c.Param("surveyID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid survey id"))
		return
	}
	p := paginationFromQuery(c)
	responses, total, err := h.svc.ListResponses(c.Request.Context(), surveyID, p)
	if err != nil {
		respondError(c, err)
		return
	}
	items := make([]dto.SurveyResponseResponse, len(responses))
	for i, r := range responses {
		items[i] = dto.SurveyResponseToResponse(r)
	}
	c.JSON(http.StatusOK, dto.Paginated(items, total, p.Page, p.PageSize))
}

func (h *surveyHandler) SubmitResponse(c *gin.Context) {
	surveyID, err := uuid.Parse(c.Param("surveyID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid survey id"))
		return
	}
	var req dto.SubmitSurveyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.FailMessages(dto.CodeValidationError, bindingErrors(err)))
		return
	}
	resp := &domain.SurveyResponse{SurveyID: surveyID}
	for _, a := range req.Answers {
		resp.Answers = append(resp.Answers, &domain.SurveyAnswer{
			QuestionID: a.QuestionID, OptionID: a.OptionID, AnswerText: a.AnswerText,
		})
	}
	if err := h.svc.SubmitResponse(c.Request.Context(), resp); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusCreated, dto.OK(dto.SurveyResponseToResponse(resp)))
}
