package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/hhung06/digimap-backend/internal/domain"
	"github.com/hhung06/digimap-backend/internal/dto"
	"github.com/hhung06/digimap-backend/internal/enricher"
	"github.com/hhung06/digimap-backend/internal/handler/middleware"
	"github.com/hhung06/digimap-backend/internal/service"
)

type surveyHandler struct {
	svc       service.SurveyService
	enrichers *enricher.Registry
}

func newSurveyHandler(svc service.SurveyService, enrichers *enricher.Registry) *surveyHandler {
	return &surveyHandler{svc: svc, enrichers: enrichers}
}

// @Summary     List surveys
// @Description List surveys for a venue
// @Tags        surveys
// @Produce     json
// @Security    BearerAuth
// @Param       id        path     string true  "Venue ID"
// @Param       page      query    int    false "Page number"
// @Param       page_size query    int    false "Page size"
// @Param       keyword      query    string false "Search keyword"
// @Param       status       query    int    false "Survey status"
// @Param       publish_type query    int    false "Survey publish type"
// @Success     200       {object} dto.Response{data=[]interface{},metadata=dto.PaginationMeta}
// @Failure     400       {object} dto.Response
// @Failure     401       {object} dto.Response
// @Router      /venues/{id}/surveys [get]
func (h *surveyHandler) List(c *gin.Context) {
	venueID, err := parseVenueID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid venue id"))
		return
	}
	p := paginationFromQuery(c)
	filter := domain.SurveyListFilter{Keyword: strings.TrimSpace(c.Query("keyword"))}
	if raw := strings.TrimSpace(c.Query("status")); raw != "" {
		status, err := strconv.Atoi(raw)
		if err != nil {
			respondError(c, domain.NewValidation(map[string]string{"status": "must be an integer"}))
			return
		}
		filter.Status = &status
	}
	if raw := strings.TrimSpace(c.Query("publish_type")); raw != "" {
		publishType, err := strconv.Atoi(raw)
		if err != nil {
			respondError(c, domain.NewValidation(map[string]string{"publish_type": "must be an integer"}))
			return
		}
		filter.PublishType = &publishType
	}
	surveys, total, err := h.svc.List(c.Request.Context(), venueID, filter, p)
	if err != nil {
		respondError(c, err)
		return
	}
	var extras map[string]any
	if h.enrichers != nil {
		extras, _ = h.enrichers.EnrichForVenue(c.Request.Context(), venueID, enricher.ResourceSurvey)
	}
	items := make([]any, len(surveys))
	for i, s := range surveys {
		items[i] = enricher.MergeInto(dto.SurveyToResponse(s), extras)
	}
	c.JSON(http.StatusOK, dto.Paginated(items, total, p.Page, p.PageSize))
}

// @Summary     Get survey
// @Description Get a survey by ID
// @Tags        surveys
// @Produce     json
// @Security    BearerAuth
// @Param       id       path     string true "Venue ID"
// @Param       surveyID path     string true "Survey ID"
// @Success     200      {object} dto.Response{data=dto.SurveyResponse}
// @Failure     400      {object} dto.Response
// @Failure     401      {object} dto.Response
// @Failure     404      {object} dto.Response
// @Router      /venues/{id}/surveys/{surveyID} [get]
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

// @Summary     Create survey
// @Description Create a new survey (requires editor role)
// @Tags        surveys
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       id   path     string                    true "Venue ID"
// @Param       body body     dto.CreateSurveyRequest   true "Survey details"
// @Success     201  {object} dto.Response{data=dto.SurveyResponse}
// @Failure     400  {object} dto.Response
// @Failure     401  {object} dto.Response
// @Failure     403  {object} dto.Response
// @Router      /venues/{id}/surveys [post]
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

// @Summary     Update survey
// @Description Update a survey (requires editor role)
// @Tags        surveys
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       id       path     string                  true "Venue ID"
// @Param       surveyID path     string                  true "Survey ID"
// @Param       body     body     dto.UpdateSurveyRequest true "Survey details"
// @Success     200      {object} dto.Response{data=dto.SurveyResponse}
// @Failure     400      {object} dto.Response
// @Failure     401      {object} dto.Response
// @Failure     403      {object} dto.Response
// @Router      /venues/{id}/surveys/{surveyID} [put]
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
	s, err := h.svc.Get(c.Request.Context(), id)
	if err != nil {
		respondError(c, err)
		return
	}
	req.ApplyTo(s)
	if err := h.svc.Update(c.Request.Context(), s); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(dto.SurveyToResponse(s)))
}

// @Summary     Delete survey
// @Description Delete a survey (requires editor role)
// @Tags        surveys
// @Produce     json
// @Security    BearerAuth
// @Param       id       path string true "Venue ID"
// @Param       surveyID path string true "Survey ID"
// @Success     204
// @Failure     400 {object} dto.Response
// @Failure     401 {object} dto.Response
// @Failure     403 {object} dto.Response
// @Router      /venues/{id}/surveys/{surveyID} [delete]
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

