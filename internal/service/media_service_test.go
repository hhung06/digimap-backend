package service_test

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"hash/crc32"
	"image"
	"image/color"
	"image/png"
	"io"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hhung06/digimap-backend/internal/domain"
	"github.com/hhung06/digimap-backend/internal/platform/storage"
	"github.com/hhung06/digimap-backend/internal/service"
)

type mediaPutCall struct {
	key           string
	contentType   string
	contentLength int64
	body          []byte
}

type mediaStorerSpy struct {
	puts        []mediaPutCall
	deleted     []string
	presignKeys []string
	presignTTL  time.Duration
	presignErr  error
}

func (s *mediaStorerSpy) PresignUpload(context.Context, string, string, time.Duration) (string, error) {
	return "", nil
}

func (s *mediaStorerSpy) PresignDownload(_ context.Context, key string, ttl time.Duration) (string, error) {
	s.presignKeys = append(s.presignKeys, key)
	s.presignTTL = ttl
	if s.presignErr != nil {
		return "", s.presignErr
	}
	return "https://assets.example.com/" + key, nil
}

func (s *mediaStorerSpy) PutMedia(_ context.Context, key, contentType string, contentLength int64, body io.Reader) error {
	b, err := io.ReadAll(body)
	if err != nil {
		return err
	}
	s.puts = append(s.puts, mediaPutCall{
		key:           key,
		contentType:   contentType,
		contentLength: contentLength,
		body:          b,
	})
	return nil
}

func (s *mediaStorerSpy) PutObject(context.Context, string, []byte) error { return nil }

func (s *mediaStorerSpy) PutEncrypted(context.Context, string, []byte, map[string]string) error {
	return nil
}

func (s *mediaStorerSpy) GetObject(context.Context, string) ([]byte, error) { return nil, nil }

func (s *mediaStorerSpy) DeleteObject(_ context.Context, key string) error {
	s.deleted = append(s.deleted, key)
	return nil
}

func TestMediaService_UploadSameFilenameTwiceProducesDifferentOwnedKeys(t *testing.T) {
	ctx := context.Background()
	storer := &mediaStorerSpy{}
	svc := service.NewMediaService(storer, "develop")
	target := service.MediaTarget{Entity: "articles", RecordID: uuid.New(), Field: "image"}
	body := validPNG(t)

	first, err := svc.Upload(ctx, target, service.MediaUpload{
		Filename:    "hero.png",
		ContentType: "image/png",
		Size:        int64(len(body)),
		Reader:      bytes.NewReader(body),
	})
	require.NoError(t, err)

	second, err := svc.Upload(ctx, target, service.MediaUpload{
		Filename:    "hero.png",
		ContentType: "image/png",
		Size:        int64(len(body)),
		Reader:      bytes.NewReader(body),
	})
	require.NoError(t, err)

	assert.NotEqual(t, first, second)
	assert.True(t, storage.OwnsMediaKey("develop", target.Entity, target.RecordID, target.Field, first))
	assert.True(t, storage.OwnsMediaKey("develop", target.Entity, target.RecordID, target.Field, second))
	require.Len(t, storer.puts, 2)
	assert.NotEqual(t, storer.puts[0].key, storer.puts[1].key)
}

func TestMediaService_UploadAllowsPDFDocument(t *testing.T) {
	ctx := context.Background()
	storer := &mediaStorerSpy{}
	svc := service.NewMediaService(storer, "develop")
	target := service.MediaTarget{Entity: "products", RecordID: uuid.New(), Field: "attachments"}
	body := []byte("%PDF-1.4\ncontent")

	key, err := svc.Upload(ctx, target, service.MediaUpload{
		Filename:    "spec.pdf",
		ContentType: "application/pdf",
		Size:        int64(len(body)),
		Reader:      bytes.NewReader(body),
	})
	require.NoError(t, err)

	assert.True(t, storage.OwnsMediaKey("develop", target.Entity, target.RecordID, target.Field, key))
	require.Len(t, storer.puts, 1)
	assert.Equal(t, "application/pdf", storer.puts[0].contentType)
	assert.Equal(t, body, storer.puts[0].body)
}

func TestMediaService_UploadSameFilenameForDifferentRecordsProducesDifferentOwnedKeys(t *testing.T) {
	ctx := context.Background()
	storer := &mediaStorerSpy{}
	svc := service.NewMediaService(storer, "develop")
	body := validPNG(t)
	firstTarget := service.MediaTarget{Entity: "articles", RecordID: uuid.New(), Field: "image"}
	secondTarget := service.MediaTarget{Entity: "articles", RecordID: uuid.New(), Field: "image"}
	upload := func(target service.MediaTarget) string {
		t.Helper()
		key, err := svc.Upload(ctx, target, service.MediaUpload{
			Filename:    "hero.png",
			ContentType: "image/png",
			Size:        int64(len(body)),
			Reader:      bytes.NewReader(body),
		})
		require.NoError(t, err)
		return key
	}

	first := upload(firstTarget)
	second := upload(secondTarget)

	assert.NotEqual(t, first, second)
	assert.True(t, storage.OwnsMediaKey("develop", firstTarget.Entity, firstTarget.RecordID, firstTarget.Field, first))
	assert.True(t, storage.OwnsMediaKey("develop", secondTarget.Entity, secondTarget.RecordID, secondTarget.Field, second))
	assert.False(t, storage.OwnsMediaKey("develop", firstTarget.Entity, firstTarget.RecordID, firstTarget.Field, second))
}

