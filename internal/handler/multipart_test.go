package handler

import (
	"bytes"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/textproto"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/hhung06/digimap-backend/internal/domain"
)

type multipartTestFile struct {
	field       string
	filename    string
	contentType string
	body        string
}

func TestMultipartBindDataValidJSON(t *testing.T) {
	c := multipartTestContext(t, map[string]string{"data": `{"name":"Expo","count":2}`}, nil)

	var dst struct {
		Name  string `json:"name"`
		Count int    `json:"count"`
	}
	if err := bindMultipartData(c, &dst, 1024); err != nil {
		t.Fatalf("bind multipart data: %v", err)
	}
	if dst.Name != "Expo" || dst.Count != 2 {
		t.Fatalf("bound data = %+v", dst)
	}
}

func TestMultipartFileReturnsOneFile(t *testing.T) {
	c := multipartTestContext(t, map[string]string{"data": `{}`}, []multipartTestFile{
		{field: "image", filename: "hero.png", contentType: "image/png", body: "png body"},
	})
	if err := bindMultipartData(c, &struct{}{}, 1024); err != nil {
		t.Fatalf("bind multipart data: %v", err)
	}

	file, found, err := multipartFile(c, "image")
	if err != nil {
		t.Fatalf("multipart file: %v", err)
	}
	if !found {
		t.Fatal("expected file to be found")
	}
	if file.Filename != "hero.png" {
		t.Fatalf("filename = %q", file.Filename)
	}
}

func TestMultipartFilesReturnsRepeatedFiles(t *testing.T) {
	c := multipartTestContext(t, map[string]string{"data": `{}`}, []multipartTestFile{
		{field: "gallery", filename: "one.png", contentType: "image/png", body: "one"},
		{field: "gallery", filename: "two.png", contentType: "image/png", body: "two"},
	})
	if err := bindMultipartData(c, &struct{}{}, 1024); err != nil {
		t.Fatalf("bind multipart data: %v", err)
	}

	files, found, err := multipartFiles(c, "gallery")
	if err != nil {
		t.Fatalf("multipart files: %v", err)
	}
	if !found {
		t.Fatal("expected files to be found")
	}
	if len(files) != 2 {
		t.Fatalf("file count = %d", len(files))
	}
}

func TestMultipartFileHelpersReturnAbsentWhenFieldMissing(t *testing.T) {
	c := multipartTestContext(t, map[string]string{"data": `{}`}, nil)
	if err := bindMultipartData(c, &struct{}{}, 1024); err != nil {
		t.Fatalf("bind multipart data: %v", err)
	}

	file, found, err := multipartFile(c, "image")
	if err != nil {
		t.Fatalf("multipart file: %v", err)
	}
	if found || file != nil {
		t.Fatalf("single file found = %v, file = %v", found, file)
	}

	files, found, err := multipartFiles(c, "gallery")
	if err != nil {
		t.Fatalf("multipart files: %v", err)
	}
	if found || files != nil {
		t.Fatalf("files found = %v, files = %v", found, files)
	}
}

func TestMultipartBindDataAllowsFileLargerThanDataLimit(t *testing.T) {
	c := multipartTestContext(t, map[string]string{"data": `{}`}, []multipartTestFile{
		{field: "image", filename: "large.png", contentType: "image/png", body: strings.Repeat("x", 128)},
	})

	if err := bindMultipartData(c, &struct{}{}, 16); err != nil {
		t.Fatalf("bind multipart data: %v", err)
	}

	file, found, err := multipartFile(c, "image")
	if err != nil {
		t.Fatalf("multipart file: %v", err)
	}
	if !found {
		t.Fatal("expected file larger than data limit to be found")
	}
	if file.Size != 128 {
		t.Fatalf("file size = %d, want 128", file.Size)
	}
}

func TestMultipartBindDataRejectsMalformedJSON(t *testing.T) {
	c := multipartTestContext(t, map[string]string{"data": `{"name":`}, nil)

	err := bindMultipartData(c, &struct{}{}, 1024)
	if !errors.Is(err, domain.ErrValidation) {
		t.Fatalf("error = %v, want validation", err)
	}
}

