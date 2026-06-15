package storage

import "github.com/google/uuid"

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
