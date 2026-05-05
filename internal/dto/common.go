package dto

// AppCode is the business-level status code included in every response.
// HTTP status codes signal the transport layer; AppCode signals the business layer.
// Clients should treat any code != 0 as an error regardless of the HTTP status.
type AppCode int

const (
	CodeSuccess           AppCode = 0    // Request processed successfully.
	CodeValidationError   AppCode = 1000 // Invalid input: missing field, wrong type, or format violation.
	CodeAuthRequired      AppCode = 1001 // Request requires authentication; no valid token provided.
	CodePermissionDenied  AppCode = 1002 // Authenticated but lacks permission to perform this action.
	CodeNotFound          AppCode = 1003 // Requested resource does not exist.
	CodeConflict          AppCode = 1004 // Duplicate resource or incompatible state.
	CodeBusinessLogic     AppCode = 1005 // Valid request but invalid in the current context.
	CodeRateLimited       AppCode = 1006 // Too many requests; client should back off.
	CodeInternalError     AppCode = 1007 // Unexpected server-side error.
)

// Response is the standard envelope for all endpoints.
//
//	Success (single):   {"code": 0, "data": {...}}
//	Success (list):     {"code": 0, "data": [...], "metadata": {"total_records": N, ...}}
//	Error:              {"code": 1003, "messages": ["venue not found"]}
type Response struct {
	Code     AppCode     `json:"code"`
	Data     interface{} `json:"data,omitempty"`
	Metadata interface{} `json:"metadata,omitempty"`
	Messages []string    `json:"messages,omitempty"`
}

// PaginationMeta holds page metadata returned with every list response.
type PaginationMeta struct {
	TotalRecords int  `json:"total_records"`
	TotalPages   int  `json:"total_pages"`
	CurrentPage  int  `json:"current_page"`
	PageSize     int  `json:"page_size"`
	HasNextPage  bool `json:"has_next_page"`
}

// PaginationQuery holds query parameters for list endpoints.
type PaginationQuery struct {
	Page     int `form:"page"`
	PageSize int `form:"page_size"`
}

// OK returns a successful Response wrapping data.
func OK(data interface{}) Response {
	return Response{Code: CodeSuccess, Data: data}
}

// Paginated returns a successful paginated Response.
func Paginated(items interface{}, total int64, page, pageSize int) Response {
	totalPages := int(total) / pageSize
	if int(total)%pageSize != 0 {
		totalPages++
	}
	return Response{
		Code: CodeSuccess,
		Data: items,
		Metadata: PaginationMeta{
			TotalRecords: int(total),
			TotalPages:   totalPages,
			CurrentPage:  page,
			PageSize:     pageSize,
			HasNextPage:  page < totalPages,
		},
	}
}

// OKMessage returns a successful Response with no data, only a human-readable message.
func OKMessage(msg string) Response {
	return Response{Code: CodeSuccess, Messages: []string{msg}}
}

// Fail returns an error Response with a single message.
func Fail(code AppCode, msg string) Response {
	return Response{Code: code, Messages: []string{msg}}
}

// FailMessages returns an error Response with multiple messages.
func FailMessages(code AppCode, msgs []string) Response {
	return Response{Code: code, Messages: msgs}
}
