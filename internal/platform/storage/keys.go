package storage

import (
	"path"
	"strings"

	"github.com/google/uuid"
)

// Key builders mirror Django's generate_object_key (indoormap-backend/utils/s3services.py:32).
// All keys are prefixed with env (e.g., "production", "staging", "local").

// DigimapV2Key returns the per-language encrypted bundle key for snapshot publish.
// Django: {env}/bundles/public/{venue_id}.digiapp.{lang}
func DigimapV2Key(env string, venueID uuid.UUID, lang string) string {
	return env + "/bundles/public/" + venueID.String() + ".digiapp." + lang
}

// TopLocationKey returns the encrypted top-location bundle key.
// Django: {env}/top_location/public/{venue_id}.digimap
func TopLocationKey(env string, venueID uuid.UUID) string {
	return env + "/top_location/public/" + venueID.String() + ".digimap"
}

// MemoKey returns the encrypted memo bundle key.
// Django: {env}/location_memo/public/{venue_id}.digimap
func MemoKey(env string, venueID uuid.UUID) string {
	return env + "/location_memo/public/" + venueID.String() + ".digimap"
}

// LatestBundleKey returns the plaintext version-pointer key (not encrypted).
// Django: {env}/latest-bundle/{venue_id}.json
func LatestBundleKey(env string, venueID uuid.UUID) string {
	return env + "/latest-bundle/" + venueID.String() + ".json"
}

// SnapshotDraftKey returns the key for a stored snapshot draft JSON blob.
// Django: {env}/{venue_id}/snapshots/{state}/{snapshot_id}.json
func SnapshotDraftKey(env string, venueID, snapshotID uuid.UUID) string {
	return env + "/" + venueID.String() + "/snapshots/draft/" + snapshotID.String() + ".json"
}

// LevelBundleKey returns the S3 key for a per-level bundle blob.
func LevelBundleKey(env string, venueID, snapshotID, levelID uuid.UUID) string {
	return env + "/" + venueID.String() + "/level-bundles/" + snapshotID.String() + "/" + levelID.String() + ".json"
}

// GlobalThemeKey returns the S3 key for a global venue theme JSON file.
// Django: {env}/venue_themes/global/{name}.json
func GlobalThemeKey(env, name string) string {
	return env + "/venue_themes/global/" + name + ".json"
}

// CustomThemeKey returns the S3 key for a venue-scoped custom theme JSON file.
// Django: {env}/venue_themes/custom/{venue_id}/{name}.json
func CustomThemeKey(env string, venueID uuid.UUID, name string) string {
	return env + "/venue_themes/custom/" + venueID.String() + "/" + name + ".json"
}

// LibraryAssetKey returns the S3 key for a library asset.
// Django equivalent: {env}/library/{id}_{filename}
func LibraryAssetKey(env string, id uuid.UUID, filename string) string {
	if env == "" || filename == "" {
		return ""
	}
	return env + "/library/" + id.String() + "_" + filename
}

// Asset2DKey returns the S3 key for a 2D asset (canvas image).
// ext should not include the leading dot (e.g. "png", "jpg").
func Asset2DKey(env string, id uuid.UUID, ext string) string {
	if env == "" || ext == "" {
		return ""
	}
	return env + "/assets/" + id.String() + "." + ext
}

// BaseKey returns the viewer-fetched mesh bundle key.
// Django: {env}/base/public/{venue_id}.digimap
func BaseKey(env string, venueID uuid.UUID) string {
	return env + "/base/public/" + venueID.String() + ".digimap"
}

// OverviewKey returns the viewer-fetched paths/theme bundle key.
// Django: {env}/overview/public/{venue_id}.digimap
func OverviewKey(env string, venueID uuid.UUID) string {
	return env + "/overview/public/" + venueID.String() + ".digimap"
}

// LocationSimpleKey returns the viewer-fetched minimal-location bundle key.
// Django: {env}/location_simple/public/{venue_id}.digimap
func LocationSimpleKey(env string, venueID uuid.UUID) string {
	return env + "/location_simple/public/" + venueID.String() + ".digimap"
}

// MetadataKey returns the viewer-fetched full metadata bundle key.
// Django: {env}/metadata/public/{venue_id}.digimap
func MetadataKey(env string, venueID uuid.UUID) string {
	return env + "/metadata/public/" + venueID.String() + ".digimap"
}

// LocalizedKey returns the per-language overlay key fetched by the viewer for non-English locales.
// Django: {env}/bundles/public/{venue_id}.digimap.{lang}
func LocalizedKey(env string, venueID uuid.UUID, lang string) string {
	return env + "/bundles/public/" + venueID.String() + ".digimap." + lang
}

// Asset3DKey returns the S3 key for a 3D asset file.
// suffix is "" for the main file, "_thumbnail" for thumbnail, "_material" for material.
// ext should not include the leading dot (e.g. "glb", "jpg").
// Django equivalent: {env}/assets_3D/{id}{suffix}.{ext}
func Asset3DKey(env string, id uuid.UUID, suffix, ext string) string {
	if env == "" || ext == "" {
		return ""
	}
	return env + "/assets_3D/" + id.String() + suffix + "." + ext
}

// MediaKey returns an immutable key for one backend-owned media upload.
func MediaKey(env, entity string, recordID uuid.UUID, field string, uploadID uuid.UUID, ext string) string {
	if !isCanonicalMediaSegment(env) || !isCanonicalMediaSegment(entity) || !isCanonicalMediaSegment(field) {
		return ""
	}
	normalizedExt, ok := normalizeMediaExtension(ext)
	if !ok {
		return ""
	}
	return path.Join(env, "media", entity, recordID.String(), field, uploadID.String()+normalizedExt)
}

// OwnsMediaKey reports whether key belongs to the specified record field.
func OwnsMediaKey(env, entity string, recordID uuid.UUID, field, key string) bool {
	if !isCanonicalMediaSegment(env) || !isCanonicalMediaSegment(entity) || !isCanonicalMediaSegment(field) {
		return false
	}
	parts := strings.Split(key, "/")
	if len(parts) != 6 || parts[0] != env || parts[1] != "media" || parts[2] != entity || parts[4] != field {
		return false
	}
	keyRecordID, err := uuid.Parse(parts[3])
	if err != nil || keyRecordID != recordID || keyRecordID.String() != parts[3] {
		return false
	}
	nameParts := strings.Split(parts[5], ".")
	if len(nameParts) != 2 {
		return false
	}
	uploadID, err := uuid.Parse(nameParts[0])
	if err != nil || uploadID.String() != nameParts[0] {
		return false
	}
	normalizedExt, ok := normalizeMediaExtension(nameParts[1])
	return ok && normalizedExt == "."+nameParts[1]
}

func isCanonicalMediaSegment(segment string) bool {
	return segment != "" && segment != "." && segment != ".." &&
		!strings.ContainsAny(segment, `/\`) && path.Clean(segment) == segment
}

func normalizeMediaExtension(ext string) (string, bool) {
	ext = strings.TrimPrefix(ext, ".")
	if ext == "" {
		return "", false
	}
	for _, r := range ext {
		if (r < 'a' || r > 'z') && (r < 'A' || r > 'Z') && (r < '0' || r > '9') {
			return "", false
		}
	}
	return "." + strings.ToLower(ext), true
}
