package i18n

import (
	"time"

	"github.com/Abraxas-365/hada-commerce/internal/kernel"
)

// Translation represents a single translated field value for an entity in a specific locale.
type Translation struct {
	ID         kernel.TranslationID `json:"id"          db:"id"`
	TenantID   kernel.TenantID      `json:"tenant_id"   db:"tenant_id"`
	EntityType string               `json:"entity_type" db:"entity_type"`
	EntityID   string               `json:"entity_id"   db:"entity_id"`
	Locale     string               `json:"locale"      db:"locale"`
	Field      string               `json:"field"       db:"field"`
	Value      string               `json:"value"       db:"value"`
	CreatedAt  time.Time            `json:"created_at"  db:"created_at"`
	UpdatedAt  time.Time            `json:"updated_at"  db:"updated_at"`
}

// TranslationBundle groups all translated fields for a single entity+locale combination.
type TranslationBundle struct {
	EntityType string            `json:"entity_type"`
	EntityID   string            `json:"entity_id"`
	Locale     string            `json:"locale"`
	Fields     map[string]string `json:"fields"`
}

// LocaleInfo describes a supported locale.
type LocaleInfo struct {
	Code      string `json:"code"`
	Name      string `json:"name"`
	IsDefault bool   `json:"is_default"`
}

// SupportedLocales is the static map of supported locales.
var SupportedLocales = map[string]LocaleInfo{
	"en":    {Code: "en", Name: "English", IsDefault: true},
	"es":    {Code: "es", Name: "Spanish", IsDefault: false},
	"fr":    {Code: "fr", Name: "French", IsDefault: false},
	"de":    {Code: "de", Name: "German", IsDefault: false},
	"pt-BR": {Code: "pt-BR", Name: "Portuguese (Brazil)", IsDefault: false},
	"ja":    {Code: "ja", Name: "Japanese", IsDefault: false},
	"zh":    {Code: "zh", Name: "Chinese", IsDefault: false},
	"ko":    {Code: "ko", Name: "Korean", IsDefault: false},
	"it":    {Code: "it", Name: "Italian", IsDefault: false},
	"nl":    {Code: "nl", Name: "Dutch", IsDefault: false},
	"ar":    {Code: "ar", Name: "Arabic", IsDefault: false},
}

// ValidEntityTypes contains the allowed entity type values.
var ValidEntityTypes = map[string]bool{
	"product":    true,
	"category":   true,
	"collection": true,
	"page":       true,
}

// ValidFields contains the allowed translatable field names.
var ValidFields = map[string]bool{
	"name":             true,
	"description":      true,
	"meta_title":       true,
	"meta_description": true,
}
