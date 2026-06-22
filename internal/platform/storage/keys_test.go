package storage

import (
	"testing"

	"github.com/google/uuid"
)

func TestMediaKeyIsUniqueAndScoped(t *testing.T) {
	recordID := uuid.New()
	firstUploadID := uuid.New()
	first := MediaKey("develop", "articles", recordID, "images", firstUploadID, "PNG")
	second := MediaKey("develop", "articles", recordID, "images", uuid.New(), ".PNG")

	if first == second {
		t.Fatal("expected different upload IDs to produce different media keys")
	}
	want := "develop/media/articles/" + recordID.String() + "/images/" + firstUploadID.String() + ".png"
	if first != want {
		t.Fatalf("expected canonical key %q, got %q", want, first)
	}
}

func TestMediaKeyRejectsInvalidSegmentsAndExtensions(t *testing.T) {
	recordID := uuid.New()
	uploadID := uuid.New()
	tests := []struct {
		name   string
		env    string
		entity string
		field  string
		ext    string
	}{
		{name: "empty environment", env: "", entity: "articles", field: "images", ext: "png"},
		{name: "environment traversal", env: "..", entity: "articles", field: "images", ext: "png"},
		{name: "entity slash", env: "develop", entity: "content/articles", field: "images", ext: "png"},
		{name: "entity backslash", env: "develop", entity: `content\articles`, field: "images", ext: "png"},
		{name: "field dot", env: "develop", entity: "articles", field: ".", ext: "png"},
		{name: "field traversal", env: "develop", entity: "articles", field: "../images", ext: "png"},
		{name: "empty extension", env: "develop", entity: "articles", field: "images", ext: ""},
		{name: "dot extension", env: "develop", entity: "articles", field: "images", ext: "."},
		{name: "punctuated extension", env: "develop", entity: "articles", field: "images", ext: "p-ng"},
		{name: "compound extension", env: "develop", entity: "articles", field: "images", ext: "tar.gz"},
		{name: "extension traversal", env: "develop", entity: "articles", field: "images", ext: "../png"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MediaKey(tt.env, tt.entity, recordID, tt.field, uploadID, tt.ext); got != "" {
				t.Fatalf("expected invalid input to return empty key, got %q", got)
			}
		})
	}
}

func TestOwnsMediaKeyRequiresExactCanonicalStructure(t *testing.T) {
	recordID := uuid.New()
	uploadID := uuid.New()
	key := MediaKey("develop", "articles", recordID, "images", uploadID, "png")

	if !OwnsMediaKey("develop", "articles", recordID, "images", key) {
		t.Fatal("expected record to own its media key")
	}

	malformed := map[string]string{
		"another environment":       "staging/media/articles/" + recordID.String() + "/images/" + uploadID.String() + ".png",
		"another entity":            "develop/media/events/" + recordID.String() + "/images/" + uploadID.String() + ".png",
		"another record":            "develop/media/articles/" + uuid.New().String() + "/images/" + uploadID.String() + ".png",
		"record prefix confusion":   "develop/media/articles/" + recordID.String() + "-other/images/" + uploadID.String() + ".png",
		"another field":             "develop/media/articles/" + recordID.String() + "/thumbnail/" + uploadID.String() + ".png",
		"entity traversal":          "develop/media/../articles/" + recordID.String() + "/images/" + uploadID.String() + ".png",
		"field traversal":           "develop/media/articles/" + recordID.String() + "/../images/" + uploadID.String() + ".png",
		"malformed upload UUID":     "develop/media/articles/" + recordID.String() + "/images/not-a-uuid.png",
		"missing extension":         "develop/media/articles/" + recordID.String() + "/images/" + uploadID.String(),
		"uppercase extension":       "develop/media/articles/" + recordID.String() + "/images/" + uploadID.String() + ".PNG",
		"nonalphanumeric extension": "develop/media/articles/" + recordID.String() + "/images/" + uploadID.String() + ".p-ng",
	}

	for name, candidate := range malformed {
		t.Run(name, func(t *testing.T) {
			if OwnsMediaKey("develop", "articles", recordID, "images", candidate) {
				t.Fatalf("expected malformed or foreign key %q to be rejected", candidate)
			}
		})
	}

	invalidTargets := []struct {
		name   string
		env    string
		entity string
		field  string
	}{
		{name: "invalid environment", env: "../develop", entity: "articles", field: "images"},
		{name: "invalid entity", env: "develop", entity: "../articles", field: "images"},
		{name: "invalid field", env: "develop", entity: "articles", field: `images\other`},
	}
	for _, tt := range invalidTargets {
		t.Run(tt.name, func(t *testing.T) {
			if OwnsMediaKey(tt.env, tt.entity, recordID, tt.field, key) {
				t.Fatal("expected invalid ownership target to be rejected")
			}
		})
	}
}
