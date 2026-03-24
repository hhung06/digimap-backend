package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/hhung06/digimap-backend/internal/domain"
)

// paginationFromQuery extracts page/page_size from query params with sensible defaults.
func paginationFromQuery(c *gin.Context) domain.Pagination {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	p := domain.Pagination{Page: page, PageSize: pageSize}
	p.Normalize()
	return p
}

// bindingErrors converts Gin/validator binding errors into a flat list of
// human-readable strings in the format "field: reason".
func bindingErrors(err error) []string {
	var ve validator.ValidationErrors
	if ok := isValidationErrors(err, &ve); ok {
		msgs := make([]string, 0, len(ve))
		for _, fe := range ve {
			msgs = append(msgs, fe.Field()+": "+validationMsg(fe))
		}
		return msgs
	}
	return []string{err.Error()}
}

func isValidationErrors(err error, ve *validator.ValidationErrors) bool {
	if e, ok := err.(validator.ValidationErrors); ok {
		*ve = e
		return true
	}
	return false
}

func validationMsg(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "required"
	case "email":
		return "invalid email format"
	case "min":
		return "must be at least " + fe.Param() + " characters"
	case "max":
		return "must be at most " + fe.Param() + " characters"
	case "uuid":
		return "invalid UUID format"
	default:
		return fe.Tag()
	}
}
