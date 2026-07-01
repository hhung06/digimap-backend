package cmd

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"mime"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/spf13/cobra"

	"github.com/hhung06/digimap-backend/config"
	"github.com/hhung06/digimap-backend/internal/platform/database"
	"github.com/hhung06/digimap-backend/internal/platform/storage"
	"github.com/hhung06/digimap-backend/internal/service"
)

var repairProductAttachmentsCmd = &cobra.Command{
	Use:   "repair-product-attachments",
	Short: "Upload base64 product attachment rows to media storage",
	RunE:  runRepairProductAttachments,
}

func init() {
	rootCmd.AddCommand(repairProductAttachmentsCmd)
}

func runRepairProductAttachments(_ *cobra.Command, _ []string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}
	ctx := context.Background()
	pool, err := database.NewPool(ctx, cfg.Database)
	if err != nil {
		return fmt.Errorf("connect database: %w", err)
	}
	defer pool.Close()

	var storer storage.Storer
	if cfg.AWS.S3AssetsBucket != "" {
		storer, err = storage.NewS3Storer(cfg.AWS, cfg.AWS.S3AssetsBucket)
		if err != nil {
			return fmt.Errorf("init s3 assets storer: %w", err)
		}
	} else {
		storer = storage.NewLogStorer()
	}
	mediaSvc := service.NewMediaService(storer, cfg.App.Environment)

	rows, err := pool.Query(ctx, `
		SELECT id, product_id, title, file
		FROM product_attachments
		WHERE deleted_at IS NULL AND file LIKE 'data:%;base64,%'`)
	if err != nil {
		return err
	}
	defer rows.Close()

	repaired := 0
	for rows.Next() {
		var id, productID uuid.UUID
		var title, raw *string
		if err := rows.Scan(&id, &productID, &title, &raw); err != nil {
			return err
		}
		if raw == nil {
			continue
		}
		upload, err := productAttachmentUploadFromDataURL(*raw, title)
		if err != nil {
			return fmt.Errorf("attachment %s: %w", id, err)
		}
		target := service.MediaTarget{Entity: "products", RecordID: productID, Field: "attachments"}
		key, err := mediaSvc.Upload(ctx, target, upload)
		if err != nil {
			return fmt.Errorf("attachment %s upload: %w", id, err)
		}
		if _, err := pool.Exec(ctx, `UPDATE product_attachments SET file = $2 WHERE id = $1`, id, key); err != nil {
			_ = mediaSvc.DeleteOwned(ctx, target, key)
			return fmt.Errorf("attachment %s update: %w", id, err)
		}
		repaired++
	}
	if err := rows.Err(); err != nil {
		return err
	}
	fmt.Printf("repaired product attachments: %d\n", repaired)
	return nil
}

func productAttachmentUploadFromDataURL(raw string, title *string) (service.MediaUpload, error) {
	meta, encoded, ok := strings.Cut(raw, ",")
	if !ok || !strings.HasPrefix(meta, "data:") || !strings.Contains(meta, ";base64") {
		return service.MediaUpload{}, fmt.Errorf("invalid data URL")
	}
	contentType := strings.TrimPrefix(strings.Split(meta, ";")[0], "data:")
	body, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return service.MediaUpload{}, err
	}
	exts, _ := mime.ExtensionsByType(contentType)
	ext := ".bin"
	if len(exts) > 0 {
		ext = exts[0]
	}
	name := "attachment" + ext
	if title != nil && strings.TrimSpace(*title) != "" {
		name = strings.TrimSpace(*title)
		if filepath.Ext(name) == "" {
			name += ext
		}
	}
	return service.MediaUpload{
		Filename:    name,
		ContentType: contentType,
		Size:        int64(len(body)),
		Reader:      bytes.NewReader(body),
	}, nil
}