func TestMediaService_UploadRejectsMIMEMismatch(t *testing.T) {
	ctx := context.Background()
	storer := &mediaStorerSpy{}
	svc := service.NewMediaService(storer, "develop")
	body := validPNG(t)

	_, err := svc.Upload(ctx, service.MediaTarget{Entity: "articles", RecordID: uuid.New(), Field: "image"}, service.MediaUpload{
		Filename:    "hero.jpg",
		ContentType: "image/jpeg",
		Size:        int64(len(body)),
		Reader:      bytes.NewReader(body),
	})

	require.Error(t, err)
	assert.Empty(t, storer.puts)
}

func TestMediaService_UploadRejectsOversizedInput(t *testing.T) {
	ctx := context.Background()
	storer := &mediaStorerSpy{}
	svc := service.NewMediaService(storer, "develop")

	_, err := svc.Upload(ctx, service.MediaTarget{Entity: "articles", RecordID: uuid.New(), Field: "image"}, service.MediaUpload{
		Filename:    "hero.png",
		ContentType: "image/png",
		Size:        10<<20 + 1,
		Reader:      bytes.NewReader(validPNG(t)),
	})

	require.Error(t, err)
	assert.Empty(t, storer.puts)
}

func TestMediaService_UploadRejectsUndecodableImage(t *testing.T) {
	ctx := context.Background()
	storer := &mediaStorerSpy{}
	svc := service.NewMediaService(storer, "develop")
	body := []byte{0x89, 'P', 'N', 'G', 0x0d, 0x0a, 0x1a, 0x0a}

	_, err := svc.Upload(ctx, service.MediaTarget{Entity: "articles", RecordID: uuid.New(), Field: "image"}, service.MediaUpload{
		Filename:    "hero.png",
		ContentType: "image/png",
		Size:        int64(len(body)),
		Reader:      bytes.NewReader(body),
	})

	require.Error(t, err)
	assert.Empty(t, storer.puts)
}

func TestMediaService_UploadRejectsImageDimensionsBeforeDecode(t *testing.T) {
	ctx := context.Background()
	storer := &mediaStorerSpy{}
	svc := service.NewMediaService(storer, "develop")
	body := pngWithDimensions(t, 10000, 10000)

	_, err := svc.Upload(ctx, service.MediaTarget{Entity: "articles", RecordID: uuid.New(), Field: "image"}, service.MediaUpload{
		Filename:    "huge.png",
		ContentType: "image/png",
		Size:        int64(len(body)),
		Reader:      bytes.NewReader(body),
	})

	require.Error(t, err)
	var appErr *domain.AppError
	require.True(t, errors.As(err, &appErr))
	assert.Equal(t, "dimensions too large", appErr.Details["image"])
	assert.Empty(t, storer.puts)
}

func TestMediaService_URLSignsOwnedKey(t *testing.T) {
	ctx := context.Background()
	storer := &mediaStorerSpy{}
	svc := service.NewMediaService(storer, "develop")
	target := service.MediaTarget{Entity: "articles", RecordID: uuid.New(), Field: "image"}
	key := storage.MediaKey("develop", target.Entity, target.RecordID, target.Field, uuid.New(), "png")

	url := svc.URL(ctx, target, key)

	require.NotNil(t, url)
	assert.Equal(t, "https://assets.example.com/"+key, *url)
	assert.Equal(t, 15*time.Minute, storer.presignTTL)
	assert.Equal(t, []string{key}, storer.presignKeys)
}

func TestMediaService_URLReturnsNilForForeignKeyAndDoesNotPresign(t *testing.T) {
	ctx := context.Background()
	target := service.MediaTarget{Entity: "articles", RecordID: uuid.New(), Field: "image"}

	tests := []struct {
		name string
		key  string
	}{
		{
			name: "another record",
			key:  storage.MediaKey("develop", target.Entity, uuid.New(), target.Field, uuid.New(), "png"),
		},
		{
			name: "another entity",
			key:  storage.MediaKey("develop", "events", target.RecordID, target.Field, uuid.New(), "png"),
		},
		{
			name: "another field",
			key:  storage.MediaKey("develop", target.Entity, target.RecordID, "thumbnail", uuid.New(), "png"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storer := &mediaStorerSpy{}
			svc := service.NewMediaService(storer, "develop")

			url := svc.URL(ctx, target, tt.key)

			require.Nil(t, url)
			assert.Empty(t, storer.presignKeys)
		})
	}
}

