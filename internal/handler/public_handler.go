package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/hhung06/digimap-backend/internal/domain"
	"github.com/hhung06/digimap-backend/internal/dto"
	"github.com/hhung06/digimap-backend/internal/repository"
	"github.com/hhung06/digimap-backend/internal/service"
)

type publicHandler struct {
	venues         service.VenueService
	surveys        service.SurveyService
	appUsers       repository.AppUserRepository
	visitorSurveys service.VisitorSurveySubmissionService
	searchOptions  service.SearchOptionsService
}

func newPublicHandler(venues service.VenueService, surveys service.SurveyService, appUsers repository.AppUserRepository, visitorSurveys service.VisitorSurveySubmissionService, searchOptions ...service.SearchOptionsService) *publicHandler {
	var searchOptionSvc service.SearchOptionsService
	if len(searchOptions) > 0 {
		searchOptionSvc = searchOptions[0]
	}
	return &publicHandler{venues: venues, surveys: surveys, appUsers: appUsers, visitorSurveys: visitorSurveys, searchOptions: searchOptionSvc}
}

func (h *publicHandler) VenueInformation(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(1000, "invalid venue id"))
		return
	}

	key := c.Query("public_key")
	if key == "" {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "public_key is required"))
		return
	}

	v, err := h.venues.Get(c.Request.Context(), id)
	if err != nil {
		respondError(c, err)
		return
	}
	if v.PublicKey != key {
		c.JSON(http.StatusUnauthorized, dto.Fail(dto.CodeAuthRequired, "invalid public_key"))
		return
	}
	c.JSON(http.StatusOK, dto.OK(dto.VenueToResponse(v)))
}

func (h *publicHandler) GetSurvey(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(1000, "invalid survey id"))
		return
	}
	s, err := h.surveys.Get(c.Request.Context(), id)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(dto.SurveyToResponse(s)))
}

func (h *publicHandler) SubmitSurveyResponse(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(1000, "invalid survey id"))
		return
	}
	var req dto.SubmitSurveyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(1000, err.Error()))
		return
	}
	resp := &domain.SurveyResponse{
		SurveyID: id,
		Answers:  make([]*domain.SurveyAnswer, len(req.Answers)),
	}
	for i, a := range req.Answers {
		answerText := a.AnswerText
		resp.Answers[i] = &domain.SurveyAnswer{
			QuestionID: a.QuestionID,
			OptionID:   a.OptionID,
			AnswerText: &answerText,
		}
	}
	if err := h.surveys.SubmitPublicResponse(c.Request.Context(), resp); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusCreated, dto.OK(dto.SurveyResponseToResponse(resp)))
}

func (h *publicHandler) VenueInfo(c *gin.Context) {
	key := c.Query("public_key")
	if key == "" {
		c.JSON(http.StatusBadRequest, dto.Fail(1000, "public_key is required"))
		return
	}
	v, err := h.venues.GetByPublicKey(c.Request.Context(), key)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(dto.VenueToResponse(v)))
}

func (h *publicHandler) SubmitVisitorSurvey(c *gin.Context) {
	var req dto.SubmitVisitorSurveyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.FailMessages(dto.CodeValidationError, bindingErrors(err)))
		return
	}
	visitor, err := h.visitorSurveys.Submit(c.Request.Context(), service.VisitorSurveySubmission{
		PublicKey:      req.PublicKey,
		FullName:       req.FullName,
		Email:          req.Email,
		PhoneNumber:    req.PhoneNumber,
		VisitorType:    req.VisitorType,
		BusinessName:   req.BusinessName,
		Interests:      req.Interests,
		OtherInterests: req.OtherInterests,
		IsConsented:    req.IsConsented,
		IPAddress:      c.ClientIP(),
		UserAgent:      strings.TrimSpace(c.GetHeader("User-Agent")),
	})
	if err != nil {
		respondError(c, err)
		return
	}
	if visitor == nil {
		c.JSON(http.StatusInternalServerError, dto.Fail(dto.CodeInternalError, "visitor submission returned no app user"))
		return
	}
	c.JSON(http.StatusCreated, dto.OK(gin.H{"app_user_id": visitor.ID}))
}

// SearchOptions returns venue-scoped search filter options.
// @Summary     Public search options
// @Description Return venue-scoped exhibitor/location or product search options without requiring authentication.
// @Tags        public
// @Produce     json
// @Param       venueId path string true "Venue ID"
// @Param       origin query string false "Option origin" Enums(products, exhibitors, locations)
// @Param       lang query string false "Language code" default(en)
// @Success     200 {object} dto.Response
// @Failure     400 {object} dto.Response
// @Failure     404 {object} dto.Response
// @Failure     429 {object} dto.Response
// @Router      /public/venues/{venueId}/search-options [get]
func (h *publicHandler) SearchOptions(c *gin.Context) {
	venueID, err := uuid.Parse(c.Param("venueId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid venue id"))
		return
	}
	options, err := h.searchOptions.Get(c.Request.Context(), venueID, c.Query("origin"), c.Query("lang"))
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(options))
}
