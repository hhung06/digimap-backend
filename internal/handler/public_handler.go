package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/hhung06/digimap-backend/internal/domain"
	"github.com/hhung06/digimap-backend/internal/dto"
	"github.com/hhung06/digimap-backend/internal/service"
)

type publicHandler struct {
	venues        service.VenueService
	surveys       service.SurveyService
	productPlazas service.ProductPlazaService
}

func newPublicHandler(venues service.VenueService, surveys service.SurveyService, productPlazas service.ProductPlazaService) *publicHandler {
	return &publicHandler{venues: venues, surveys: surveys, productPlazas: productPlazas}
}

func (h *publicHandler) VenueInformation(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(1000, "invalid venue id"))
		return
	}
	v, err := h.venues.Get(c.Request.Context(), id)
	if err != nil {
		respondError(c, err)
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
		resp.Answers[i] = &domain.SurveyAnswer{
			QuestionID: a.QuestionID,
			OptionID:   a.OptionID,
			AnswerText: a.AnswerText,
		}
	}
	if err := h.surveys.SubmitResponse(c.Request.Context(), resp); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.OK(nil))
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