func TestMediaService_URLReturnsNilForBlankKeyAndDoesNotPresign(t *testing.T) {
	ctx := context.Background()
	storer := &mediaStorerSpy{}
	svc := service.NewMediaService(storer, "develop")
	target := service.MediaTarget{Entity: "articles", RecordID: uuid.New(), Field: "image"}

	url := svc.URL(ctx, target, "")

	require.Nil(t, url)
	assert.Empty(t, storer.presignKeys)
}

func TestMediaService_URLReturnsNilWhenPresignFails(t *testing.T) {
	ctx := context.Background()
	storer := &mediaStorerSpy{presignErr: assert.AnError}
	svc := service.NewMediaService(storer, "develop")
	target := service.MediaTarget{Entity: "articles", RecordID: uuid.New(), Field: "image"}
	key := storage.MediaKey("develop", target.Entity, target.RecordID, target.Field, uuid.New(), "png")

	url := svc.URL(ctx, target, key)

	require.Nil(t, url)
	assert.Equal(t, 15*time.Minute, storer.presignTTL)
	assert.Equal(t, []string{key}, storer.presignKeys)
}

func TestMediaService_DeleteOwnedRejectsAnotherRecordsKey(t *testing.T) {
	ctx := context.Background()
	storer := &mediaStorerSpy{}
	svc := service.NewMediaService(storer, "develop")
	target := service.MediaTarget{Entity: "articles", RecordID: uuid.New(), Field: "image"}
	otherKey := storage.MediaKey("develop", target.Entity, uuid.New(), target.Field, uuid.New(), "png")

	err := svc.DeleteOwned(ctx, target, otherKey)

	require.Error(t, err)
	assert.Empty(t, storer.deleted)
}

func TestMediaService_DeleteOwnedDeletesOwnedKey(t *testing.T) {
	ctx := context.Background()
	storer := &mediaStorerSpy{}
	svc := service.NewMediaService(storer, "develop")
	target := service.MediaTarget{Entity: "articles", RecordID: uuid.New(), Field: "image"}
	key := storage.MediaKey("develop", target.Entity, target.RecordID, target.Field, uuid.New(), "png")

	err := svc.DeleteOwned(ctx, target, key)

	require.NoError(t, err)
	require.Equal(t, []string{key}, storer.deleted)
}

func TestMediaService_UploadPassesValidatedContentToPutMedia(t *testing.T) {
	ctx := context.Background()
	storer := &mediaStorerSpy{}
	svc := service.NewMediaService(storer, "develop")
	body := validPNG(t)

	key, err := svc.Upload(ctx, service.MediaTarget{Entity: "articles", RecordID: uuid.New(), Field: "image"}, service.MediaUpload{
		Filename:    "hero.bin",
		ContentType: "image/png",
		Size:        int64(len(body)),
		Reader:      bytes.NewReader(body),
	})

	require.NoError(t, err)
	require.Len(t, storer.puts, 1)
	assert.Equal(t, key, storer.puts[0].key)
	assert.Equal(t, "image/png", storer.puts[0].contentType)
	assert.Equal(t, int64(len(body)), storer.puts[0].contentLength)
	assert.Equal(t, body, storer.puts[0].body)
}

func validPNG(t *testing.T) []byte {
	t.Helper()

	img := image.NewRGBA(image.Rect(0, 0, 1, 1))
	img.Set(0, 0, color.RGBA{R: 255, A: 255})
	var buf bytes.Buffer
	require.NoError(t, png.Encode(&buf, img))
	return buf.Bytes()
}

func pngWithDimensions(t *testing.T, width, height uint32) []byte {
	t.Helper()

	var body bytes.Buffer
	body.Write([]byte{0x89, 'P', 'N', 'G', 0x0d, 0x0a, 0x1a, 0x0a})
	writePNGChunk(t, &body, "IHDR", func() []byte {
		data := make([]byte, 13)
		binary.BigEndian.PutUint32(data[0:4], width)
		binary.BigEndian.PutUint32(data[4:8], height)
		data[8] = 8
		data[9] = 2
		return data
	}())
	writePNGChunk(t, &body, "IEND", nil)
	return body.Bytes()
}

func writePNGChunk(t *testing.T, dst *bytes.Buffer, name string, data []byte) {
	t.Helper()

	require.Len(t, name, 4)
	require.NoError(t, binary.Write(dst, binary.BigEndian, uint32(len(data))))
	dst.WriteString(name)
	dst.Write(data)
	crc := crc32.ChecksumIEEE(append([]byte(name), data...))
	require.NoError(t, binary.Write(dst, binary.BigEndian, crc))
}
