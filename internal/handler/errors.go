package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hhung06/digimap-backend/internal/domain"
	"github.com/hhung06/digimap-backend/internal/dto"
)

// respondError maps a service/domain error to the correct HTTP status and
// AppCode, then writes a JSON error response. Handlers call this for every
// non-nil error returned from the service layer.
func respondError(c *gin.Context, err error) {
	var appErr *domain.AppError
	if errors.As(err, &appErr) {
		status, code := mapDomainError(appErr.Err)

		// Validation errors include field-level detail as extra messages
		if len(appErr.Details) > 0 {
			msgs := make([]string, 0, len(appErr.Details)+1)
			msgs = append(msgs, appErr.Message)
			for field, reason := range appErr.Details {
				msgs = append(msgs, field+": "+reason)
			}
			c.AbortWithStatusJSON(status, dto.FailMessages(code, msgs))
			return
		}

		c.AbortWithStatusJSON(status, dto.Fail(code, appErr.Message))
		return
	}

	// Unwrapped sentinel errors
	status, code := mapDomainError(err)
	c.AbortWithStatusJSON(status, dto.Fail(code, err.Error()))
}

// mapDomainError returns the HTTP status and AppCode for a given sentinel error.
func mapDomainError(err error) (int, dto.AppCode) {
	switch {
	case errors.Is(err, domain.ErrNotFound):
		return http.StatusNotFound, dto.CodeNotFound
	case errors.Is(err, domain.ErrUnauthorized):
		return http.StatusUnauthorized, dto.CodeAuthRequired
	case errors.Is(err, domain.ErrForbidden):
		return http.StatusForbidden, dto.CodePermissionDenied
	case errors.Is(err, domain.ErrConflict):
		return http.StatusConflict, dto.CodeConflict
	case errors.Is(err, domain.ErrValidation), errors.Is(err, domain.ErrBadRequest):
		return http.StatusBadRequest, dto.CodeValidationError
	default:
		return http.StatusInternalServerError, dto.CodeInternalError
	}
}
