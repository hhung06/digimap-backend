package dto

import "github.com/google/uuid"

// MediaUploadRequest is the JSON value of the multipart "data" field.
type MediaUploadRequest struct {
	Entity   string    `json:"entity"`
	RecordID uuid.UUID `json:"record_id"`
	Field    string    `json:"field"`
}

// MediaUploadResponse is returned after a successful upload.
// Key is the S3 object key stored in the base field.
// URL is a presigned download URL (nil if the key does not match the expected layout).
type MediaUploadResponse struct {
	Key string  `json:"key"`
	URL *string `json:"url"`
}