// @Summary     Create survey question
// @Description Add a question to a survey (requires editor role)
// @Tags        surveys
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       id       path     string               true "Venue ID"
// @Param       surveyID path     string               true "Survey ID"
// @Param       body     body     dto.QuestionRequest  true "Question details"
// @Success     201      {object} dto.Response{data=dto.QuestionResponse}
// @Failure     400      {object} dto.Response
// @Failure     401      {object} dto.Response
// @Failure     403      {object} dto.Response
// @Router      /venues/{id}/surveys/{surveyID}/questions [post]
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

// @Summary     Update survey question
// @Description Update a survey question (requires editor role)
// @Tags        surveys
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       id         path     string               true "Venue ID"
// @Param       surveyID   path     string               true "Survey ID"
// @Param       questionID path     string               true "Question ID"
// @Param       body       body     dto.QuestionRequest  true "Question details"
// @Success     200        {object} dto.Response{data=dto.QuestionResponse}
// @Failure     400        {object} dto.Response
// @Failure     401        {object} dto.Response
// @Failure     403        {object} dto.Response
// @Router      /venues/{id}/surveys/{surveyID}/questions/{questionID} [put]
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

// @Summary     Delete survey question
// @Description Delete a survey question (requires editor role)
// @Tags        surveys
// @Produce     json
// @Security    BearerAuth
// @Param       id         path string true "Venue ID"
// @Param       surveyID   path string true "Survey ID"
// @Param       questionID path string true "Question ID"
// @Success     204
// @Failure     400 {object} dto.Response
// @Failure     401 {object} dto.Response
// @Failure     403 {object} dto.Response
// @Router      /venues/{id}/surveys/{surveyID}/questions/{questionID} [delete]
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

// @Summary     Create question option
// @Description Add an option to a survey question (requires editor role)
// @Tags        surveys
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       id         path     string              true "Venue ID"
// @Param       surveyID   path     string              true "Survey ID"
// @Param       questionID path     string              true "Question ID"
// @Param       body       body     dto.OptionRequest   true "Option details"
// @Success     201        {object} dto.Response{data=dto.OptionResponse}
// @Failure     400        {object} dto.Response
// @Failure     401        {object} dto.Response
// @Failure     403        {object} dto.Response
// @Router      /venues/{id}/surveys/{surveyID}/questions/{questionID}/options [post]
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

// @Summary     Update question option
// @Description Update a survey question option (requires editor role)
// @Tags        surveys
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       id         path     string             true "Venue ID"
// @Param       surveyID   path     string             true "Survey ID"
// @Param       questionID path     string             true "Question ID"
// @Param       optionID   path     string             true "Option ID"
// @Param       body       body     dto.OptionRequest  true "Option details"
// @Success     200        {object} dto.Response{data=dto.OptionResponse}
// @Failure     400        {object} dto.Response
// @Failure     401        {object} dto.Response
// @Failure     403        {object} dto.Response
// @Router      /venues/{id}/surveys/{surveyID}/questions/{questionID}/options/{optionID} [put]
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

// @Summary     Delete question option
// @Description Delete a survey question option (requires editor role)
// @Tags        surveys
// @Produce     json
// @Security    BearerAuth
// @Param       id         path string true "Venue ID"
// @Param       surveyID   path string true "Survey ID"
// @Param       questionID path string true "Question ID"
// @Param       optionID   path string true "Option ID"
// @Success     204
// @Failure     400 {object} dto.Response
// @Failure     401 {object} dto.Response
// @Failure     403 {object} dto.Response
// @Router      /venues/{id}/surveys/{surveyID}/questions/{questionID}/options/{optionID} [delete]
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

// @Summary     List survey responses
// @Description List responses for a survey
// @Tags        surveys
// @Produce     json
// @Security    BearerAuth
// @Param       id        path     string true  "Venue ID"
// @Param       surveyID  path     string true  "Survey ID"
// @Param       page      query    int    false "Page number"
// @Param       page_size query    int    false "Page size"
// @Success     200       {object} dto.Response{data=[]interface{},metadata=dto.PaginationMeta}
// @Failure     400       {object} dto.Response
// @Failure     401       {object} dto.Response
// @Router      /venues/{id}/surveys/{surveyID}/responses [get]
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

// @Summary     Submit survey response
// @Description Submit a response to a survey (public endpoint, no auth required)
// @Tags        surveys
// @Accept      json
// @Produce     json
// @Param       id       path     string                   true "Venue ID"
// @Param       surveyID path     string                   true "Survey ID"
// @Param       body     body     dto.SubmitSurveyRequest  true "Survey response"
// @Success     201      {object} dto.Response{data=dto.SurveyResponseResponse}
// @Failure     400      {object} dto.Response
// @Router      /venues/{id}/surveys/{surveyID}/responses [post]
func (h *surveyHandler) SubmitResponse(c *gin.Context) {
	venueID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid venue id"))
		return
	}
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
		answerText := a.AnswerText
		resp.Answers = append(resp.Answers, &domain.SurveyAnswer{
			QuestionID: a.QuestionID, OptionID: a.OptionID, AnswerText: &answerText,
		})
	}
	if err := h.svc.SubmitVenueResponse(c.Request.Context(), venueID, resp); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusCreated, dto.OK(dto.SurveyResponseToResponse(resp)))
}