func TestMultipartBindDataRejectsMissingData(t *testing.T) {
	c := multipartTestContext(t, nil, nil)

	err := bindMultipartData(c, &struct{}{}, 1024)
	if !errors.Is(err, domain.ErrValidation) {
		t.Fatalf("error = %v, want validation", err)
	}
}

func TestMultipartBindDataRejectsOversizeDataField(t *testing.T) {
	c := multipartTestContext(t, map[string]string{"data": `{"name":"too large"}`}, nil)

	err := bindMultipartData(c, &struct{}{}, 8)
	if !errors.Is(err, domain.ErrValidation) {
		t.Fatalf("error = %v, want validation", err)
	}
}

func TestMultipartBindDataRejectsOversizeRequest(t *testing.T) {
	c := multipartTestContext(t, map[string]string{"data": `{}`}, []multipartTestFile{
		{field: "image", filename: "too-large.png", contentType: "image/png", body: strings.Repeat("x", int(maxMultipartRequestBytes)+1)},
	})

	err := bindMultipartData(c, &struct{}{}, 1024)
	if !errors.Is(err, domain.ErrValidation) {
		t.Fatalf("error = %v, want validation", err)
	}
}

func TestMultipartFileRejectsMultipleFiles(t *testing.T) {
	c := multipartTestContext(t, map[string]string{"data": `{}`}, []multipartTestFile{
		{field: "image", filename: "one.png", contentType: "image/png", body: "one"},
		{field: "image", filename: "two.png", contentType: "image/png", body: "two"},
	})
	if err := bindMultipartData(c, &struct{}{}, 1024); err != nil {
		t.Fatalf("bind multipart data: %v", err)
	}

	_, _, err := multipartFile(c, "image")
	if !errors.Is(err, domain.ErrValidation) {
		t.Fatalf("error = %v, want validation", err)
	}
}

func TestMultipartMediaUploadOpensFileAndPopulatesMetadata(t *testing.T) {
	c := multipartTestContext(t, map[string]string{"data": `{}`}, []multipartTestFile{
		{field: "image", filename: "hero.png", contentType: "image/png", body: "file body"},
	})
	if err := bindMultipartData(c, &struct{}{}, 1024); err != nil {
		t.Fatalf("bind multipart data: %v", err)
	}
	file, found, err := multipartFile(c, "image")
	if err != nil {
		t.Fatalf("multipart file: %v", err)
	}
	if !found {
		t.Fatal("expected file to be found")
	}

	upload, closer, err := mediaUpload(file)
	if err != nil {
		t.Fatalf("media upload: %v", err)
	}
	defer closer.Close()

	if upload.Filename != "hero.png" {
		t.Fatalf("filename = %q", upload.Filename)
	}
	if upload.ContentType != "image/png" {
		t.Fatalf("content type = %q", upload.ContentType)
	}
	if upload.Size != int64(len("file body")) {
		t.Fatalf("size = %d", upload.Size)
	}
	body, err := io.ReadAll(upload.Reader)
	if err != nil {
		t.Fatalf("read upload: %v", err)
	}
	if string(body) != "file body" {
		t.Fatalf("body = %q", body)
	}
}

func multipartTestContext(t *testing.T, fields map[string]string, files []multipartTestFile) *gin.Context {
	t.Helper()

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	for key, value := range fields {
		if err := writer.WriteField(key, value); err != nil {
			t.Fatalf("write field: %v", err)
		}
	}
	for _, file := range files {
		header := make(textproto.MIMEHeader)
		header.Set("Content-Disposition", `form-data; name="`+file.field+`"; filename="`+file.filename+`"`)
		header.Set("Content-Type", file.contentType)
		part, err := writer.CreatePart(header)
		if err != nil {
			t.Fatalf("create part: %v", err)
		}
		if _, err := part.Write([]byte(file.body)); err != nil {
			t.Fatalf("write file: %v", err)
		}
	}
	if err := writer.Close(); err != nil {
		t.Fatalf("close writer: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/", &body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	rec := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(rec)
	c.Request = req
	return c
}
